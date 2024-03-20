// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package core

import (
	"encoding/binary"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/util"
)

type StandardRequestItem struct {
	// SpecificationType 变量规范
	// 对于读/写消息，它总是具有值0x12
	// 字节大小：1
	// 字节序数：0
	SpecificationType byte
	// LengthOfFollowing 其余部分的长度规范
	// 字节大小：1
	// 字节序数：1
	LengthOfFollowing byte
	// SyntaxId 寻址模式和项结构其余部分的格式，它具有任意类型寻址的常量值0x10
	// 字节大小：1
	// 字节序数：2
	SyntaxId common.SyntaxID
	// VariableType 变量的类型和长度BIT，BYTE，WORD，DWORD，COUNTER
	// 字节大小：1
	// 字节序数：3
	VariableType common.ParamVariableType
	// Count 读取长度
	// 字节大小：2
	// 字节序数：4-5
	Count uint16
	// DbNumber DB编号
	// 如果访问的不是DB区域，此处为0x0000
	// 字节大小：2
	// 字节序数：6-7
	DbNumber uint16
	// Area 存储区类型
	// 字节大小：1
	// 字节序数：8
	Area common.AreaType
	// ByteAddress 字节地址
	// 位于开始字节地址address中3个字节，从第4位开始计数
	// 字节大小：3
	// 字节序数：9-11
	ByteAddress int
	// BitAddress 位地址
	// 位于开始字节地址address中3个字节的最后3位
	BitAddress int
}

func NewStandardRequestItem(area common.AreaType, dbNumber int, variableType common.ParamVariableType, byteAddress int, bitAddress int, count int) *StandardRequestItem {
	return &StandardRequestItem{
		SpecificationType: 0x12,
		LengthOfFollowing: 0x0A,
		SyntaxId:          common.SiAny,
		VariableType:      variableType,
		Count:             uint16(count),
		DbNumber:          uint16(dbNumber),
		Area:              area,
		ByteAddress:       byteAddress,
		BitAddress:        bitAddress,
	}
}

func (s *StandardRequestItem) Len() int {
	return common.StandardRequestItemLen
}

func (s *StandardRequestItem) ToBytes() []byte {
	res := make([]byte, 0, s.Len())
	res = append(res, s.SpecificationType, s.LengthOfFollowing, byte(s.SyntaxId), byte(s.VariableType))
	res = append(res, util.NumberToBytes(s.Count)...)
	res = append(res, util.NumberToBytes(s.DbNumber)...)
	res = append(res, byte(s.Area))
	res = append(res, util.NumberToBytes(uint32(s.ByteAddress<<3 + s.BitAddress))[1:]...)
	return res
}

func StandardRequestItemFromBytes(bytes []byte) (*StandardRequestItem, error) {
	return StandardRequestItemFromBytesWithOffset(bytes, 0)
}

func StandardRequestItemFromBytesWithOffset(bytes []byte, offset int) (*StandardRequestItem, error) {
	if len(bytes) < offset+common.StandardRequestItemLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "StandardRequestItem", common.StandardRequestItemLen)
	}
	u := binary.BigEndian.Uint32(append([]byte{0x00}, bytes[offset+9:offset+12]...))
	return &StandardRequestItem{
		SpecificationType: bytes[offset],
		LengthOfFollowing: bytes[offset+1],
		SyntaxId:          common.SyntaxID(bytes[offset+2]),
		VariableType:      common.ParamVariableType(bytes[offset+3]),
		Count:             binary.BigEndian.Uint16(bytes[offset+4:]),
		DbNumber:          binary.BigEndian.Uint16(bytes[offset+6:]),
		Area:              common.AreaType(bytes[offset+8]),
		ByteAddress:       int(u >> 3),
		BitAddress:        int(u & 0x07),
	}, nil
}

func (s *StandardRequestItem) GetSpecificationType() byte {
	return s.SpecificationType
}

func (s *StandardRequestItem) GetLengthOfFollowing() byte {
	return s.LengthOfFollowing
}

func (s *StandardRequestItem) GetSyntaxId() common.SyntaxID {
	return s.SyntaxId
}

func (s *StandardRequestItem) SetSpecificationType(b byte) {
	s.SpecificationType = b
}

func (s *StandardRequestItem) SetLengthOfFollowing(b byte) {
	s.LengthOfFollowing = b
}

func (s *StandardRequestItem) SetSyntaxId(id common.SyntaxID) {
	s.SyntaxId = id
}

type NckRequestItem struct {
	// SpecificationType 变量规范
	// 对于读/写消息，它总是具有值0x12
	// 字节大小：1
	// 字节序数：0
	SpecificationType byte
	// LengthOfFollowing 其余部分的长度规范
	// 字节大小：1
	// 字节序数：1
	LengthOfFollowing byte
	// SyntaxId 寻址模式和项结构其余部分的格式，它具有任意类型寻址的常量值0x10
	// 字节大小：1
	// 字节序数：2
	SyntaxId common.SyntaxID
	// Area NCK区域
	// 字节大小：1
	// 字节序数：3
	Area common.AreaType
	// Unit 通道编号
	// 字节大小：1
	// 字节序数：3
	Unit byte
	// ColumnNumber 列编号
	// 字节大小：2
	// 字节序数：4-5
	ColumnNumber uint16
	// LineNumber 行编号
	// 字节大小：2
	// 字节序数：6-7
	LineNumber uint16
	// Module 模块名
	// 字节大小：1
	// 字节序数：8
	Module common.NckModule
	// LineCount 行个数
	// 字节大小：1
	// 字节序数：9
	LineCount uint8
}

func NckRequestItemFromBytes(bytes []byte) (*NckRequestItem, error) {
	return NckRequestItemFromBytesWithOffset(bytes, 0)
}

func NckRequestItemFromBytesWithOffset(bytes []byte, offset int) (*NckRequestItem, error) {
	if len(bytes) < offset+common.NckRequestItemLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "NckRequestItem", common.NckRequestItemLen)
	}
	return &NckRequestItem{
		SpecificationType: bytes[offset],
		LengthOfFollowing: bytes[offset+1],
		SyntaxId:          common.SyntaxID(bytes[offset+2]),
		Area:              common.AreaType((bytes[offset+3] & 0xE0) >> 5),
		Unit:              bytes[offset+3] & 0x1F,
		ColumnNumber:      binary.BigEndian.Uint16(bytes[offset+4:]),
		LineNumber:        binary.BigEndian.Uint16(bytes[offset+6:]),
		Module:            common.NckModule(bytes[offset+8]),
		LineCount:         bytes[offset+9],
	}, nil
}

func (n *NckRequestItem) Len() int {
	return common.NckRequestItemLen
}

func (n *NckRequestItem) ToBytes() []byte {
	res := make([]byte, 0, n.Len())
	res = append(res, n.SpecificationType, n.LengthOfFollowing, byte(n.SyntaxId), (byte(n.Area)<<5&0xE0)|(n.Unit&0x1F), n.Unit)
	res = append(res, util.NumberToBytes(n.ColumnNumber)...)
	res = append(res, util.NumberToBytes(n.LineNumber)...)
	res = append(res, byte(n.Module), n.LineCount)
	return res
}

func (n *NckRequestItem) GetSpecificationType() byte {
	return n.SpecificationType
}

func (n *NckRequestItem) GetLengthOfFollowing() byte {
	return n.LengthOfFollowing
}

func (n *NckRequestItem) GetSyntaxId() common.SyntaxID {
	return n.SyntaxId
}

func (n *NckRequestItem) SetSpecificationType(b byte) {
	n.SpecificationType = b
}

func (n *NckRequestItem) SetLengthOfFollowing(b byte) {
	n.LengthOfFollowing = b
}

func (n *NckRequestItem) SetSyntaxId(id common.SyntaxID) {
	n.SyntaxId = id
}
