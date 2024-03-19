// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import "gs7/model"

type RequestContext interface {
	PutError(err error)
	PutResponse(data *model.PDU)
	GetResponse() (data *model.PDU, err error)
	GetRequestId() uint16
	GetRequest() *model.PDU
}

type StandardRequestContext struct {
	RequestId uint16
	Request   *model.PDU
	Response  chan *model.PDU
	Error     chan error
}

func (s *StandardRequestContext) PutError(err error) {
	s.Error <- err
}

func (s *StandardRequestContext) PutResponse(data *model.PDU) {
	s.Response <- data
}

func (s *StandardRequestContext) GetRequestId() uint16 {
	return s.RequestId
}

func (s *StandardRequestContext) GetRequest() *model.PDU {
	return s.Request
}

func (s *StandardRequestContext) GetResponse() (res *model.PDU, err error) {
	for {
		select {
		case res = <-s.Response:
			return
		case err = <-s.Error:
			return
		}
	}
}

type ConnectRequestContext struct {
	Request  *model.PDU
	Response chan *model.PDU
	Error    chan error
}

func (c *ConnectRequestContext) PutError(err error) {
	c.Error <- err
}

func (c *ConnectRequestContext) PutResponse(data *model.PDU) {
	c.Response <- data
}

func (c *ConnectRequestContext) GetRequestId() uint16 {
	return 0
}

func (c *ConnectRequestContext) GetRequest() *model.PDU {
	return c.Request
}

func (c *ConnectRequestContext) GetResponse() (res *model.PDU, err error) {
	for {
		select {
		case res = <-c.Response:
			return
		case err = <-c.Error:
			return
		}
	}
}
