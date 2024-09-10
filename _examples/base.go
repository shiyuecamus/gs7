package _examples

import (
	"github.com/shiyuecamus/gs7"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/logging"
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
		Rack(rack).
		Slot(slot).
		Build()
	if _, err := c.Connect().Wait(); err != nil {
		logger.Errorf("Failed to connect PLC, host: %s, port: %d, error: %s", host, port, err)
		return
	}
	defer c.Disconnect()

	// read
	wait, err := c.ReadParsed("DB1.X0.0").Wait()
	if err != nil {
		logger.Errorf("Failed to read bit, error: %s", err)
		return
	}
	logger.Infof("Read bit success with value: %s", wait)

	// write
	err = c.WriteRaw("DB1.X0.0", gs7.Bit(true).ToBytes()).Wait()
	if err != nil {
		logger.Errorf("Failed to read bit, error: %s", err)
		return
	}

	// check
	wait, err = c.ReadParsed("DB1.X0.0").Wait()
	if err != nil {
		logger.Errorf("Failed to read bit, error: %s", err)
		return
	}
	logger.Infof("Read bit success with value: %s", wait)
}
