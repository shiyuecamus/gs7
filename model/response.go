// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"encoding/binary"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/util"
	"github.com/spf13/cast"
)

type DataItem struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 变量类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Count 数据长度
	// 按位进行计算的，如果是字节数据读取需要进行 /8 或 *8操作
	// 如果是位数据，不需要任何额外操作
	Count uint16
	// 数据内容
	Data []byte
}

func NewAckDataItem(bytes []byte, variableType common.DataVariableType) *DataItem {
	return &DataItem{
		ReturnCode:   common.RcSuccess,
		VariableType: variableType,
		Count:        uint16(len(bytes)),
		Data:         bytes,
	}
}

func NewReqDataItem(bytes []byte, variableType common.DataVariableType) *DataItem {
	return &DataItem{
		ReturnCode:   common.RcReserved,
		VariableType: variableType,
		Count:        uint16(len(bytes)),
		Data:         bytes,
	}
}

func NewReqDataItemByBool(b bool) *DataItem {
	return &DataItem{
		ReturnCode:   common.RcReserved,
		VariableType: common.DvtBit,
		Count:        1,
		Data:         []byte{cast.ToUint8(b)},
	}
}

func NewReqDataItemByByte(b byte) *DataItem {
	return &DataItem{
		ReturnCode:   common.RcReserved,
		VariableType: common.DvtByteWordDword,
		Count:        1,
		Data:         []byte{b},
	}
}

func NewReqDataItemByBytes(bytes []byte) *DataItem {
	return &DataItem{
		ReturnCode:   common.RcReserved,
		VariableType: common.DvtByteWordDword,
		Count:        uint16(len(bytes)),
		Data:         bytes,
	}
}

func DataItemFromBytes(bytes []byte) (*DataItem, error) {
	if len(bytes) < common.DataItemMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "DataItem", common.DataItemMinLen)
	}
	d := &DataItem{
		ReturnCode:   common.ReturnCode(bytes[0]),
		VariableType: common.DataVariableType(bytes[1]),
	}
	switch d.VariableType {
	case common.DvtNull, common.DvtByteWordDword, common.DvtInt:
		d.Count = binary.BigEndian.Uint16(bytes[2:]) / 8
		break
	case common.DvtBit, common.DvtDint, common.DvtReal, common.DvtOctetString:
		d.Count = binary.BigEndian.Uint16(bytes[2:])
		break
	default:
		return nil, common.ErrorWithCode(common.ErrVariableTypeUnrecognized, d.VariableType)
	}
	if d.VariableType != common.DvtNull {
		d.Data = bytes[4 : 4+d.Count]
	}
	return d, nil
}

func (d *DataItem) Len() int {
	return common.DataItemMinLen + len(d.Data)
}

func (d *DataItem) ToBytes() []byte {
	res := make([]byte, 0, d.Len())
	res = append(res, byte(d.ReturnCode), byte(d.VariableType))
	switch d.VariableType {
	case common.DvtNull, common.DvtByteWordDword, common.DvtInt:
		res = append(res, util.NumberToBytes(d.Count*8)...)
		break
	case common.DvtBit, common.DvtDint, common.DvtReal, common.DvtOctetString:
		res = append(res, util.NumberToBytes(d.Count)...)
		break
	}
	res = append(res, d.Data...)
	return res
}

func (d *DataItem) GetReturnCode() common.ReturnCode {
	return d.ReturnCode
}

func (d *DataItem) SetReturnCode(code common.ReturnCode) {
	d.ReturnCode = code
}

type ReturnItem struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
}

func NewReturnItem(code common.ReturnCode) *ReturnItem {
	return &ReturnItem{
		ReturnCode: code,
	}
}

func ReturnItemFromBytes(bytes []byte) (*ReturnItem, error) {
	if len(bytes) < common.ReturnItemLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "ReturnItem", common.ReturnItemLen)
	}
	return &ReturnItem{
		ReturnCode: common.ReturnCode(bytes[0]),
	}, nil
}

func (r *ReturnItem) Len() int {
	return common.ReturnItemLen
}

func (r *ReturnItem) ToBytes() []byte {
	res := make([]byte, 0, r.Len())
	res = append(res, byte(r.ReturnCode))
	return res
}

func (r *ReturnItem) GetReturnCode() common.ReturnCode {
	return r.ReturnCode
}

func (r *ReturnItem) SetReturnCode(code common.ReturnCode) {
	r.ReturnCode = code
}
