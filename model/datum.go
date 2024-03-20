// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/util"
	"strconv"
	"strings"
	"time"
)

type ReadWriteDatum struct {
	// ReturnItems 数据项
	ReturnItems []common.ResponseItem
}

func NewReadWriteDatum(items []common.ResponseItem) *ReadWriteDatum {
	return &ReadWriteDatum{
		ReturnItems: items,
	}
}

func ReadWriteDatumFromBytes(bytes []byte, messageType common.MessageType, functionCode common.FunctionCode) (*ReadWriteDatum, error) {
	rwd := &ReadWriteDatum{}
	if len(bytes) == 0 {
		return rwd, nil
	}
	offset := 0
	remain := bytes
	for {
		var dataItem common.ResponseItem
		var err error
		// Handle write operations response result specifically
		if messageType == common.MtAckData && functionCode == common.FcWrite {
			dataItem, err = ReturnItemFromBytes(remain)
			if err != nil {
				return nil, err
			}
			rwd.ReturnItems = append(rwd.ReturnItems, dataItem)
			offset += dataItem.Len()
		} else {
			dataItem, err = DataItemFromBytes(remain)
			if err != nil {
				return nil, err
			}
			rwd.ReturnItems = append(rwd.ReturnItems, dataItem)
			offset += dataItem.Len()
			// When the data is not the last one, if the length of data is odd, S7 protocol will pad an extra byte to keep it even (the last one does not need to be padded if its length is odd)
			if dataItem.Len()%2 == 1 {
				offset++
			}
		}
		if offset >= len(bytes) {
			break
		}
		remain = bytes[offset:]
	}
	return rwd, nil
}

func (r *ReadWriteDatum) Len() int {
	if len(r.ReturnItems) == 0 {
		return 0
	}
	sum := 0
	for i := 0; i < len(r.ReturnItems); i++ {
		item := r.ReturnItems[i]
		length := item.Len()
		sum += length
		// 当数据不是最后一个的时候
		// 如果数据长度为奇数 S7协议会多填充一个字节
		// 使其数量保持为偶数（最后一个奇数长度数据不需要填充）
		_, isDataItem := item.(*DataItem)
		if i != len(r.ReturnItems)-1 && length%2 == 1 && isDataItem {
			sum++
		}
	}
	return sum
}

func (r *ReadWriteDatum) ToBytes() []byte {
	if len(r.ReturnItems) == 0 {
		return []byte{0x00}
	}
	res := make([]byte, 0, r.Len())
	for i := 0; i < len(r.ReturnItems); i++ {
		item := r.ReturnItems[i]
		length := item.Len()
		res = append(res, item.ToBytes()...)
		_, isDataItem := item.(*DataItem)
		if i != len(r.ReturnItems)-1 && length%2 == 1 && isDataItem {
			res = append(res, 0x00)
		}
	}
	return res
}

type UpDownloadDatum struct {
	// Reserved 保留
	// 字节大小：2
	// 字节序数：0-1
	Reserved []byte
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// UnkonwnBytes 未知
	// 字节大小：2
	// 字节序数：4-5
	UnkonwnBytes []byte
	// Data 数据部分
	Data []byte
}

func NewUpDownloadDatum(bytes []byte) *UpDownloadDatum {
	return &UpDownloadDatum{
		Length:       uint16(len(bytes)),
		UnkonwnBytes: []byte{0x00, 0xFB},
		Data:         bytes,
	}
}

func (u *UpDownloadDatum) Len() int {
	return common.UpDownloadDatumMinLen + len(u.Data)
}

func (u *UpDownloadDatum) ToBytes() []byte {
	res := make([]byte, u.Length)
	res = append(res, util.NumberToBytes(u.Length)...)
	res = append(res, util.NumberToBytes(u.UnkonwnBytes)...)
	res = append(res, u.Data...)
	return res
}

func UpDownloadDatumFromBytes(bytes []byte) (*UpDownloadDatum, error) {
	if len(bytes) < common.UpDownloadDatumMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "UpDownloadDatum", common.UpDownloadDatumMinLen)
	}
	l := binary.BigEndian.Uint16(bytes[:2])
	return &UpDownloadDatum{
		Length:       l,
		UnkonwnBytes: bytes[2:4],
		Data:         bytes[4 : 4+int(l)],
	}, nil
}

type UserdataDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
}

func NewUserdataDatum() *UserdataDatum {
	return &UserdataDatum{
		ReturnCode:   common.RcSuccess,
		VariableType: common.DvtOctetString,
		Length:       0,
	}
}

func UserdataDatumFromBytes(bytes []byte) (*UserdataDatum, error) {
	if len(bytes) < common.UserdataDatumLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "UserdataDatum", common.UserdataDatumLen)
	}
	return &UserdataDatum{
		ReturnCode:   common.ReturnCode(bytes[0]),
		VariableType: common.DataVariableType(bytes[1]),
		Length:       binary.BigEndian.Uint16(bytes[2:]),
	}, nil
}

func (c *UserdataDatum) Len() int {
	return common.UserdataDatumLen
}

func (c *UserdataDatum) ToBytes() []byte {
	res := make([]byte, 0, c.Len())
	res = append(res, byte(c.ReturnCode), byte(c.VariableType))
	res = append(res, util.NumberToBytes(c.Length)...)
	return res
}

type SetPasswordDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// Password 密码
	// 字节大小：8
	// 字节序数：4-11
	Password string
}

func NewSetPasswordDatum(pwd string) *SetPasswordDatum {
	return &SetPasswordDatum{
		ReturnCode:   common.RcSuccess,
		VariableType: common.DvtOctetString,
		Length:       8,
		Password:     pwd,
	}
}

func (s *SetPasswordDatum) Len() int {
	return common.SetPasswordDatumLen
}

func (s *SetPasswordDatum) ToBytes() []byte {
	res := make([]byte, 0, s.Len())
	res = append(res, byte(s.ReturnCode), byte(s.VariableType))
	res = append(res, util.NumberToBytes(s.Length)...)
	if len(s.Password) < 8 {
		s.Password += strings.Repeat("0", 8-len(s.Password))
	}
	bs := make([]byte, 0, len(s.Password))
	for i := 0; i < len(s.Password); i++ {
		b := s.Password[i]
		if i < 2 {
			b = b ^ 0x55
		} else {
			b = b ^ 0x55 ^ bs[i-2]
		}
		bs = append(bs, b)
	}
	res = append(res, bs...)
	return res
}

type ClockAckDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// Reserved 保留
	// 字节大小：1
	// 字节序数：4
	Reserved byte
	// Year1 年份前两位
	// 字节大小：1
	// 字节序数：5
	Year1 uint8
	// Year2 年份后两位
	// 字节大小：1
	// 字节序数：6
	Year2 uint8
	// Month 月份
	// 字节大小：1
	// 字节序数：7
	Month uint8
	// Day 天
	// 字节大小：1
	// 字节序数：8
	Day uint8
	// Hour 小时
	// 字节大小：1
	// 字节序数：9
	Hour uint8
	// Minute 分
	// 字节大小：1
	// 字节序数：10
	Minute uint8
	// Second 秒
	// 字节大小：1
	// 字节序数：11
	Second uint8
	// MilliSecond 毫秒
	// 字节大小：2
	// 字节序数：12-13
	MilliSecond uint16
}

func NewClockAckDatum(t time.Time) *ClockAckDatum {
	yearStr := strconv.Itoa(t.Year())
	year1, _ := hex.DecodeString(yearStr[:2])
	year2, _ := hex.DecodeString(yearStr[2:4])
	month, _ := hex.DecodeString(fmt.Sprintf("%02d", t.Month()))
	day, _ := hex.DecodeString(fmt.Sprintf("%02d", t.Day()))
	hour, _ := hex.DecodeString(fmt.Sprintf("%02d", t.Hour()))
	minute, _ := hex.DecodeString(fmt.Sprintf("%02d", t.Minute()))
	second, _ := hex.DecodeString(fmt.Sprintf("%02d", t.Second()))
	return &ClockAckDatum{
		ReturnCode:   common.RcSuccess,
		VariableType: common.DvtOctetString,
		Length:       10,
		Reserved:     0,
		Year1:        year1[0],
		Year2:        year2[0],
		Month:        month[0],
		Day:          day[0],
		Hour:         hour[0],
		Minute:       minute[0],
		Second:       second[0],
		MilliSecond:  uint16(t.Nanosecond() / 1000000),
	}
}

func ClockAckDatumFromBytes(bytes []byte) (*ClockAckDatum, error) {
	if len(bytes) < common.ClockReadAckDatumLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "ClockAckDatum", common.ClockReadAckDatumLen)
	}
	return &ClockAckDatum{
		ReturnCode:   common.ReturnCode(bytes[0]),
		VariableType: common.DataVariableType(bytes[1]),
		Length:       binary.BigEndian.Uint16(bytes[2:]),
		Reserved:     bytes[4],
		Year1:        bytes[5],
		Year2:        bytes[6],
		Month:        bytes[7],
		Day:          bytes[8],
		Hour:         bytes[9],
		Minute:       bytes[10],
		Second:       bytes[11],
		MilliSecond:  binary.BigEndian.Uint16(bytes[12:]),
	}, nil
}

func (c *ClockAckDatum) Len() int {
	return common.ClockReadAckDatumLen
}

func (c *ClockAckDatum) ToBytes() []byte {
	res := make([]byte, 0, c.Len())
	res = append(res, byte(c.ReturnCode), byte(c.VariableType))
	res = append(res, util.NumberToBytes(c.Length)...)
	res = append(res, c.Reserved, c.Year1, c.Year2, c.Month, c.Day, c.Hour, c.Minute, c.Second)
	res = append(res, util.NumberToBytes(c.MilliSecond)...)
	return res
}

type BlockInfoDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// BlockType 块类型
	// 字节大小：2
	// 字节序数：4-5
	BlockType common.BlockType
	// BlockNumber 数据块编号，范围00000-99999
	// 字节大小：5
	// 字节序数：6-10
	BlockNumber int
	// DestinationFileSystem 目标文件系统
	// 字节大小：1
	// 字节序数：11
	DestinationFileSystem common.DestinationFileSystem
}

func NewBlockInfoDatum(bt common.BlockType, dfs common.DestinationFileSystem, blockNumber int) *BlockInfoDatum {
	return &BlockInfoDatum{
		ReturnCode:            common.RcSuccess,
		VariableType:          common.DvtOctetString,
		Length:                8,
		BlockType:             bt,
		BlockNumber:           blockNumber,
		DestinationFileSystem: dfs,
	}
}

func (b *BlockInfoDatum) Len() int {
	return common.BlockInfoDatumLen
}

func (b *BlockInfoDatum) ToBytes() []byte {
	res := make([]byte, 0, b.Len())
	res = append(res, byte(b.ReturnCode), byte(b.VariableType))
	res = append(res, util.NumberToBytes(b.Length)...)
	res = append(res, util.NumberToBytes(uint16(b.BlockType))...)
	res = append(res, []byte(fmt.Sprintf("%05d", b.BlockNumber))...)
	res = append(res, byte(b.DestinationFileSystem))
	return res
}

type BlockListTypeDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// BlockType 块类型
	// 字节大小：2
	// 字节序数：4-5
	BlockType common.BlockType
}

func NewBlockListTypeDatum(bt common.BlockType) *BlockListTypeDatum {
	return &BlockListTypeDatum{
		ReturnCode:   common.RcSuccess,
		VariableType: common.DvtOctetString,
		Length:       2,
		BlockType:    bt,
	}
}

func (b *BlockListTypeDatum) Len() int {
	return common.BlockListTypeDatumLen
}

func (b *BlockListTypeDatum) ToBytes() []byte {
	res := make([]byte, 0, b.Len())
	res = append(res, byte(b.ReturnCode), byte(b.VariableType))
	res = append(res, util.NumberToBytes(b.Length)...)
	res = append(res, util.NumberToBytes(uint16(b.BlockType))...)
	return res
}

type BlockInfoAckDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// BlockType 块类型
	// 字节大小：2
	// 字节序数：4-5
	BlockType uint16
	// LengthOfInfo 信息长度
	// 字节大小：2
	// 字节序数：6-7
	LengthOfInfo uint16
	// Reserved1 保留
	// 字节大小：2
	// 字节序数：8-9
	Reserved1 []byte
	// Constant 常量
	// 字节大小：2
	// 字节序数：10-11
	Constant []byte
	// Reserved2 保留
	// 字节大小：1
	// 字节序数：12
	Reserved2 byte
	// Flags 标记
	// 字节大小：1
	// 字节序数：13
	Flags byte
	// Language 语言
	// 字节大小：1
	// 字节序数：14
	Language byte
	// SubBlkType
	// 字节大小：1
	// 字节序数：15
	SubBlkType byte
	// BlockNumber 序号
	// 字节大小：2
	// 字节序数：16-17
	BlockNumber uint16
	// LengthLoadMemory 加载内存长度
	// 字节大小：4
	// 字节序数：18-21
	LengthLoadMemory uint32
	// BlockSecurity 安全类型
	// 字节大小：4
	// 字节序数：22-25
	BlockSecurity uint32
	// CodeTimestamp 代码时间戳
	// 字节大小：6
	// 字节序数：26-31
	CodeTimestamp []byte
	// InterfaceTimestamp 接口时间戳
	// 字节大小：6
	// 字节序数：32-37
	InterfaceTimestamp []byte
	// SBBLength SBB长度
	// 字节大小：2
	// 字节序数：38-39
	SBBLength uint16
	// ADDLength ADD长度
	// 字节大小：2
	// 字节序数：40-41
	ADDLength uint16
	// LocalDataLength 本地数据长度
	// 字节大小：2
	// 字节序数：42-43
	LocalDataLength uint16
	// MC7CodeLength MC7代码长度
	// 字节大小：2
	// 字节序数：44-45
	MC7CodeLength uint16
	// Auth 作者
	// 字节大小：8
	// 字节序数：46-53
	Auth []byte
	// Family family
	// 字节大小：8
	// 字节序数：54-61
	Family []byte
	// Header 名称
	// 字节大小：8
	// 字节序数：62-69
	Header []byte
	// Version 版本
	// 字节大小：1
	// 字节序数：70
	Version uint8
	// Reserved3 保留
	// 字节大小：1
	// 字节序数：71
	Reserved3 byte
	// CheckSum 校验和
	// 字节大小：2
	// 字节序数：72-73
	CheckSum uint16
	// Reserved4 保留
	// 字节大小：4
	// 字节序数：74-77
	Reserved4 []byte
	// Reserved5 保留
	// 字节大小：4
	// 字节序数：78-81
	Reserved5 []byte
}

func BlockInfoAckDatumFromBytes(bytes []byte) (*BlockInfoAckDatum, error) {
	length := len(bytes)
	if length < common.BlockAckDatumMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "BlockInfoAckDatum", common.BlockAckDatumMinLen)
	}
	b := &BlockInfoAckDatum{
		ReturnCode:   common.ReturnCode(bytes[0]),
		VariableType: common.DataVariableType(bytes[1]),
		Length:       binary.BigEndian.Uint16(bytes[2:]),
	}
	if b.Length == 78 {
		b.BlockType = binary.BigEndian.Uint16(bytes[4:])
		b.LengthOfInfo = binary.BigEndian.Uint16(bytes[6:])
		b.Reserved1 = bytes[8:10]
		b.Constant = bytes[10:12]
		b.Reserved2 = bytes[12]
		b.Flags = bytes[13]
		b.Language = bytes[14]
		b.SubBlkType = bytes[15]
		b.BlockNumber = binary.BigEndian.Uint16(bytes[16:])
		b.LengthLoadMemory = binary.BigEndian.Uint32(bytes[18:])
		b.BlockSecurity = binary.BigEndian.Uint32(bytes[22:])
		b.CodeTimestamp = bytes[26:32]
		b.InterfaceTimestamp = bytes[32:38]
		b.SBBLength = binary.BigEndian.Uint16(bytes[38:])
		b.ADDLength = binary.BigEndian.Uint16(bytes[40:])
		b.LocalDataLength = binary.BigEndian.Uint16(bytes[42:])
		b.MC7CodeLength = binary.BigEndian.Uint16(bytes[44:])
		b.Auth = bytes[46:54]
		b.Family = bytes[54:62]
		b.Header = bytes[62:70]
		b.Version = bytes[70]
		b.Reserved3 = bytes[71]
		b.CheckSum = binary.BigEndian.Uint16(bytes[72:])
		b.Reserved4 = bytes[74:78]
		b.Reserved5 = bytes[78:82]
	}
	return b, nil
}

func (b *BlockInfoAckDatum) Len() int {
	return common.BlockAckDatumMinLen
}

func (b *BlockInfoAckDatum) ToBytes() []byte {
	res := make([]byte, 0, b.Len())
	res = append(res, byte(b.ReturnCode), byte(b.VariableType))
	res = append(res, util.NumberToBytes(b.Length)...)
	res = append(res, util.NumberToBytes(b.BlockType)...)
	res = append(res, util.NumberToBytes(b.LengthOfInfo)...)
	res = append(res, b.Reserved1...)
	res = append(res, b.Constant...)
	res = append(res, b.Reserved2, b.Flags, b.Language, b.SubBlkType)
	res = append(res, util.NumberToBytes(b.BlockNumber)...)
	res = append(res, util.NumberToBytes(b.LengthLoadMemory)...)
	res = append(res, util.NumberToBytes(b.BlockSecurity)...)
	res = append(res, b.CodeTimestamp...)
	res = append(res, b.InterfaceTimestamp...)
	res = append(res, util.NumberToBytes(b.SBBLength)...)
	res = append(res, util.NumberToBytes(b.ADDLength)...)
	res = append(res, util.NumberToBytes(b.LocalDataLength)...)
	res = append(res, util.NumberToBytes(b.MC7CodeLength)...)
	res = append(res, b.Auth...)
	res = append(res, b.Family...)
	res = append(res, b.Header...)
	res = append(res, b.Version, b.Reserved3)
	res = append(res, util.NumberToBytes(b.CheckSum)...)
	res = append(res, b.Reserved4...)
	res = append(res, b.Reserved5...)
	return res
}

type BlockListTypeAckDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// Types 块类型信息
	Types []ListBlockTypeInfo
}

func BlockListTypeAckDatumFromBytes(bytes []byte) (*BlockListTypeAckDatum, error) {
	length := len(bytes)
	if length < common.BlockAckDatumMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "BlockListTypeAckDatum", common.BlockAckDatumMinLen)
	}
	b := &BlockListTypeAckDatum{
		ReturnCode:   common.ReturnCode(bytes[0]),
		VariableType: common.DataVariableType(bytes[1]),
		Length:       binary.BigEndian.Uint16(bytes[2:]),
		Types:        make([]ListBlockTypeInfo, 0),
	}
	for i := 0; i < int(b.Length); i += 4 {
		b.Types = append(b.Types, ListBlockTypeInfo{
			Number:   binary.BigEndian.Uint16(bytes[i+4:]),
			Flags:    bytes[i+4],
			Language: bytes[i+5],
		})
	}
	return b, nil
}

func (b *BlockListTypeAckDatum) Len() int {
	return common.BlockAckDatumMinLen
}

func (b *BlockListTypeAckDatum) ToBytes() []byte {
	res := make([]byte, 0, b.Len())
	res = append(res, byte(b.ReturnCode), byte(b.VariableType))
	res = append(res, util.NumberToBytes(b.Length)...)
	for _, t := range b.Types {
		res = append(res, util.NumberToBytes(t.Number)...)
		res = append(res, t.Flags, t.Language)
	}
	return res
}

type BlockListAckDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// Blocks 块信息（数量）
	Blocks []ListBlockInfo
}

func BlockListAckDatumFromBytes(bytes []byte) (*BlockListAckDatum, error) {
	length := len(bytes)
	if length < common.BlockAckDatumMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "BlockListAckDatum", common.BlockAckDatumMinLen)
	}
	b := &BlockListAckDatum{
		ReturnCode:   common.ReturnCode(bytes[0]),
		VariableType: common.DataVariableType(bytes[1]),
		Length:       binary.BigEndian.Uint16(bytes[2:]),
		Blocks:       make([]ListBlockInfo, 0),
	}
	for i := 0; i < int(b.Length); i += 4 {
		b.Blocks = append(b.Blocks, ListBlockInfo{
			Type:  common.BlockType(binary.BigEndian.Uint16(bytes[i+4:])),
			Count: binary.BigEndian.Uint16(bytes[i+6:]),
		})
	}
	return b, nil
}

func (b *BlockListAckDatum) Len() int {
	return common.BlockAckDatumMinLen
}

func (b *BlockListAckDatum) ToBytes() []byte {
	res := make([]byte, 0, b.Len())
	res = append(res, byte(b.ReturnCode), byte(b.VariableType))
	res = append(res, util.NumberToBytes(b.Length)...)
	for _, block := range b.Blocks {
		res = append(res, util.NumberToBytes(uint16(block.Type))...)
		res = append(res, util.NumberToBytes(block.Count)...)
	}
	return res
}

type ReadSzlDatum struct {
	ReturnCode   common.ReturnCode
	VariableType common.DataVariableType
	Length       uint16
	Id           uint16
	Index        uint16
}

func NewReadSzlDatum(szlId uint16, szlIndex uint16) *ReadSzlDatum {
	return &ReadSzlDatum{
		ReturnCode:   common.RcSuccess,
		VariableType: common.DvtOctetString,
		Length:       4,
		Id:           szlId,
		Index:        szlIndex,
	}
}

func (r *ReadSzlDatum) Len() int {
	return common.ReadSzlDatumLen
}

func (r *ReadSzlDatum) ToBytes() []byte {
	res := make([]byte, 0, r.Len())
	res = append(res, byte(r.ReturnCode), byte(r.VariableType))
	res = append(res, util.NumberToBytes(r.Length)...)
	res = append(res, util.NumberToBytes(r.Id)...)
	res = append(res, util.NumberToBytes(r.Index)...)
	return res
}

type ReadSzlAckDatum struct {
	// ReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	ReturnCode common.ReturnCode
	// VariableType 数据类型
	// 字节大小：1
	// 字节序数：1
	VariableType common.DataVariableType
	// Length 长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
	// Id szl id
	// 字节大小：2
	// 字节序数：4-5
	Id uint16
	// Index szl index
	// 字节大小：2
	// 字节序数：6-7
	Index uint16
	// PartLength 数据部分长度
	// 字节大小：2
	// 字节序数：8-9
	PartLength uint16
	// PartCount 数据部分数量
	// 字节大小：2
	// 字节序数：10-11
	PartCount uint16
	// Parts 数据部分
	Parts [][]byte
}

func ReadSzlAckDatumFromBytes(bytes []byte) (*ReadSzlAckDatum, error) {
	length := len(bytes)
	if length < common.ReadSzlAckDatumMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "ReadSzlAckDatum", common.ReadSzlAckDatumMinLen)
	}
	r := &ReadSzlAckDatum{
		ReturnCode:   common.ReturnCode(bytes[0]),
		VariableType: common.DataVariableType(bytes[1]),
		Length:       binary.BigEndian.Uint16(bytes[2:]),
	}
	if r.Length > 0 {
		r.Id = binary.BigEndian.Uint16(bytes[4:])
		r.Index = binary.BigEndian.Uint16(bytes[6:])
		r.PartLength = binary.BigEndian.Uint16(bytes[8:])
		r.PartCount = binary.BigEndian.Uint16(bytes[10:])
		r.Parts = make([][]byte, 0)
		if r.PartCount > 0 {
			offset := common.ReadSzlAckDatumMinLen
			for i := 0; i < int(r.PartCount); i++ {
				if length >= offset+int(r.PartLength) {
					bs := bytes[offset : offset+int(r.PartLength)]
					r.Parts = append(r.Parts, bs)
					offset += int(r.PartLength)
				} else {
					return nil, common.ErrorWithCode(common.ErrCliSzlPartsInvalid)
				}
			}
		}
	}
	return r, nil
}

func (r *ReadSzlAckDatum) Len() int {
	i := common.ReadSzlAckDatumMinLen
	if r.Length > 0 {
		i += 8
		for _, part := range r.Parts {
			i += len(part)
		}
	}
	return i
}

func (r *ReadSzlAckDatum) ToBytes() []byte {
	res := make([]byte, 0, r.Len())
	res = append(res, byte(r.ReturnCode), byte(r.VariableType))
	res = append(res, util.NumberToBytes(r.Length)...)
	if r.Length > 0 {
		res = append(res, util.NumberToBytes(r.Id)...)
		res = append(res, util.NumberToBytes(r.Index)...)
		res = append(res, util.NumberToBytes(r.PartLength)...)
		res = append(res, util.NumberToBytes(r.PartCount)...)
		for _, part := range r.Parts {
			res = append(res, part...)
		}
	}
	return res
}
