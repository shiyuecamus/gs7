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
		// connect and read Timeout
		// default value 5s
		Timeout(time.Duration(5) * time.Second). // connect and read timeout default: 5s
		// connectRetry automatically retry when attempting to connect failed
		// default value false
		ConnectRetry(true).
		// maximum connection retry interval.
		// each attempt will be multiplied by 2
		// default value 10s
		RetryInterval(time.Duration(5) * time.Second).
		// maximum connection wait time
		// default value 5min
		MaxRetryBackoff(time.Duration(1) * time.Minute).
		// reconnect on connection lost
		// default value false
		AutoReconnect(true).
		// maximum connection reconnect interval.
		// each attempt will be multiplied by 2
		// default value 10s
		ReconnectInterval(time.Duration(5) * time.Second).
		// maximum reconnect times
		// if set to -1, infinite attempts will be made to reconnect
		// default value 5
		MaxReconnectBackoff(time.Duration(1) * time.Minute).
		Build()

	if _, err := c.Connect().Wait(); err != nil {
		logger.Errorf("Failed to connect PLC, host: %s, port: %d, error: %s", host, port, err)
		return
	}
	defer c.Disconnect()
}
