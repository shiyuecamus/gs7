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
		Port(102).
		Rack(rack).
		Slot(slot).
		Build()

	if _, err := c.Connect().Wait(); err != nil {
		logger.Errorf("Failed to connect PLC, host: %s, port: %d, error: %s", host, port, err)
		return
	}
	defer c.Disconnect()

	unitInfo, _ := c.GetUnitInfo().Wait()
	protectionInfo, _ := c.GetProtectionInfo().Wait()
	plcStatus, _ := c.GetPlcStatus().Wait()
	catalog, _ := c.GetCatalog().Wait()
	communicationInfo, _ := c.GetCommunicationInfo().Wait()
	logging.Infof("%s", unitInfo)
	logging.Infof("%s", protectionInfo)
	logging.Infof("%s", plcStatus)
	logging.Infof("%s", catalog)
	logging.Infof("%s", communicationInfo)
}
