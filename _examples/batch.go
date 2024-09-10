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

	if _, err := c.Connect().Wait(); err != nil {
		logger.Errorf("Failed to connect PLC, host: %s, port: %d, error: %s", host, port, err)
		return
	}
	defer c.Disconnect()

	BatchRead(c, logger)
	BatchWrite(c, logger)
	BatchRead(c, logger)
}

var addresses = []string{
	"DB1.X0.0",
	"DB1.B1",
	"DB1.C2",
	"DB1.I4",
	"DB1.W6",
	"DB1.DI8",
	"DB1.DW12",
	"DB1.R16",
	"DB1.T20",
	"DB1.D24",
	"DB1.TOD26",
	"DB1.DT30",
	"DB1.DTL38",
	"DB1.ST50",
	"DB1.S52",
	"DB1.WS308",
}

func BatchRead(c gs7.Client, logger logging.Logger) {
	v, err := c.ReadBatchRaw(addresses).Wait()
	if err != nil {
		logger.Errorf("Failed to batch read, error: %s", err)
		return
	}
	for _, bytes := range v {
		logger.Infof("%x", bytes)
	}
}

func BatchWrite(c gs7.Client, logger logging.Logger) {
	data := make([][]byte, 0, len(addresses))
	data = append(data, gs7.Bit(true).ToBytes())
	data = append(data, gs7.Byte(0x16).ToBytes())
	data = append(data, gs7.Char('t').ToBytes())
	data = append(data, gs7.Int(-88).ToBytes())
	data = append(data, gs7.Word(88).ToBytes())
	data = append(data, gs7.DInt(-188).ToBytes())
	data = append(data, gs7.DWord(188).ToBytes())
	data = append(data, gs7.Real(88.88).ToBytes())
	data = append(data, gs7.Time(time.Duration(88)*time.Millisecond).ToBytes())
	data = append(data, gs7.Date(time.Now()).ToBytes())
	data = append(data, gs7.TimeOfDay(time.Now()).ToBytes())
	data = append(data, gs7.DateTime(time.Now()).ToBytes())
	data = append(data, gs7.DateTimeLong(time.Now()).ToBytes())
	data = append(data, gs7.S5Time(time.Duration(80)*time.Millisecond).ToBytes())
	data = append(data, gs7.String("batch read").ToBytes(c.GetPduLength()))
	data = append(data, gs7.WString("批处理写入").ToBytes(c.GetPduLength()))
	err := c.WriteRawBatch(addresses, data).Wait()
	if err != nil {
		logger.Errorf("Failed to batch read, error: %s", err)
		return
	}
}
