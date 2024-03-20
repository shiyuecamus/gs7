package _examples

import (
	"github.com/shiyuecamus/gs7"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/logging"
	"time"
)

func main() {
	const (
		host = "192.168.0.1"
		port = 102
		rack = 0
		slot = 1
	)
	logger := logging.GetDefaultLogger()

	c := gs7.NewClientBuilder().
		PlcType(common.S1500).
		Host(host).
		Port(port).
		Port(102).
		Rack(rack).
		Slot(slot).
		Build()

	_, err := c.Connect().Wait()
	if err != nil {
		logger.Errorf("Failed to connect PLC, host: %s, port: %d, error: %s", host, port, err)
		return
	}

	// Bit/X
	BitRwAndParse(c, logger)

	// Byte/B
	ByteRwAndParse(c, logger)

	// Char/C
	CharRwAndParse(c, logger)

	// Int/I
	IntRwAndParse(c, logger)

	// Word/W
	WordRwAndParse(c, logger)

	// DInt/DI
	DIntRwAndParse(c, logger)

	// DWord/DW
	DWordRwAndParse(c, logger)

	// Real/R
	RealRwAndParse(c, logger)

	// Time/T
	TimeRwAndParse(c, logger)

	// Date/D
	DateRwAndParse(c, logger)

	// TIMEOFDAY/TOD
	TodRwAndParse(c, logger)

	// DATETIME/DT
	DateTimeRwAndParse(c, logger)

	// DATETIMELONG/DTL
	DtlRwAndParse(c, logger)

	// STIME/ST
	S5TimeRwAndParse(c, logger)

	// STRING/S
	StringRwAndParse(c, logger)

	// WSTRING/WS
	WStringRwAndParse(c, logger)

	// C
	CounterRwAndParse(c, logger)

	// T
	TimerRwAndParse(c, logger)
}

func BitRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.X0.0").Wait()
	if err != nil {
		logger.Errorf("Failed to read bit, error: %s", err)
		return
	}

	b, err := gs7.BitFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to bit failed, error: %s", err)
		return
	}
	logger.Infof("read bit raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.BIT0.0", gs7.Bit(false).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.X0.0").Wait()
	if err != nil {
		logger.Errorf("Failed to read bit, error: %s", err)
		return
	}

	b, err = gs7.BitFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to bit failed, error: %s", err)
		return
	}
	logger.Infof("read bit raw success, bytes: %x, parsed value: %s", wait, b)
}

func ByteRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.B1").Wait()
	if err != nil {
		logger.Errorf("Failed to read byte, error: %s", err)
		return
	}

	b, err := gs7.ByteFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to byte failed, error: %s", err)
		return
	}
	logger.Infof("read byte raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.BYTE1", gs7.Byte(0x1C).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.B1").Wait()
	if err != nil {
		logger.Errorf("Failed to read byte, error: %s", err)
		return
	}

	b, err = gs7.ByteFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to byte failed, error: %s", err)
		return
	}
	logger.Infof("read byte raw success, bytes: %x, parsed value: %s", wait, b)
}

func CharRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.C2").Wait()
	if err != nil {
		logger.Errorf("Failed to read char, error: %s", err)
		return
	}

	b, err := gs7.CharFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to char failed, error: %s", err)
		return
	}
	logger.Infof("read char raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.CHAR2", gs7.Char('a').ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.C2").Wait()
	if err != nil {
		logger.Errorf("Failed to read char, error: %s", err)
		return
	}

	b, err = gs7.CharFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to char failed, error: %s", err)
		return
	}
	logger.Infof("read char raw success, bytes: %x, parsed value: %s", wait, b)
}

func IntRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.I4").Wait()
	if err != nil {
		logger.Errorf("Failed to read int, error: %s", err)
		return
	}

	b, err := gs7.IntFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to int failed, error: %s", err)
		return
	}
	logger.Infof("read int raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.INT4", gs7.Int(88).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.I4").Wait()
	if err != nil {
		logger.Errorf("Failed to read int, error: %s", err)
		return
	}

	b, err = gs7.IntFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to int failed, error: %s", err)
		return
	}
	logger.Infof("read int raw success, bytes: %x, parsed value: %s", wait, b)
}

func WordRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.W6").Wait()
	if err != nil {
		logger.Errorf("Failed to read word, error: %s", err)
		return
	}

	b, err := gs7.WordFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to word failed, error: %s", err)
		return
	}
	logger.Infof("read word raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.WORD6", gs7.Word(99).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.W6").Wait()
	if err != nil {
		logger.Errorf("Failed to read word, error: %s", err)
		return
	}

	b, err = gs7.WordFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to word failed, error: %s", err)
		return
	}
	logger.Infof("read word raw success, bytes: %x, parsed value: %s", wait, b)
}

func DIntRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.DI8").Wait()
	if err != nil {
		logger.Errorf("Failed to read dint, error: %s", err)
		return
	}

	b, err := gs7.DIntFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to dint failed, error: %s", err)
		return
	}
	logger.Infof("read dint raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.DINT8", gs7.DInt(188).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.DI8").Wait()
	if err != nil {
		logger.Errorf("Failed to read dint, error: %s", err)
		return
	}

	b, err = gs7.DIntFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to dint failed, error: %s", err)
		return
	}
	logger.Infof("read dint raw success, bytes: %x, parsed value: %s", wait, b)
}

func DWordRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.DW12").Wait()
	if err != nil {
		logger.Errorf("Failed to read dword, error: %s", err)
		return
	}

	b, err := gs7.DWordFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to dword failed, error: %s", err)
		return
	}
	logger.Infof("read dword raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.DWORD12", gs7.DWord(999).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.DW12").Wait()
	if err != nil {
		logger.Errorf("Failed to read dword, error: %s", err)
		return
	}

	b, err = gs7.DWordFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to dword failed, error: %s", err)
		return
	}
	logger.Infof("read dword raw success, bytes: %x, parsed value: %s", wait, b)
}

func RealRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.R16").Wait()
	if err != nil {
		logger.Errorf("Failed to read real, error: %s", err)
		return
	}

	b, err := gs7.RealFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to real failed, error: %s", err)
		return
	}
	logger.Infof("read real raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.REAL16", gs7.Real(355.538).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.R16").Wait()
	if err != nil {
		logger.Errorf("Failed to read real, error: %s", err)
		return
	}

	b, err = gs7.RealFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to real failed, error: %s", err)
		return
	}
	logger.Infof("read real raw success, bytes: %x, parsed value: %s", wait, b)
}

func TimeRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.T20").Wait()
	if err != nil {
		logger.Errorf("Failed to read time, error: %s", err)
		return
	}

	b, err := gs7.TimeFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to time failed, error: %s", err)
		return
	}
	logger.Infof("read time raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.TIME20", gs7.Time(time.Duration(355)*time.Millisecond).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.T20").Wait()
	if err != nil {
		logger.Errorf("Failed to read time, error: %s", err)
		return
	}

	b, err = gs7.TimeFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to time failed, error: %s", err)
		return
	}
	logger.Infof("read time raw success, bytes: %x, parsed value: %s", wait, b)
}

func DateRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.D24").Wait()
	if err != nil {
		logger.Errorf("Failed to read date, error: %s", err)
		return
	}

	b, err := gs7.DateFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to date failed, error: %s", err)
		return
	}
	logger.Infof("read date raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.DATE24", gs7.Date(time.Now()).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.D24").Wait()
	if err != nil {
		logger.Errorf("Failed to read date, error: %s", err)
		return
	}

	b, err = gs7.DateFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to date failed, error: %s", err)
		return
	}
	logger.Infof("read date raw success, bytes: %x, parsed value: %s", wait, b)
}

func TodRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.TOD26").Wait()
	if err != nil {
		logger.Errorf("Failed to read tod, error: %s", err)
		return
	}

	b, err := gs7.TimeOfDayFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to tod failed, error: %s", err)
		return
	}
	logger.Infof("read tod raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.TIMEOFDAY26", gs7.TimeOfDay(time.Now()).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.TOD26").Wait()
	if err != nil {
		logger.Errorf("Failed to read tod, error: %s", err)
		return
	}

	b, err = gs7.TimeOfDayFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to tod failed, error: %s", err)
		return
	}
	logger.Infof("read tod raw success, bytes: %x, parsed value: %s", wait, b)
}

func DateTimeRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.DT30").Wait()
	if err != nil {
		logger.Errorf("Failed to read datetime, error: %s", err)
		return
	}

	b, err := gs7.DateTimeFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to datetime failed, error: %s", err)
		return
	}
	logger.Infof("read datetime raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.DATETIME30", gs7.DateTime(time.Now()).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.DT30").Wait()
	if err != nil {
		logger.Errorf("Failed to read datetime, error: %s", err)
		return
	}

	b, err = gs7.DateTimeFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to datetime failed, error: %s", err)
		return
	}
	logger.Infof("read datetime raw success, bytes: %x, parsed value: %s", wait, b)
}

func DtlRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.DTL38").Wait()
	if err != nil {
		logger.Errorf("Failed to read dtl, error: %s", err)
		return
	}

	b, err := gs7.DateTimeLongFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to dtl failed, error: %s", err)
		return
	}
	logger.Infof("read dtl raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.DATETIMELONG38", gs7.DateTimeLong(time.Now()).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.DTL38").Wait()
	if err != nil {
		logger.Errorf("Failed to read dtl, error: %s", err)
		return
	}

	b, err = gs7.DateTimeLongFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to dtl failed, error: %s", err)
		return
	}
	logger.Infof("read dtl raw success, bytes: %x, parsed value: %s", wait, b)
}

func S5TimeRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.ST50").Wait()
	if err != nil {
		logger.Errorf("Failed to read s5time, error: %s", err)
		return
	}

	b, err := gs7.S5TimeFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to s5time failed, error: %s", err)
		return
	}
	logger.Infof("read s5time raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.STIME50", gs7.S5Time(time.Duration(40)*time.Millisecond).ToBytes()).Wait()

	wait, err = c.ReadRaw("DB1.ST50").Wait()
	if err != nil {
		logger.Errorf("Failed to read s5time, error: %s", err)
		return
	}

	b, err = gs7.S5TimeFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to s5time failed, error: %s", err)
		return
	}
	logger.Infof("read s5time raw success, bytes: %x, parsed value: %s", wait, b)
}

func StringRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.S52").Wait()
	if err != nil {
		logger.Errorf("Failed to read string, error: %s", err)
		return
	}

	b, err := gs7.StringFromBytes(wait, c.GetPlcType())
	if err != nil {
		logger.Errorf("Parse raw to string failed, error: %s", err)
		return
	}
	logger.Infof("read string raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.STRING52", gs7.String("test string!!").ToBytes(c.GetPduLength())).Wait()

	wait, err = c.ReadRaw("DB1.S52").Wait()
	if err != nil {
		logger.Errorf("Failed to read string, error: %s", err)
		return
	}

	b, err = gs7.StringFromBytes(wait, c.GetPlcType())
	if err != nil {
		logger.Errorf("Parse raw to string failed, error: %s", err)
		return
	}
	logger.Infof("read string raw success, bytes: %x, parsed value: %s", wait, b)
}

func WStringRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("DB1.WS308").Wait()
	if err != nil {
		logger.Errorf("Failed to read wstring, error: %s", err)
		return
	}

	b, err := gs7.WStringFromBytes(wait, c.GetPlcType())
	if err != nil {
		logger.Errorf("Parse raw to wstring failed, error: %s", err)
		return
	}
	logger.Infof("read wstring raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("DB1.WSTRING308", gs7.WString("测试一下中文").ToBytes(c.GetPduLength())).Wait()

	wait, err = c.ReadRaw("DB1.WS308").Wait()
	if err != nil {
		logger.Errorf("Failed to read wstring, error: %s", err)
		return
	}

	b, err = gs7.WStringFromBytes(wait, c.GetPlcType())
	if err != nil {
		logger.Errorf("Parse raw to wstring failed, error: %s", err)
		return
	}
	logger.Infof("read wstring raw success, bytes: %x, parsed value: %s", wait, b)
}

func CounterRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("C0").Wait()
	if err != nil {
		logger.Errorf("Failed to read counter, error: %s", err)
		return
	}

	b, err := gs7.CounterFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to counter failed, error: %s", err)
		return
	}
	logger.Infof("read counter raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("C0", gs7.Counter(126).ToBytes()).Wait()

	wait, err = c.ReadRaw("C0").Wait()
	if err != nil {
		logger.Errorf("Failed to read counter, error: %s", err)
		return
	}

	b, err = gs7.CounterFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to counter failed, error: %s", err)
		return
	}
	logger.Infof("read counter raw success, bytes: %x, parsed value: %s", wait, b)
}

func TimerRwAndParse(c gs7.Client, logger logging.Logger) {
	wait, err := c.ReadRaw("T0").Wait()
	if err != nil {
		logger.Errorf("Failed to read timer, error: %s", err)
		return
	}

	b, err := gs7.TimerFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to timer failed, error: %s", err)
		return
	}
	logger.Infof("read timer raw success, bytes: %x, parsed value: %s", wait, b)

	_ = c.WriteRaw("T0", gs7.Timer(time.Duration(123)*time.Millisecond).ToBytes()).Wait()

	wait, err = c.ReadRaw("T0").Wait()
	if err != nil {
		logger.Errorf("Failed to read timer, error: %s", err)
		return
	}

	b, err = gs7.TimerFromBytes(wait)
	if err != nil {
		logger.Errorf("Parse raw to timer failed, error: %s", err)
		return
	}
	logger.Infof("read timer raw success, bytes: %x, parsed value: %s", wait, b)
}
