// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import (
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/core"
	"github.com/shiyuecamus/gs7/util"
	"sync"
	"time"
)

type ErrorToken interface {
	// Wait will wait indefinitely for the ErrorToken to complete, ie the Publishing
	// to be sent and confirmed receipt from the broker.
	Wait() error

	// WaitTimeout takes a time.Duration to wait for the flow associated with the
	// ErrorToken to complete, returns true if it returned before the timeout or
	// returns false if the timeout occurred. In the case of a timeout the ErrorToken
	// does not have an error set in case the caller wishes to wait again.
	WaitTimeout(time.Duration) error

	Async(func(error))

	// Done returns a channel that is closed when the flow associated
	// with the ErrorToken completes. Clients should call Error after the
	// channel is closed to check if the flow completed successfully.
	//
	// Done is provided for use in select statements. Simple use cases may
	// use Wait or WaitTimeout.
	Done() <-chan struct{}
}

type ValueToken[V interface{}] interface {
	// Wait will wait indefinitely for the ValueToken to complete, ie the Publishing
	// to be sent and confirmed receipt from the broker.
	Wait() (V, error)

	// WaitTimeout takes a time.Duration to wait for the flow associated with the
	// ValueToken to complete, returns true if it returned before the timeout or
	// returns false if the timeout occurred. In the case of a timeout the ValueToken
	// does not have an error set in case the caller wishes to wait again.
	WaitTimeout(time.Duration) (V, error)

	Async(func(V, error))

	// Done returns a channel that is closed when the flow associated
	// with the ValueToken completes. Clients should call Error after the
	// channel is closed to check if the flow completed successfully.
	//
	// Done is provided for use in select statements. Simple use cases may
	// use Wait or WaitTimeout.
	Done() <-chan struct{}
}

type TokenErrorSetter interface {
	setError(error)
}

type TokenCompleter interface {
	TokenErrorSetter
	flowComplete()
}

type TokenType byte

const (
	TtSimple TokenType = iota
	TtConnect
	TtRead
	TtPdu
	TtSingleRawRead
	TtBatchRawRead
	TtSingleParsedRead
	TtBatchParsedRead
	TtUpload
	TtSzlIds
	TtCatalog
	TtPlcStatus
	TtUnitInfo
	TtCommunicationInfo
	TtProtectionInfo
	TtBlockList
	TtBlockListType
	TtBlockInfo
	TtClockRead
	TtBaseRead
)

func NewToken(tt TokenType) TokenCompleter {
	switch tt {
	case TtConnect:
		return &ConnectToken{baseToken[Client]{complete: make(chan struct{})}}
	case TtSimple:
		return &SimpleToken{complete: make(chan struct{})}
	case TtRead:
		return &ReadToken{baseToken[[]*core.DataItem]{complete: make(chan struct{})}}
	case TtPdu:
		return &PduToken{baseToken[*core.PDU]{complete: make(chan struct{})}}
	case TtSingleRawRead:
		return &SingleRawReadToken{baseToken[RawInfo]{complete: make(chan struct{})}}
	case TtBatchRawRead:
		return &BatchRawReadToken{baseToken[[]RawInfo]{complete: make(chan struct{})}}
	case TtSingleParsedRead:
		return &SingleParsedReadToken{baseToken[any]{complete: make(chan struct{})}}
	case TtBatchParsedRead:
		return &BatchParsedReadToken{baseToken[[]any]{complete: make(chan struct{})}}
	case TtUpload:
		return &UploadToken{baseToken[[]byte]{complete: make(chan struct{})}}
	case TtSzlIds:
		return &SzlIdsToken{baseToken[[]uint16]{complete: make(chan struct{})}}
	case TtCatalog:
		return &CatalogToken{baseToken[core.Catalog]{complete: make(chan struct{})}}
	case TtPlcStatus:
		return &PlcStatusToken{baseToken[core.PlcStatus]{complete: make(chan struct{})}}
	case TtUnitInfo:
		return &UnitInfoToken{baseToken[core.UnitInfo]{complete: make(chan struct{})}}
	case TtCommunicationInfo:
		return &CommunicationInfoToken{baseToken[core.CommunicationInfo]{complete: make(chan struct{})}}
	case TtProtectionInfo:
		return &ProtectionInfoToken{baseToken[core.ProtectionInfo]{complete: make(chan struct{})}}
	case TtBlockList:
		return &BlockListToken{baseToken[[]core.ListBlockInfo]{complete: make(chan struct{})}}
	case TtBlockListType:
		return &BlockListTypeToken{baseToken[[]core.ListBlockTypeInfo]{complete: make(chan struct{})}}
	case TtBlockInfo:
		return &BlockInfoToken{baseToken[core.BlockInfo]{complete: make(chan struct{})}}
	case TtClockRead:
		return &ClockReadToken{baseToken[time.Time]{complete: make(chan struct{})}}
	case TtBaseRead:
		return &BaseReadToken{baseToken[[]byte]{complete: make(chan struct{})}}
	default:
		return nil
	}
}

type SimpleToken struct {
	m        sync.RWMutex
	complete chan struct{}
	err      error
}

func (s *SimpleToken) Wait() error {
	<-s.complete
	return s.err
}

func (s *SimpleToken) WaitTimeout(d time.Duration) error {
	timer := time.NewTimer(d)
	select {
	case <-s.complete:
		if !timer.Stop() {
			<-timer.C
		}
		return s.err
	case <-timer.C:
	}

	return common.ErrorWithCode(common.ErrTcpRequestTimeout)
}

func (s *SimpleToken) Async(ack func(err error)) {
	go func() {
		<-s.complete
		util.Invoke(ack, []interface{}{(*error)(nil)}, s.err)
	}()
}

func (s *SimpleToken) Done() <-chan struct{} {
	return s.complete
}

func (s *SimpleToken) setError(e error) {
	s.m.Lock()
	s.err = e
	s.flowComplete()
	s.m.Unlock()
}

func (s *SimpleToken) flowComplete() {
	select {
	case <-s.complete:
	default:
		close(s.complete)
	}
}

type baseToken[V interface{}] struct {
	m        sync.RWMutex
	complete chan struct{}
	err      error
	v        V
}

func (b *baseToken[V]) Wait() (V, error) {
	<-b.complete
	return b.v, b.err
}

// WaitTimeout implements the ValueToken WaitTimeout method.
func (b *baseToken[V]) WaitTimeout(d time.Duration) (V, error) {
	timer := time.NewTimer(d)
	select {
	case <-b.complete:
		if !timer.Stop() {
			<-timer.C
		}
		return b.v, b.err
	case <-timer.C:
	}

	return b.v, common.ErrorWithCode(common.ErrTcpRequestTimeout)
}

func (b *baseToken[V]) Async(ack func(v V, err error)) {
	go func() {
		<-b.complete
		util.Invoke(ack, []interface{}{(*V)(nil), (*error)(nil)}, b.v, b.err)
	}()
}

// Done implements the ValueToken Done method.
func (b *baseToken[V]) Done() <-chan struct{} {
	return b.complete
}

func (b *baseToken[V]) setError(e error) {
	b.m.Lock()
	b.err = e
	b.flowComplete()
	b.m.Unlock()
}

func (b *baseToken[V]) flowComplete() {
	select {
	case <-b.complete:
	default:
		close(b.complete)
	}
}

type ConnectToken struct {
	baseToken[Client]
}

type ReadToken struct {
	baseToken[[]*core.DataItem]
}

type PduToken struct {
	baseToken[*core.PDU]
}

type UploadToken struct {
	baseToken[[]byte]
}

type SzlIdsToken struct {
	baseToken[[]uint16]
}

type CatalogToken struct {
	baseToken[core.Catalog]
}

type PlcStatusToken struct {
	baseToken[core.PlcStatus]
}

type UnitInfoToken struct {
	baseToken[core.UnitInfo]
}

type CommunicationInfoToken struct {
	baseToken[core.CommunicationInfo]
}

type ProtectionInfoToken struct {
	baseToken[core.ProtectionInfo]
}

type BlockListToken struct {
	baseToken[[]core.ListBlockInfo]
}

type BlockListTypeToken struct {
	baseToken[[]core.ListBlockTypeInfo]
}

type BlockInfoToken struct {
	baseToken[core.BlockInfo]
}

type ClockReadToken struct {
	baseToken[time.Time]
}

type BaseReadToken struct {
	baseToken[[]byte]
}

type SingleRawReadToken struct {
	baseToken[RawInfo]
}

type BatchRawReadToken struct {
	baseToken[[]RawInfo]
}

type SingleParsedReadToken struct {
	baseToken[any]
}

type BatchParsedReadToken struct {
	baseToken[[]any]
}
