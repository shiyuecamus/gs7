// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import (
	"github.com/panjf2000/gnet/v2"
	"gs7/common"
	"gs7/logging"
	"gs7/model"
	"gs7/util"
	"sync"
	"time"
)

const (
	isoHeaderSize = common.TpktLen + common.CotpDataLen // TPKT+COTP Header Size
	minPduSize    = 1
)

type s7TcpClient struct {
	*gnet.BuiltinEventEngine
	logger            logging.Logger
	eng               gnet.Engine
	requestContextMap sync.Map
	isoConnectChan    chan RequestContext
	isoDisconnectChan chan RequestContext
	onOpen            OnOpen
	onClose           OnClose
	validate          PduValidate
	// Connect
	timeout time.Duration
}

type OnOpen func(c gnet.Conn)
type OnClose func(c gnet.Conn, err error)
type PduValidate func(tpkt common.TPKT) error

func newTcpClient(logger logging.Logger, timeout time.Duration,
	onOpen OnOpen, onClose OnClose, pduValidate PduValidate) *s7TcpClient {
	return &s7TcpClient{
		logger:            logger,
		isoConnectChan:    make(chan RequestContext, 1),
		isoDisconnectChan: make(chan RequestContext, 1),
		onOpen:            onOpen,
		onClose:           onClose,
		validate:          pduValidate,
		timeout:           timeout,
	}
}

func (t *s7TcpClient) handleRequestContext(context RequestContext) (err error) {
	_, ok := t.requestContextMap.Load(context.GetRequestId())
	if ok {
		err = common.ErrorWithCode(common.ErrTcpRequestProcessing, context.GetRequestId())
		return
	}
	t.requestContextMap.Store(context.GetRequestId(), context)
	time.AfterFunc(t.timeout, func() {
		ctx, ok := t.requestContextMap.LoadAndDelete(context.GetRequestId())
		if ok {
			ctx.(RequestContext).PutError(common.ErrorWithCode(common.ErrTcpRequestTimeout))
		}
	})
	return
}

func (t *s7TcpClient) handleConnectRequestContext(context RequestContext) (err error) {
	select {
	case t.isoConnectChan <- context:
		break
	default:
		err = common.ErrorWithCode(common.ErrTcpRequestProcessing, context.GetRequestId())
		return
	}
	time.AfterFunc(t.timeout, func() {
		select {
		case reqCtx := <-t.isoConnectChan:
			reqCtx.PutError(common.ErrorWithCode(common.ErrTcpRequestTimeout))
			break
		default:
			break
		}
	})
	return
}

func (t *s7TcpClient) handleDisconnectRequestContext(context RequestContext) (err error) {
	select {
	case t.isoDisconnectChan <- context:
		break
	default:
		err = common.ErrorWithCode(common.ErrTcpRequestProcessing, context.GetRequestId())
		return
	}
	time.AfterFunc(t.timeout, func() {
		select {
		case reqCtx := <-t.isoDisconnectChan:
			reqCtx.PutError(common.ErrorWithCode(common.ErrTcpRequestTimeout))
			break
		default:
			break
		}
	})
	return
}

func (t *s7TcpClient) OnBoot(eng gnet.Engine) (action gnet.Action) {
	t.logger.Infof("S7 tcp client on boot.")
	t.eng = eng
	return
}

func (t *s7TcpClient) OnShutdown(gnet.Engine) {
	t.logger.Infof("S7 tcp client on shutdown.")
}

func (t *s7TcpClient) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	t.logger.Infof("S7 tcp client connection [%s] did open", c.RemoteAddr().String())
	util.Invoke(t.onOpen, []interface{}{(*gnet.Conn)(nil)}, c)
	return
}

func (t *s7TcpClient) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	t.logger.Infof("S7 tcp client connection [%s] did closed with error: %v", c.RemoteAddr().String(), err)
	util.Invoke(t.onClose, []interface{}{(*gnet.Conn)(nil), (*error)(nil)}, c, err)
	return
}

func (t *s7TcpClient) OnTraffic(c gnet.Conn) (action gnet.Action) {
	tpktBuf, err := c.Next(common.TpktLen)
	if tpktBuf == nil || err != nil {
		t.logger.Warnf("S7 tcp client received invalid package")
		return
	}
	tpkt, err := model.TPKTFromBytes(tpktBuf)
	if err != nil {
		t.logger.Warnf("S7 tcp client parse tpkt failed with error: [%v]", err)
		return
	}

	nextBuf, err := c.Next(int(tpkt.GetLength()) - common.TpktLen)
	if nextBuf == nil || err != nil {
		t.logger.Warnf("S7 tcp client parse package main failed with error: [%v]", err)
		return
	}
	total := append(tpktBuf, nextBuf...)
	t.logger.Debugf("S7 client received: % x", total)
	ack, err := model.DataFromBytes(total)
	var ctx RequestContext

	switch ack.GetCOTP().GetPduType() {
	case common.PtConnectConfirm:
		select {
		case ctx = <-t.isoConnectChan:
			break
		default:
			t.logger.Infof("S7 tcp client discard timeout connect response")
			return
		}
	case common.PtReject:
		select {
		case ctx = <-t.isoConnectChan:
			ctx.PutError(common.ErrorWithCode(common.ErrTcpRequestRejected))
			return
		default:
			t.logger.Infof("S7 tcp client discard timeout connect response")
			return
		}
	case common.PtDisconnectConfirm:
		select {
		case ctx = <-t.isoDisconnectChan:
			break
		default:
			t.logger.Infof("S7 tcp client discard timeout disconnect response")
			return
		}
	default:
		if ack.GetHeader() != nil {
			if value, ok := t.requestContextMap.LoadAndDelete(ack.GetHeader().GetPduReference()); ok {
				ctx = value.(RequestContext)
			}
		}
	}
	if ctx != nil {
		if err != nil {
			ctx.PutError(err)
			return
		}
		if t.validate != nil {
			if err = t.validate(tpkt); err != nil {
				ctx.PutError(err)
				return
			}
		}
		ctx.PutResponse(ack)
	}
	return
}
