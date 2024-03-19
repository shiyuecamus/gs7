// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package logging

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Flusher is the callback function which flushes any buffered log entries to the underlying writer.
// It is usually called before the s7 process exits.
type Flusher = func() error

var (
	mu                  sync.RWMutex
	defaultLogger       Logger
	defaultLoggingLevel Level
	defaultFlusher      Flusher
)

// Level is the alias of zapcore.Level.
type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

func init() {
	lvl := os.Getenv("S7_LOGGING_LEVEL")
	if len(lvl) > 0 {
		loggingLevel, err := strconv.ParseInt(lvl, 10, 8)
		if err != nil {
			panic("invalid S7_LOGGING_LEVEL, " + err.Error())
		}
		defaultLoggingLevel = Level(loggingLevel)
	}

	// Initializes the inside default logger of s7.
	fileName := os.Getenv("S7_LOGGING_FILE")
	if len(fileName) > 0 {
		var err error
		defaultLogger, defaultFlusher, err = CreateLoggerAsLocalFile(fileName, defaultLoggingLevel)
		if err != nil {
			panic("invalid S7_LOGGING_FILE, " + err.Error())
		}
	} else {
		core := zapcore.NewCore(getDevEncoder(), zapcore.Lock(os.Stdout), defaultLoggingLevel)
		zapLogger := zap.New(core,
			zap.Development(),
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.ErrorOutput(zapcore.Lock(os.Stderr)))
		defaultLogger = zapLogger.Sugar()
	}
}

type prefixEncoder struct {
	zapcore.Encoder

	prefix  string
	bufPool buffer.Pool
}

func (e *prefixEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := e.bufPool.Get()

	buf.AppendString(e.prefix)
	buf.AppendString(" ")

	logEntry, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(logEntry.Bytes())
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func getDevEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return &prefixEncoder{
		Encoder: zapcore.NewConsoleEncoder(encoderConfig),
		prefix:  "[s7]",
		bufPool: buffer.NewPool(),
	}
}

func getProdEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return &prefixEncoder{
		Encoder: zapcore.NewConsoleEncoder(encoderConfig),
		prefix:  "[s7]",
		bufPool: buffer.NewPool(),
	}
}

// GetDefaultLogger returns the default logger.
func GetDefaultLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return defaultLogger
}

// GetDefaultFlusher returns the default flusher.
func GetDefaultFlusher() Flusher {
	mu.RLock()
	defer mu.RUnlock()
	return defaultFlusher
}

// SetDefaultLoggerAndFlusher sets the default logger and its flusher.
func SetDefaultLoggerAndFlusher(logger Logger, flusher Flusher) {
	mu.Lock()
	defaultLogger, defaultFlusher = logger, flusher
	mu.Unlock()
}

// LogLevel tells what the default logging level is.
func LogLevel() string {
	return strings.ToUpper(defaultLoggingLevel.String())
}

// CreateLoggerAsLocalFile setups the logger by local file path.
func CreateLoggerAsLocalFile(localFilePath string, logLevel Level) (logger Logger, flush func() error, err error) {
	if len(localFilePath) == 0 {
		return nil, nil, errors.New("invalid local logger path")
	}

	// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
	lumberJackLogger := &lumberjack.Logger{
		Filename:   localFilePath,
		MaxSize:    100, // megabytes
		MaxBackups: 2,
		MaxAge:     15, // days
	}

	encoder := getProdEncoder()
	ws := zapcore.AddSync(lumberJackLogger)
	zapcore.Lock(ws)

	levelEnabler := zap.LevelEnablerFunc(func(level Level) bool {
		return level >= logLevel
	})
	core := zapcore.NewCore(encoder, ws, levelEnabler)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(ErrorLevel))
	logger = zapLogger.Sugar()
	flush = zapLogger.Sync
	return
}

// Cleanup does something windup for logger, like closing, flushing, etc.
func Cleanup() {
	mu.RLock()
	if defaultFlusher != nil {
		_ = defaultFlusher()
	}
	mu.RUnlock()
}

// Error prints err if it's not nil.
func Error(err error) {
	if err != nil {
		mu.RLock()
		defaultLogger.Errorf("error occurs during runtime, %v", err)
		mu.RUnlock()
	}
}

// Debugf logs messages at DEBUG level.
func Debugf(format string, args ...interface{}) {
	mu.RLock()
	defaultLogger.Debugf(format, args...)
	mu.RUnlock()
}

// Infof logs messages at INFO level.
func Infof(format string, args ...interface{}) {
	mu.RLock()
	defaultLogger.Infof(format, args...)
	mu.RUnlock()
}

// Warnf logs messages at WARN level.
func Warnf(format string, args ...interface{}) {
	mu.RLock()
	defaultLogger.Warnf(format, args...)
	mu.RUnlock()
}

// Errorf logs messages at ERROR level.
func Errorf(format string, args ...interface{}) {
	mu.RLock()
	defaultLogger.Errorf(format, args...)
	mu.RUnlock()
}

// Fatalf logs messages at FATAL level.
func Fatalf(format string, args ...interface{}) {
	mu.RLock()
	defaultLogger.Fatalf(format, args...)
	mu.RUnlock()
}

// Logger is used for logging formatted messages.
type Logger interface {
	// Debugf logs messages at DEBUG level.
	Debugf(format string, args ...interface{})
	// Infof logs messages at INFO level.
	Infof(format string, args ...interface{})
	// Warnf logs messages at WARN level.
	Warnf(format string, args ...interface{})
	// Errorf logs messages at ERROR level.
	Errorf(format string, args ...interface{})
	// Fatalf logs messages at FATAL level.
	Fatalf(format string, args ...interface{})
}
