// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import (
	"errors"
	"sync"
)

// Status - Manage the connection status

// Multiple go routines will want to access/set this. Previously status was implemented as a `uint32` and updated
// with a mixture of atomic functions and a mutex (leading to some deadlock type issues that were very hard to debug).

// In this new implementation `connectionStatus` takes over managing the state and provides functions that allow the
// client to request a move to a particular state (it may reject these requests!). In some cases the 'state' is
// transitory, for example `connecting`, in those cases a function will be returned that allows the client to move
// to a more static state (`disconnected` or `connected`).

// This "belts-and-braces" may be a little over the top but issues with the status have caused a number of difficult
// to trace bugs in the past and the likelihood that introducing a new system would introduce bugs seemed high!
// I have written this in a way that should make it very difficult to misuse it (but it does make things a little
// complex with functions returning functions that return functions!).

type Status uint32

const (
	disconnected Status = iota
	disconnecting
	connecting
	reconnecting
	connected
)

func (s Status) String() string {
	switch s {
	case disconnected:
		return "disconnected"
	case disconnecting:
		return "disconnecting"
	case connecting:
		return "connecting"
	case reconnecting:
		return "reconnecting"
	case connected:
		return "connected"
	default:
		return "invalid"
	}
}

type connCompletedFn func(success bool) error
type disconnectCompletedFn func()
type connectionLostHandledFn func(bool) (connCompletedFn, error)

var (
	errAbortConnection                = errors.New("disconnect called whist connection attempt in progress")
	errAlreadyConnectedOrReconnecting = errors.New("status is already connected or reconnecting")
	errStatusMustBeDisconnected       = errors.New("status can only transition to connecting from disconnected")
	errAlreadyDisconnected            = errors.New("status is already disconnected")
	errDisconnectionRequested         = errors.New("disconnection was requested whilst the action was in progress")
	errDisconnectionInProgress        = errors.New("disconnection already in progress")
)

type connectionStatus struct {
	sync.RWMutex
	status        Status
	willReconnect bool

	actionCompleted chan struct{}
}

func (c *connectionStatus) ConnectionStatus() Status {
	c.RLock()
	defer c.RUnlock()
	return c.status
}

func (c *connectionStatus) ConnectionStatusRetry() (Status, bool) {
	c.RLock()
	defer c.RUnlock()
	return c.status, c.willReconnect
}

func (c *connectionStatus) Connecting() (connCompletedFn, error) {
	c.Lock()
	defer c.Unlock()
	if c.status == connected || c.status == reconnecting {
		return nil, errAlreadyConnectedOrReconnecting
	}
	if c.status != disconnected {
		return nil, errStatusMustBeDisconnected
	}
	c.status = connecting
	c.actionCompleted = make(chan struct{})
	return c.connected, nil
}

func (c *connectionStatus) connected(success bool) error {
	c.Lock()
	defer func() {
		close(c.actionCompleted)
		c.actionCompleted = nil
		c.Unlock()
	}()

	if c.status == disconnecting {
		return errAbortConnection
	}
	if success {
		c.status = connected
	} else {
		c.status = disconnected
	}
	return nil
}

func (c *connectionStatus) Disconnecting() (disconnectCompletedFn, error) {
	c.Lock()
	if c.status == disconnected {
		c.Unlock()
		return nil, errAlreadyDisconnected
	}
	if c.status == disconnecting {
		c.willReconnect = false
		disConnectDone := c.actionCompleted
		c.Unlock()
		<-disConnectDone
		return nil, errAlreadyDisconnected
	}

	prevStatus := c.status
	c.status = disconnecting

	if prevStatus == connecting || prevStatus == reconnecting {
		connectDone := c.actionCompleted
		c.Unlock()
		<-connectDone

		if prevStatus == reconnecting && !c.willReconnect {
			return nil, errAlreadyDisconnected // Following connectionLost process we will be disconnected
		}
		c.Lock()
	}
	c.actionCompleted = make(chan struct{})
	c.Unlock()
	return c.disconnectionCompleted, nil
}

func (c *connectionStatus) disconnectionCompleted() {
	c.Lock()
	defer c.Unlock()
	c.status = disconnected
	close(c.actionCompleted)
	c.actionCompleted = nil
}

func (c *connectionStatus) ConnectionLost(willReconnect bool) (connectionLostHandledFn, error) {
	c.Lock()
	defer c.Unlock()
	if c.status == disconnected {
		return nil, errAlreadyDisconnected
	}
	if c.status == disconnecting { // its expected that connection lost will be called during the disconnection process
		return nil, errDisconnectionInProgress
	}

	c.willReconnect = willReconnect
	prevStatus := c.status
	c.status = disconnecting

	if prevStatus == connecting || prevStatus == reconnecting {
		connectDone := c.actionCompleted
		c.Unlock()
		<-connectDone
		c.Lock()
		if !willReconnect {
			// In this case the connection will always be aborted so there is nothing more for us to do
			return nil, errAlreadyDisconnected
		}
	}
	c.actionCompleted = make(chan struct{})

	return c.getConnectionLostHandler(willReconnect), nil
}

func (c *connectionStatus) getConnectionLostHandler(reconnectRequested bool) connectionLostHandledFn {
	return func(proceed bool) (connCompletedFn, error) {
		c.Lock()
		defer c.Unlock()

		if !c.willReconnect || !proceed {
			c.status = disconnected
			close(c.actionCompleted)
			c.actionCompleted = nil
			if !reconnectRequested || !proceed {
				return nil, nil
			}
			return nil, errDisconnectionRequested
		}

		c.status = reconnecting
		return c.connected, nil
	}
}

func (c *connectionStatus) forceConnectionStatus(s Status) {
	c.Lock()
	defer c.Unlock()
	c.status = s
}
