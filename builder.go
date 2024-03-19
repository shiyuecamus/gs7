// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import (
	"gs7/common"
	"gs7/logging"
	"gs7/util"
	"sync"
	"time"
)

type ClientBuilder struct {
	logger    logging.Logger
	plcType   common.PlcType
	host      string
	port      int
	rack      int
	slot      int
	pduLength int
	// timeout connect and read Timeout
	// default value 5s
	timeout time.Duration
	// autoReconnect reconnect on connection lost
	// default value false
	autoReconnect bool
	// reconnectInterval maximum connection reconnect interval.
	// each attempt will be multiplied by 2
	// default value 10s
	reconnectInterval time.Duration
	// maxReconnectTimes maximum reconnect times
	// if set to -1, infinite attempts will be made to reconnect
	// default value 5min
	maxReconnectTimes int
	// maxReconnectBackoff maximum reconnect wait time
	// default value 5min
	maxReconnectBackoff time.Duration
	// connectRetry automatically retry when attempting to connect failed
	// default value false
	connectRetry bool
	// retryInterval maximum connection retry interval.
	// each attempt will be multiplied by 2
	// default value 10s
	retryInterval time.Duration
	// maxRetries maximum connection retries
	// if set to -1, infinite attempts will be made to retry
	// default value 5
	maxRetries int
	// maxRetryBackoff maximum connection wait time
	// default value 5min
	maxRetryBackoff time.Duration
	onConnected     func(c Client)
	onUnActive      func(c Client, err error)
}

func NewClientBuilder() ClientBuilder {
	return ClientBuilder{}
}

func (b ClientBuilder) PlcType(plcType common.PlcType) ClientBuilder {
	b.plcType = plcType
	return b
}

func (b ClientBuilder) Host(host string) ClientBuilder {
	b.host = host
	return b
}

func (b ClientBuilder) Port(port int) ClientBuilder {
	b.port = port
	return b
}

func (b ClientBuilder) Rack(rack int) ClientBuilder {
	b.rack = rack
	return b
}

func (b ClientBuilder) Slot(slot int) ClientBuilder {
	b.slot = slot
	return b
}

func (b ClientBuilder) PduLength(pduLength int) ClientBuilder {
	b.pduLength = pduLength
	return b
}

func (b ClientBuilder) Timeout(timeout time.Duration) ClientBuilder {
	b.timeout = timeout
	return b
}

func (b ClientBuilder) AutoReconnect(autoReconnect bool) ClientBuilder {
	b.autoReconnect = autoReconnect
	return b
}

func (b ClientBuilder) ReconnectInterval(interval time.Duration) ClientBuilder {
	b.reconnectInterval = interval
	return b
}

func (b ClientBuilder) MaxReconnectTimes(maxReconnectTimes int) ClientBuilder {
	b.maxReconnectTimes = maxReconnectTimes
	return b
}

func (b ClientBuilder) MaxReconnectBackoff(maxReconnectBackoff time.Duration) ClientBuilder {
	b.maxReconnectBackoff = maxReconnectBackoff
	return b
}

func (b ClientBuilder) ConnectRetry(connectRetry bool) ClientBuilder {
	b.connectRetry = connectRetry
	return b
}

func (b ClientBuilder) RetryInterval(interval time.Duration) ClientBuilder {
	b.retryInterval = interval
	return b
}

func (b ClientBuilder) MaxRetries(maxRetries int) ClientBuilder {
	b.maxRetries = maxRetries
	return b
}

func (b ClientBuilder) MaxRetryBackoff(maxRetryBackoff time.Duration) ClientBuilder {
	b.maxRetryBackoff = maxRetryBackoff
	return b
}

func (b ClientBuilder) Logger(logger logging.Logger) ClientBuilder {
	b.logger = logger
	return b
}

func (b ClientBuilder) OnConnected(onConnected func(c Client)) ClientBuilder {
	b.onConnected = onConnected
	return b
}

func (b ClientBuilder) OnUnActive(onUnActive func(c Client, err error)) ClientBuilder {
	b.onUnActive = onUnActive
	return b
}

const (
	DefaultPduLength        = 480
	Localhost        string = "127.0.0.1"
	DefaultPort      int    = 102
)

func (b ClientBuilder) Build() Client {
	s := &client{
		m:                   new(sync.RWMutex),
		plcType:             b.plcType,
		host:                util.StrOrDefault(b.host, Localhost),
		port:                util.IntOrDefault(b.port, DefaultPort),
		rack:                b.rack,
		slot:                b.slot,
		pduLength:           util.IntOrDefault(b.pduLength, DefaultPduLength),
		timeout:             util.DurationOrDefault(b.timeout, time.Duration(5)*time.Second),
		autoReconnect:       b.autoReconnect,
		reconnectInterval:   util.DurationOrDefault(b.reconnectInterval, time.Duration(10)*time.Second),
		maxReconnectTimes:   util.IntOrDefault(b.maxReconnectTimes, 5),
		maxReconnectBackoff: util.DurationOrDefault(b.maxReconnectBackoff, time.Duration(5)*time.Minute),
		connectRetry:        b.connectRetry,
		retryInterval:       util.DurationOrDefault(b.retryInterval, time.Duration(10)*time.Second),
		maxRetries:          util.IntOrDefault(b.maxRetries, 5),
		maxRetryBackoff:     util.DurationOrDefault(b.maxRetryBackoff, time.Duration(5)*time.Minute),
		logger:              util.AnyOrDefault(b.logger, logging.GetDefaultLogger()).(logging.Logger),
		onConnected:         b.onConnected,
		onDisconnected:      b.onUnActive,
	}
	return s.init()
}

func (b ClientBuilder) BuildAndConnect() *ConnectToken {
	c := b.Build()
	return c.Connect()
}
