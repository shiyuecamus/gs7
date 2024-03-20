// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package core

import (
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/util"
)

type COTPData struct {
	// Length 长度（但并不包含length这个字段）
	// 字节大小：1
	// 字节序数：0
	Length uint8
	// PduType PDU类型（PtData 数据）
	// 字节大小：1
	// 字节序数：1
	PduType common.PduType
	// TpduNumber TPDU编号
	// 字节大小：1，后面7位
	// 字节序数：2
	TpduNumber byte
	// LastDataUnit 是否最后一个数据单元
	// 字节大小：1，最高位，7位
	// 字节序数：2
	LastDataUnit bool
}

func NewCOTPData() *COTPData {
	return &COTPData{
		Length:       0x02,
		PduType:      common.PtData,
		TpduNumber:   0x00,
		LastDataUnit: true,
	}
}

func COTPDataFromBytes(bytes []byte) (*COTPData, error) {
	if len(bytes) < common.CotpDataLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "COTPData", common.CotpDataLen)
	}

	return &COTPData{
		Length:       bytes[0],
		PduType:      common.PduType(bytes[1]),
		TpduNumber:   bytes[2] & 0x7F,
		LastDataUnit: util.GetBoolAt(bytes[2], 7),
	}, nil
}

func (c *COTPData) Len() int {
	return common.CotpDataLen
}

func (c *COTPData) ToBytes() []byte {
	res := make([]byte, 0, c.Len())
	res = append(res, c.Length, byte(c.PduType))
	b := util.SetBoolAt(0x00, 7, c.LastDataUnit)
	res = append(res, b|c.TpduNumber)
	return res
}

func (c *COTPData) GetLength() byte {
	return c.Length
}

func (c *COTPData) GetPduType() common.PduType {
	return c.PduType
}

func (c *COTPData) SetLength(length byte) {
	c.Length = length
}

func (c *COTPData) SetPduType(pduType common.PduType) {
	c.PduType = pduType
}

type COTPConnection struct {
	// Length 长度（但并不包含length这个字段）
	// 字节大小：1
	// 字节序数：0
	Length uint8
	// PduType PDU类型（CRConnect Request 连接请求）
	// 字节大小：1
	// 字节序数：1
	PduType common.PduType
	// DestinationReference 目标引用，用来唯一标识目标
	// 字节大小：2
	// 字节序数：2-3
	DestinationReference []byte
	// SourceReference 源引用
	// 字节大小：2
	// 字节序数：4-5
	SourceReference []byte
	// Flags 扩展格式/流控制  前四位标识Class，  倒数第二位Extended formats
	// 倒数第一位No explicit flow control
	// 字节大小：1
	// 字节序数：6
	Flags byte
	// ParameterCodeTpduSize 参数代码TPDU-Size
	// 字节大小：1
	// 字节序数：7
	ParameterCodeTpduSize byte
	// ParameterLength1 参数长度
	// 字节大小：1
	// 字节序数：8
	ParameterLength1 byte
	// TpduSize TPDU大小 TPDU Size (2^10 = 1024)
	// 字节大小：1
	// 字节序数：9
	TpduSize byte
	// ParameterCodeSrcTsap 参数代码SRC-TASP
	// 字节大小：1
	// 字节序数：10
	ParameterCodeSrcTsap byte
	// ParameterLength2 参数长度
	// 字节大小：1
	// 字节序数：11
	ParameterLength2 byte
	// SourceTsap SourceTSAP/Rack
	// 字节大小：2
	// 字节序数：12-13
	SourceTsap []byte
	// ParameterCodeDstTsap 参数代码DST-TASP
	// 字节大小：1
	// 字节序数：14
	ParameterCodeDstTsap byte
	// ParameterLength3 参数长度
	// 字节大小：1
	// 字节序数：15
	ParameterLength3 byte
	// DestinationTsap Slot
	// 字节大小：2
	// 字节序数：16-17
	DestinationTsap []byte
}

func NewCOTPConnection() *COTPConnection {
	return &COTPConnection{
		Length:                0x00,
		PduType:               common.PtConnectRequest,
		DestinationReference:  []byte{0x00, 0x00},
		SourceReference:       []byte{0x00, 0x01},
		Flags:                 0x00,
		ParameterCodeTpduSize: 0xC0,
		ParameterLength1:      0x01,
		TpduSize:              0x01,
		ParameterCodeSrcTsap:  0xC1,
		ParameterLength2:      0x02,
		SourceTsap:            []byte{0x01, 0x00},
		ParameterCodeDstTsap:  0xC2,
		ParameterLength3:      0x02,
		DestinationTsap:       []byte{0x01, 0x00},
	}
}

// CRConnect Request 连接请求
// @param local  本地参数
// @param remote 远程参数
func NewCOTPConnectionForRequest(local uint16, remote uint16) *COTPConnection {
	return &COTPConnection{
		Length:                0x11,
		PduType:               common.PtConnectRequest,
		DestinationReference:  []byte{0x00, 0x00},
		SourceReference:       []byte{0x00, 0x01},
		Flags:                 0x00,
		ParameterCodeTpduSize: 0xC0,
		ParameterLength1:      0x01,
		TpduSize:              0x0A,
		ParameterCodeSrcTsap:  0xC1,
		ParameterLength2:      0x02,
		SourceTsap:            util.NumberToBytes(local),
		ParameterCodeDstTsap:  0xC2,
		ParameterLength3:      0x02,
		DestinationTsap:       util.NumberToBytes(remote),
	}
}

func NewCOTPConnectionForConfirm(request *COTPConnection) *COTPConnection {
	return &COTPConnection{
		Length:                0x11,
		PduType:               common.PtConnectConfirm,
		DestinationReference:  []byte{0x00, 0x01},
		SourceReference:       request.SourceReference,
		Flags:                 request.Flags,
		ParameterCodeTpduSize: request.ParameterCodeTpduSize,
		ParameterLength1:      request.ParameterLength1,
		TpduSize:              request.TpduSize,
		ParameterCodeSrcTsap:  request.ParameterCodeSrcTsap,
		ParameterLength2:      request.ParameterLength2,
		SourceTsap:            request.SourceTsap,
		ParameterCodeDstTsap:  request.ParameterCodeDstTsap,
		ParameterLength3:      request.ParameterLength3,
		DestinationTsap:       request.DestinationTsap,
	}
}

func COTPConnectionFromBytes(bytes []byte) (*COTPConnection, error) {
	if len(bytes) < common.CotpConnectionLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "COTPConnection", common.CotpConnectionLen)
	}
	return &COTPConnection{
		Length:                bytes[0],
		PduType:               common.PduType(bytes[1]),
		DestinationReference:  bytes[2:4],
		SourceReference:       bytes[4:6],
		Flags:                 bytes[6],
		ParameterCodeTpduSize: bytes[7],
		ParameterLength1:      bytes[8],
		TpduSize:              bytes[9],
		ParameterCodeSrcTsap:  bytes[10],
		ParameterLength2:      bytes[11],
		SourceTsap:            bytes[12:14],
		ParameterCodeDstTsap:  bytes[14],
		ParameterLength3:      bytes[15],
		DestinationTsap:       bytes[16:18],
	}, nil
}

func (c *COTPConnection) Len() int {
	return common.CotpConnectionLen
}

func (c *COTPConnection) ToBytes() []byte {
	res := make([]byte, 0, c.Len())
	res = append(res, c.Length, byte(c.PduType))
	res = append(res, c.DestinationReference...)
	res = append(res, c.SourceReference...)
	res = append(res, c.Flags, c.ParameterCodeTpduSize, c.ParameterLength1, c.TpduSize, c.ParameterCodeSrcTsap, c.ParameterLength2)
	res = append(res, c.SourceTsap...)
	res = append(res, c.ParameterCodeDstTsap, c.ParameterLength3)
	res = append(res, c.DestinationTsap...)
	return res
}

func (c *COTPConnection) GetLength() byte {
	return c.Length
}

func (c *COTPConnection) GetPduType() common.PduType {
	return c.PduType
}

func (c *COTPConnection) SetLength(length byte) {
	c.Length = length
}

func (c *COTPConnection) SetPduType(pduType common.PduType) {
	c.PduType = pduType
}
