// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"encoding/binary"
	"gs7/common"
	"gs7/util"
)

type RequestHeader struct {
	// ProtocolId 协议id
	// 字节大小：1
	// 字节序数：0
	ProtocolId byte
	// MessageType pdu（协议数据单元（Protocol Data Unit））的类型
	// 字节大小：1
	// 字节序数：1
	MessageType common.MessageType
	// Reserved 保留
	// 字节大小：2
	// 字节序数：2-3
	Reserved []byte
	// PduReference pdu的参考–由主站生成，每次新传输递增，大端
	// 字节大小：2
	// 字节序数：4-5
	PduReference uint16
	// ParameterLength 参数的长度（大端）
	// 字节大小：2
	// 字节序数：6-7
	ParameterLength uint16
	// DataLength 数据的长度（大端）
	// 字节大小：2
	// 字节序数：8-9
	DataLength uint16
}

func NewRequestHeader(requestId uint16) *RequestHeader {
	return &RequestHeader{
		ProtocolId:      0x32,
		MessageType:     common.MtJob,
		Reserved:        []byte{0x00, 0x00},
		PduReference:    requestId,
		ParameterLength: 0,
		DataLength:      0,
	}
}

func NewUserDataHeader(requestId uint16) *RequestHeader {
	return &RequestHeader{
		ProtocolId:      0x32,
		MessageType:     common.MtUserData,
		Reserved:        []byte{0x00, 0x00},
		PduReference:    requestId,
		ParameterLength: 0,
		DataLength:      0,
	}
}

func RequestHeaderFromBytes(bytes []byte) (*RequestHeader, error) {
	if len(bytes) < common.RequestHeaderLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "RequestHeader", common.RequestHeaderLen)
	}
	return &RequestHeader{
		ProtocolId:      bytes[0],
		MessageType:     common.MessageType(bytes[1]),
		Reserved:        bytes[2:4],
		PduReference:    binary.BigEndian.Uint16(bytes[4:]),
		ParameterLength: binary.BigEndian.Uint16(bytes[6:]),
		DataLength:      binary.BigEndian.Uint16(bytes[8:]),
	}, nil
}

func (h *RequestHeader) Len() int {
	return common.RequestHeaderLen
}

func (h *RequestHeader) ToBytes() []byte {
	res := make([]byte, 0, h.Len())
	res = append(res, h.ProtocolId, byte(h.MessageType))
	res = append(res, h.Reserved...)
	res = append(res, util.NumberToBytes(h.PduReference)...)
	res = append(res, util.NumberToBytes(h.ParameterLength)...)
	res = append(res, util.NumberToBytes(h.DataLength)...)
	return res
}

func (h *RequestHeader) GetProtocolId() byte {
	return h.ProtocolId
}

func (h *RequestHeader) GetMessageType() common.MessageType {
	return h.MessageType
}

func (h *RequestHeader) GetReserved() []byte {
	return h.Reserved
}

func (h *RequestHeader) GetPduReference() uint16 {
	return h.PduReference
}

func (h *RequestHeader) GetParameterLength() uint16 {
	return h.ParameterLength
}

func (h *RequestHeader) GetDataLength() uint16 {
	return h.DataLength
}

func (h *RequestHeader) SetProtocolId(b byte) {
	h.ProtocolId = b
}

func (h *RequestHeader) SetMessageType(messageType common.MessageType) {
	h.MessageType = messageType
}

func (h *RequestHeader) SetReserved(bytes []byte) {
	h.Reserved = bytes
}

func (h *RequestHeader) SetPduReference(u uint16) {
	h.PduReference = u
}

func (h *RequestHeader) SetParameterLength(u uint16) {
	h.ParameterLength = u
}

func (h *RequestHeader) SetDataLength(u uint16) {
	h.DataLength = u
}

type AckHeader struct {
	// ProtocolId 协议id
	// 字节大小：1
	// 字节序数：0
	ProtocolId byte
	// MessageType pdu（协议数据单元（Protocol Data Unit））的类型
	// 字节大小：1
	// 字节序数：1
	MessageType common.MessageType
	// Reserved 保留
	// 字节大小：2
	// 字节序数：2-3
	Reserved []byte
	// PduReference pdu的参考–由主站生成，每次新传输递增，大端
	// 字节大小：2
	// 字节序数：4-5
	PduReference uint16
	// ParameterLength 参数的长度（大端）
	// 字节大小：2
	// 字节序数：6-7
	ParameterLength uint16
	// DataLength 数据的长度（大端）
	// 字节大小：2
	// 字节序数：8-9
	DataLength uint16
	// ErrorClass 错误类型
	// 字节大小：1
	// 字节序数：10
	ErrorClass byte
	// ErrorCode 错误码
	// 本来是1个字节的，但本质上errorCode（真正） = errorClass + errorCode（原）
	// 字节大小：2
	// 字节序数：10-11
	ErrorCode []byte
}

func AckHeaderFromBytes(bytes []byte) (*AckHeader, error) {
	if len(bytes) < common.AckHeaderLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "AckHeader", common.AckHeaderLen)
	}
	return &AckHeader{
		ProtocolId:      bytes[0],
		MessageType:     common.MessageType(bytes[1]),
		Reserved:        bytes[2:4],
		PduReference:    binary.BigEndian.Uint16(bytes[4:]),
		ParameterLength: binary.BigEndian.Uint16(bytes[6:]),
		DataLength:      binary.BigEndian.Uint16(bytes[8:]),
		ErrorClass:      bytes[10],
		ErrorCode:       bytes[10:12],
	}, nil
}

func (a *AckHeader) Len() int {
	return common.AckHeaderLen
}

func (a *AckHeader) ToBytes() []byte {
	res := make([]byte, 0, a.Len())
	res = append(res, a.ProtocolId, byte(a.MessageType))
	res = append(res, a.Reserved...)
	res = append(res, util.NumberToBytes(a.PduReference)...)
	res = append(res, util.NumberToBytes(a.ParameterLength)...)
	res = append(res, util.NumberToBytes(a.DataLength)...)
	res = append(res, a.ErrorClass)
	res = append(res, util.NumberToBytes(a.ErrorCode)...)
	return res
}

func (a *AckHeader) GetProtocolId() byte {
	return a.ProtocolId
}

func (a *AckHeader) SetProtocolId(b byte) {
	a.ProtocolId = b
}

func (a *AckHeader) GetMessageType() common.MessageType {
	return a.MessageType
}

func (a *AckHeader) SetMessageType(messageType common.MessageType) {
	a.MessageType = messageType
}

func (a *AckHeader) GetReserved() []byte {
	return a.Reserved
}

func (a *AckHeader) SetReserved(bytes []byte) {
	a.Reserved = bytes
}

func (a *AckHeader) GetPduReference() uint16 {
	return a.PduReference
}

func (a *AckHeader) SetPduReference(u uint16) {
	a.PduReference = u
}

func (a *AckHeader) GetParameterLength() uint16 {
	return a.ParameterLength
}

func (a *AckHeader) SetParameterLength(u uint16) {
	a.ParameterLength = u
}

func (a *AckHeader) GetDataLength() uint16 {
	return a.DataLength
}

func (a *AckHeader) SetDataLength(u uint16) {
	a.DataLength = u
}
