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

	gs7.NewClientBuilder().
		PlcType(common.S1500).
		Host(host).
		Port(port).
		Port(102).
		Rack(rack).
		Slot(slot).
		BuildAndConnect().
		Async(func(c gs7.Client, err error) {
			if err != nil {
				logger.Errorf("Failed to connect PLC, host: %s, port: %d, error: %s", host, port, err)
				return
			}

			c.ReadParsed("DB1.X260.0").Async(func(v any, err error) {
				if err != nil {
					logger.Errorf("Failed to read bit, error: %s", err)
					return
				}
				logger.Infof("Read bit success with value: %s", v)
			})

			c.WriteRaw("DB1.X260.0", gs7.Bit(true).ToBytes()).Async(func(err error) {
				if err != nil {
					logger.Errorf("Failed to read bit, error: %s", err)
					return
				}
			})

			c.ReadParsed("DB1.X260.0").Async(func(v any, err error) {
				if err != nil {
					logger.Errorf("Failed to read bit, error: %s", err)
					return
				}
				logger.Infof("Read bit success with value: %s", v)
			})
		})

	select {}
}
