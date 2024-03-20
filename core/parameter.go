// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package core

import (
	"encoding/binary"
	"fmt"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/util"
	"github.com/spf13/cast"
)

type StandardParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
}

func NewStandardParameter(code common.FunctionCode) *StandardParameter {
	return &StandardParameter{
		FunctionCode: code,
	}
}

func (s *StandardParameter) Len() int {
	return common.StandardParameterLen
}

func (s *StandardParameter) ToBytes() []byte {
	res := make([]byte, 0, s.Len())
	res = append(res, byte(s.FunctionCode))
	return res
}

type ReadWriteParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// ItemCount Request Item结构的数量
	// 字节大小：1
	// 字节序数：1
	ItemCount uint8
	// RequestItems 可重复的请求项
	RequestItems []common.RequestItem
}

func NewReqReadWriteParameter(functionCode common.FunctionCode, items []common.RequestItem) *ReadWriteParameter {
	return &ReadWriteParameter{
		FunctionCode: functionCode,
		RequestItems: items,
		ItemCount:    uint8(len(items)),
	}
}

func NewAckReadWriteParameter(request *ReadWriteParameter) *ReadWriteParameter {
	return &ReadWriteParameter{
		FunctionCode: request.FunctionCode,
		ItemCount:    request.ItemCount,
	}
}

func ReadWriteParameterFromBytes(bytes []byte) (*ReadWriteParameter, error) {
	dataLen := len(bytes)
	if dataLen < common.ReadWriteParameterMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "ReadWriteParameter", common.ReadWriteParameterMinLen)
	}
	rwp := &ReadWriteParameter{
		FunctionCode: common.FunctionCode(bytes[0]),
		ItemCount:    bytes[1],
	}
	if rwp.ItemCount == 0 || dataLen == 2 {
		return rwp, nil
	}
	offset := 2
	for i := 0; i < cast.ToInt(rwp.ItemCount); i++ {
		if item, err := parseItem(bytes, offset); err != nil {
			return nil, err
		} else {
			rwp.RequestItems = append(rwp.RequestItems, item)
			offset += item.Len()
		}
	}
	return rwp, nil
}

func (r *ReadWriteParameter) Len() int {
	l := common.ReadWriteParameterMinLen
	for _, item := range r.RequestItems {
		l += item.Len()
	}
	return l
}

func (r *ReadWriteParameter) ToBytes() []byte {
	res := make([]byte, 0, r.Len())
	res = append(res, byte(r.FunctionCode), r.ItemCount)
	for _, item := range r.RequestItems {
		res = append(res, item.ToBytes()...)
	}
	return res
}

func parseItem(bytes []byte, offset int) (common.RequestItem, error) {
	switch syntaxID := common.SyntaxID(bytes[2+offset]); syntaxID {
	case common.SiAny:
		return StandardRequestItemFromBytesWithOffset(bytes, offset)
	case common.SiNck:
		return NckRequestItemFromBytesWithOffset(bytes, offset)
	default:
		return nil, common.ErrorWithCode(common.ErrCliRequestItemInvalid)
	}
}

type PlcStopParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// UnknownBytes 未知字节，固定参数
	// 字节大小：5
	// 字节序数：1-5
	UnknownBytes []byte
	// LengthPart 服务名长度，后续字节长度，不包含自身
	// 字节大小：1
	// 字节序数：6
	LengthPart uint8
	// PiService 程序调用的服务名
	PiService string
}

func NewPlcStopParameter() *PlcStopParameter {
	return &PlcStopParameter{
		FunctionCode: common.FcStop,
		UnknownBytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00},
		PiService:    "P_PROGRAM",
		LengthPart:   uint8(len("P_PROGRAM")),
	}
}

func PlcStopParameterFromBytes(bytes []byte) (*PlcStopParameter, error) {
	if len(bytes) < common.PlcStopParameterMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "PlcStopParameter", common.PlcStopParameterMinLen)
	}
	p := &PlcStopParameter{
		FunctionCode: common.FunctionCode(bytes[0]),
		UnknownBytes: bytes[1:6],
		LengthPart:   bytes[6],
	}
	if p.LengthPart == 0 {
		p.PiService = ""
	} else {
		p.PiService = string(bytes[7:])
	}
	return p, nil
}

func (p *PlcStopParameter) Len() int {
	return common.PlcStopParameterMinLen + len(p.PiService)
}

func (p *PlcStopParameter) ToBytes() []byte {
	res := make([]byte, 0, p.Len())
	res = append(res, byte(p.FunctionCode))
	res = append(res, p.UnknownBytes...)
	res = append(res, p.LengthPart)
	res = append(res, []byte(p.PiService)...)
	return res
}

// PlcControlParameter 启动参数
type PlcControlParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// UnknownBytes 未知字节，固定参数
	// 字节大小：7
	// 字节序数：1-7
	UnknownBytes []byte
	// ParameterBlockLength 参数块长度
	// 字节大小：2
	// 字节序数：8-9
	ParameterBlockLength uint16
	// ParameterBlock 参数块内容
	ParameterBlock common.PlcControlParamBlock
	// LengthPart 服务名长度，后续字节长度，不包含自身
	LengthPart uint8
	// PiService 程序调用的服务名
	PiService string
}

// NewHotRestartPlcControlParameter 热重启
func NewHotRestartPlcControlParameter() *PlcControlParameter {
	block := NewPlcControlStringParamBlock("")
	piService := "P_PROGRAM"
	return &PlcControlParameter{
		FunctionCode:         common.FcControl,
		UnknownBytes:         []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFD},
		ParameterBlockLength: uint16(block.Len()),
		ParameterBlock:       block,
		LengthPart:           uint8(len(piService)),
		PiService:            piService,
	}
}

// NewColdRestartPlcControlParameter 冷启动
func NewColdRestartPlcControlParameter() *PlcControlParameter {
	block := NewPlcControlStringParamBlock("C ")
	piService := "P_PROGRAM"
	return &PlcControlParameter{
		FunctionCode:         common.FcControl,
		UnknownBytes:         []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFD},
		ParameterBlockLength: uint16(block.Len()),
		ParameterBlock:       block,
		LengthPart:           uint8(len(piService)),
		PiService:            piService,
	}
}

// NewCopyRamToRomPlcControlParameter 将ram复制到rom中
func NewCopyRamToRomPlcControlParameter() *PlcControlParameter {
	block := NewPlcControlStringParamBlock("")
	piService := "_GARB"
	return &PlcControlParameter{
		FunctionCode:         common.FcControl,
		UnknownBytes:         []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFD},
		ParameterBlockLength: uint16(block.Len()),
		ParameterBlock:       block,
		LengthPart:           uint8(len(piService)),
		PiService:            piService,
	}
}

// NewCompressPlcControlParameter 压缩
func NewCompressPlcControlParameter() *PlcControlParameter {
	block := NewPlcControlStringParamBlock("EP")
	piService := "_MODU"
	return &PlcControlParameter{
		FunctionCode:         common.FcControl,
		UnknownBytes:         []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFD},
		ParameterBlockLength: uint16(block.Len()),
		ParameterBlock:       block,
		LengthPart:           uint8(len(piService)),
		PiService:            piService,
	}
}

func NewInsertPlcControlParameter(bt common.BlockType, bn int, dfs common.DestinationFileSystem) *PlcControlParameter {
	data := make([]byte, 0)
	data = append(data, util.NumberToBytes(uint16(bt))...)
	data = append(data, []byte(fmt.Sprintf("%05d", bn))...)
	data = append(data, byte(dfs))
	block := NewPlcControlInsertParamBlock([]string{string(data)})
	piService := "_INSE"
	return &PlcControlParameter{
		FunctionCode:         common.FcControl,
		UnknownBytes:         []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFD},
		ParameterBlockLength: uint16(block.Len()),
		ParameterBlock:       block,
		LengthPart:           uint8(len(piService)),
		PiService:            piService,
	}
}

func PlcControlParameterFromBytes(bytes []byte) (*PlcControlParameter, error) {
	if len(bytes) < common.PlcControlParameterMinLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "PlcControlParameter", common.PlcControlParameterMinLen)
	}
	p := &PlcControlParameter{
		FunctionCode:         common.FunctionCode(bytes[0]),
		UnknownBytes:         bytes[1:8],
		ParameterBlockLength: binary.BigEndian.Uint16(bytes[8:]),
	}
	if p.ParameterBlockLength != 0 {
		p.ParameterBlock = NewPlcControlStringParamBlock(string(bytes[10 : 10+p.ParameterBlockLength]))
	}
	p.LengthPart = bytes[10+p.ParameterBlockLength]
	if p.LengthPart != 0 {
		p.PiService = string(bytes[11 : 11+p.LengthPart])
	}
	return p, nil
}

func (p *PlcControlParameter) Len() int {
	return common.PlcControlParameterMinLen + int(p.ParameterBlockLength) + len(p.PiService)
}

func (p *PlcControlParameter) ToBytes() []byte {
	res := make([]byte, 0, p.Len())
	res = append(res, byte(p.FunctionCode))
	res = append(res, p.UnknownBytes...)
	res = append(res, util.NumberToBytes(p.ParameterBlockLength)...)
	res = append(res, p.ParameterBlock.ToBytes()...)
	res = append(res, p.LengthPart)
	res = append(res, []byte(p.PiService)...)
	return res
}

type PlcControlAckParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// UnknownByte 未知字节
	// 字节大小：1
	// 字节序数：1
	UnknownByte byte
}

func NewPlcControlAckParameter() *PlcControlAckParameter {
	return &PlcControlAckParameter{
		FunctionCode: common.FcControl,
		UnknownByte:  0x00,
	}
}

func PlcControlAckParameterFromBytes(bytes []byte) (*PlcControlAckParameter, error) {
	if len(bytes) < common.PlcControlAckParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "PlcControlAckParameter", common.PlcControlAckParameterLen)
	}
	return &PlcControlAckParameter{
		FunctionCode: common.FunctionCode(bytes[0]),
		UnknownByte:  bytes[1],
	}, nil
}

func (p *PlcControlAckParameter) Len() int {
	return common.PlcControlAckParameterLen
}

func (p *PlcControlAckParameter) ToBytes() []byte {
	res := make([]byte, 0, p.Len())
	res = append(res, byte(p.FunctionCode), p.UnknownByte)
	return res
}

type SetupComParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// Reserved 预留
	// 字节大小：1
	// 字节序数：1
	Reserved byte
	// MaxAmqCaller Ack队列的大小（主叫）（大端）
	// 字节大小：2
	// 字节序数：2-3
	MaxAmqCaller uint16
	// MaxAmqCallee Ack队列的大小（被叫）（大端）
	// 字节大小：2
	// 字节序数：4-5
	MaxAmqCallee uint16
	// PduLength PDU长度（大端）
	// 字节大小：2
	// 字节序数：6-7
	PduLength uint16
}

// NewSetupComParameter 创建默认的设置通信参数，默认最大PDU长度240
func NewSetupComParameter(pduLength uint16) *SetupComParameter {
	return &SetupComParameter{
		FunctionCode: common.FcSetupCom,
		Reserved:     0x00,
		MaxAmqCaller: 1,
		MaxAmqCallee: 1,
		PduLength:    pduLength,
	}
}

func SetupComParameterFromBytes(bytes []byte) (*SetupComParameter, error) {
	if len(bytes) < common.SetupComParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "SetupComParameter", common.SetupComParameterLen)
	}
	return &SetupComParameter{
		FunctionCode: common.FunctionCode(bytes[0]),
		Reserved:     bytes[1],
		MaxAmqCaller: binary.BigEndian.Uint16(bytes[2:]),
		MaxAmqCallee: binary.BigEndian.Uint16(bytes[4:]),
		PduLength:    binary.BigEndian.Uint16(bytes[6:]),
	}, nil
}

func (s *SetupComParameter) Len() int {
	return common.SetupComParameterLen
}

func (s *SetupComParameter) ToBytes() []byte {
	res := make([]byte, 0, s.Len())
	res = append(res, byte(s.FunctionCode), s.Reserved)
	res = append(res, util.NumberToBytes(s.MaxAmqCaller)...)
	res = append(res, util.NumberToBytes(s.MaxAmqCallee)...)
	res = append(res, util.NumberToBytes(s.PduLength)...)
	return res
}

type UploadParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
	// ErrorCode 未知字节
	// 字节大小：2
	// 字节序数：2-3
	ErrorCode []byte
	// Id 下载的Id，4个字节
	// 字节大小：4
	// 字节序数：4-7
	Id uint32
}

func NewUploadParameter(id uint32) *UploadParameter {
	return &UploadParameter{
		FunctionCode:      common.FcUpload,
		MoreDataFollowing: false,
		ErrorStatus:       false,
		ErrorCode:         []byte{0x00, 0x00},
		Id:                id,
	}
}

func UploadParameterFromBytes(bytes []byte) (*UploadParameter, error) {
	if len(bytes) < common.UploadParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "UploadParameter", common.UploadParameterLen)
	}
	return &UploadParameter{
		FunctionCode:      common.FunctionCode(bytes[0]),
		MoreDataFollowing: util.GetBoolAt(bytes[1], 0),
		ErrorStatus:       util.GetBoolAt(bytes[1], 1),
		ErrorCode:         bytes[2:4],
		Id:                binary.BigEndian.Uint32(bytes[4:8]),
	}, nil
}

func (u *UploadParameter) Len() int {
	return common.UploadParameterLen
}

func (u *UploadParameter) ToBytes() []byte {
	res := make([]byte, 0, u.Len())
	res = append(res, byte(u.FunctionCode), util.SetBoolAt(0x00, 0, u.MoreDataFollowing)|util.SetBoolAt(0x00, 1, u.ErrorStatus))
	res = append(res, u.ErrorCode...)
	res = append(res, util.NumberToBytes(u.Id)...)
	return res
}

type StartUploadParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
	// ErrorCode 未知字节
	// 字节大小：2
	// 字节序数：2-3
	ErrorCode []byte
	// Id 下载的Id，4个字节
	// 字节大小：4
	// 字节序数：4-7
	Id uint32
	// FileNameLength 文件名长度
	// 字节大小：1
	// 字节序数：8
	FileNameLength uint8
	// FileId 文件id
	// 字节大小：1
	// 字节序数：9
	FileId uint8
	// BlockType 数据块类型
	// 字节大小：2
	// 字节序数：10-11
	BlockType common.BlockType
	// BlockNumber 数据块编号，范围00000-99999
	// 字节大小：5
	// 字节序数：12-16
	BlockNumber int
	// DestinationFileSystem 目标文件系统
	// 字节大小：1
	// 字节序数：17
	DestinationFileSystem common.DestinationFileSystem
}

func NewStartUploadParameter(bt common.BlockType, dfs common.DestinationFileSystem, blockNumber int) *StartUploadParameter {
	return &StartUploadParameter{
		FunctionCode:          common.FcDownload,
		MoreDataFollowing:     false,
		ErrorStatus:           false,
		ErrorCode:             []byte{0x01, 0x00},
		Id:                    0,
		FileNameLength:        9,
		FileId:                byte('_'),
		BlockType:             bt,
		BlockNumber:           blockNumber,
		DestinationFileSystem: dfs,
	}
}

func StartUploadParameterFromBytes(bytes []byte) (*StartUploadParameter, error) {
	if len(bytes) < common.StartUploadParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "StartUploadParameter", common.StartUploadParameterLen)
	}
	return &StartUploadParameter{
		FunctionCode:          common.FunctionCode(bytes[0]),
		MoreDataFollowing:     util.GetBoolAt(bytes[1], 0),
		ErrorStatus:           util.GetBoolAt(bytes[1], 1),
		ErrorCode:             bytes[2:4],
		Id:                    binary.BigEndian.Uint32(bytes[4:8]),
		FileNameLength:        bytes[8],
		FileId:                bytes[9],
		BlockType:             common.BlockType(binary.BigEndian.Uint16(bytes[10:])),
		BlockNumber:           cast.ToInt(string(bytes[12:17])),
		DestinationFileSystem: common.DestinationFileSystem(bytes[17]),
	}, nil
}

func (s *StartUploadParameter) Len() int {
	return common.StartUploadParameterLen
}

func (s *StartUploadParameter) ToBytes() []byte {
	res := make([]byte, 0, s.Len())
	res = append(res, byte(s.FunctionCode), util.SetBoolAt(0x00, 0, s.MoreDataFollowing)|util.SetBoolAt(0x00, 1, s.ErrorStatus))
	res = append(res, s.ErrorCode...)
	res = append(res, util.NumberToBytes(s.Id)...)
	res = append(res, s.FileNameLength, s.FileId)
	res = append(res, util.NumberToBytes(uint16(s.BlockType))...)
	res = append(res, []byte(fmt.Sprintf("%05d", s.BlockNumber))...)
	res = append(res, byte(s.DestinationFileSystem))
	return res
}

type StartUploadAckParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
	// ErrorCode 未知字节
	// 字节大小：2
	// 字节序数：2-3
	ErrorCode []byte
	// Id 下载的Id，4个字节
	// 字节大小：4
	// 字节序数：4-7
	Id uint32
	// BlockLengthStringLength 即自此之后的数据长度
	// 字节大小：1
	// 字节序数：8
	BlockLengthStringLength uint8
	// BlockLength 到尾完整上传快的长度（以字节为单位）
	// 可以拆分为多个PDU
	// 字节大小：7
	// 字节序数：9-15
	BlockLength int
}

func (s *StartUploadAckParameter) Len() int {
	return common.StartUploadAckParameterLen
}

func (s *StartUploadAckParameter) ToBytes() []byte {
	res := make([]byte, 0, s.Len())
	res = append(res, byte(s.FunctionCode), util.SetBoolAt(0x00, 0, s.MoreDataFollowing)|util.SetBoolAt(0x00, 1, s.ErrorStatus))
	res = append(res, s.ErrorCode...)
	res = append(res, util.NumberToBytes(s.Id)...)
	res = append(res, s.BlockLengthStringLength)
	res = append(res, []byte(fmt.Sprintf("%07d", s.BlockLength))...)
	return res
}

func StartUploadAckParameterFromBytes(bytes []byte) (*StartUploadAckParameter, error) {
	if len(bytes) < common.StartUploadAckParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "StartUploadAckParameter", common.StartUploadAckParameterLen)
	}
	return &StartUploadAckParameter{
		FunctionCode:            common.FunctionCode(bytes[0]),
		MoreDataFollowing:       util.GetBoolAt(bytes[1], 0),
		ErrorStatus:             util.GetBoolAt(bytes[1], 1),
		ErrorCode:               bytes[2:4],
		Id:                      binary.BigEndian.Uint32(bytes[4:8]),
		BlockLengthStringLength: bytes[8],
		BlockLength:             cast.ToInt(string(bytes[9:16])),
	}, nil
}

type UploadAckParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
}

func NewUploadAckParameter() *UploadAckParameter {
	return &UploadAckParameter{
		FunctionCode:      common.FcUpload,
		MoreDataFollowing: false,
		ErrorStatus:       false,
	}
}

func (u *UploadAckParameter) Len() int {
	return common.UploadAckParameterLen
}

func (u *UploadAckParameter) ToBytes() []byte {
	res := make([]byte, 0, u.Len())
	res = append(res, byte(u.FunctionCode), util.SetBoolAt(0x00, 0, u.MoreDataFollowing)|util.SetBoolAt(0x00, 1, u.ErrorStatus))
	return res
}

func UploadAckParameterFromBytes(bytes []byte) (*UploadAckParameter, error) {
	if len(bytes) < common.UploadAckParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "UploadAckParameter", common.UploadAckParameterLen)
	}
	return &UploadAckParameter{
		FunctionCode:      common.FunctionCode(bytes[0]),
		MoreDataFollowing: util.GetBoolAt(bytes[1], 0),
		ErrorStatus:       util.GetBoolAt(bytes[1], 1),
	}, nil
}

type EndUploadParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
	// ErrorCode 未知字节
	// 字节大小：2
	// 字节序数：2-3
	ErrorCode []byte
	// Id 下载的Id，4个字节
	// 字节大小：4
	// 字节序数：4-7
	Id uint32
}

func NewEndUploadParameter(id uint32) *EndUploadParameter {
	return &EndUploadParameter{
		FunctionCode:      common.FcEndUpload,
		MoreDataFollowing: false,
		ErrorStatus:       false,
		ErrorCode:         []byte{0x00, 0x00},
		Id:                id,
	}
}

func EndUploadParameterFromBytes(bytes []byte) (*EndUploadParameter, error) {
	if len(bytes) < common.EndUploadParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "EndUploadParameter", common.EndUploadParameterLen)
	}
	return &EndUploadParameter{
		FunctionCode:      common.FunctionCode(bytes[0]),
		MoreDataFollowing: util.GetBoolAt(bytes[1], 0),
		ErrorStatus:       util.GetBoolAt(bytes[1], 1),
		ErrorCode:         bytes[2:4],
		Id:                binary.BigEndian.Uint32(bytes[4:8]),
	}, nil
}

func (e *EndUploadParameter) Len() int {
	return common.EndUploadParameterLen
}

func (e *EndUploadParameter) ToBytes() []byte {
	res := make([]byte, 0, e.Len())
	res = append(res, byte(e.FunctionCode), util.SetBoolAt(0x00, 0, e.MoreDataFollowing)|util.SetBoolAt(0x00, 1, e.ErrorStatus))
	res = append(res, e.ErrorCode...)
	res = append(res, util.NumberToBytes(e.Id)...)
	return res
}

type DownloadParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
	// ErrorCode 未知字节
	// 字节大小：2
	// 字节序数：2-3
	ErrorCode []byte
	// Id 下载的Id，4个字节
	// 字节大小：4
	// 字节序数：4-7
	Id uint32
	// FileNameLength 文件名长度
	// 字节大小：1
	// 字节序数：8
	FileNameLength uint8
	// FileId 文件id
	// 字节大小：1
	// 字节序数：9
	FileId uint8
	// BlockType 数据块类型
	// 字节大小：2
	// 字节序数：10-11
	BlockType common.BlockType
	// BlockNumber 数据块编号，范围00000-99999
	// 字节大小：5
	// 字节序数：12-16
	BlockNumber int
	// DestinationFileSystem 目标文件系统
	// 字节大小：1
	// 字节序数：17
	DestinationFileSystem common.DestinationFileSystem
}

func NewDownloadParameter(bt common.BlockType, dfs common.DestinationFileSystem, bn int, moreDataFollowing bool) *DownloadParameter {
	return &DownloadParameter{
		FunctionCode:          common.FcDownload,
		MoreDataFollowing:     moreDataFollowing,
		ErrorStatus:           false,
		ErrorCode:             []byte{0x01, 0x00},
		Id:                    0,
		FileNameLength:        9,
		FileId:                byte('_'),
		BlockType:             bt,
		BlockNumber:           bn,
		DestinationFileSystem: dfs,
	}
}

func DownloadParameterFromBytes(bytes []byte) (*DownloadParameter, error) {
	if len(bytes) < common.DownloadParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "DownloadParameter", common.DownloadParameterLen)
	}
	return &DownloadParameter{
		FunctionCode:          common.FunctionCode(bytes[0]),
		MoreDataFollowing:     util.GetBoolAt(bytes[1], 0),
		ErrorStatus:           util.GetBoolAt(bytes[1], 1),
		ErrorCode:             bytes[2:4],
		Id:                    binary.BigEndian.Uint32(bytes[4:8]),
		FileNameLength:        bytes[8],
		FileId:                bytes[9],
		BlockType:             common.BlockType(binary.BigEndian.Uint16(bytes[10:])),
		BlockNumber:           cast.ToInt(string(bytes[12:17])),
		DestinationFileSystem: common.DestinationFileSystem(bytes[17]),
	}, nil
}

func (d *DownloadParameter) Len() int {
	return common.DownloadParameterLen
}

func (d *DownloadParameter) ToBytes() []byte {
	res := make([]byte, 0, d.Len())
	res = append(res, byte(d.FunctionCode), util.SetBoolAt(0x00, 0, d.MoreDataFollowing)|util.SetBoolAt(0x00, 1, d.ErrorStatus))
	res = append(res, d.ErrorCode...)
	res = append(res, util.NumberToBytes(d.Id)...)
	res = append(res, d.FileNameLength, d.FileId)
	res = append(res, util.NumberToBytes(uint16(d.BlockType))...)
	res = append(res, []byte(fmt.Sprintf("%05d", d.BlockNumber))...)
	res = append(res, byte(d.DestinationFileSystem))
	return res
}

type StartDownloadParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
	// ErrorCode 未知字节
	// 字节大小：2
	// 字节序数：2-3
	ErrorCode []byte
	// Id 下载的Id，4个字节
	// 字节大小：4
	// 字节序数：4-7
	Id uint32
	// FileNameLength 文件名长度
	// 字节大小：1
	// 字节序数：8
	FileNameLength uint8
	// FileId 文件id
	// 字节大小：1
	// 字节序数：9
	FileId uint8
	// BlockType 数据块类型
	// 字节大小：2
	// 字节序数：10-11
	BlockType common.BlockType
	// BlockNumber 数据块编号，范围00000-99999
	// 字节大小：5
	// 字节序数：12-16
	BlockNumber int
	// DestinationFileSystem 目标文件系统
	// 字节大小：1
	// 字节序数：17
	DestinationFileSystem common.DestinationFileSystem
	// Part2Length 第二部分字符串长度
	// 字节大小：1
	// 字节序数：18
	Part2Length uint8
	// 未知字符
	// 字节大小：1
	// 字节序数：19
	UnknownChar byte
	// 未知字符
	// 字节大小：6
	// 字节序数：20-25
	LoadMemoryLength int
	// 未知字符
	// 字节大小：6
	// 字节序数：26-31
	McCodeLength int
}

func NewStartDownloadParameter(bt common.BlockType, dfs common.DestinationFileSystem, bn int, loadMemoryLength int, mcCodeLength int) *StartDownloadParameter {
	return &StartDownloadParameter{
		FunctionCode:          common.FcStartDownload,
		MoreDataFollowing:     false,
		ErrorStatus:           false,
		ErrorCode:             []byte{0x01, 0x00},
		Id:                    0,
		FileNameLength:        9,
		FileId:                byte('_'),
		BlockType:             bt,
		BlockNumber:           bn,
		DestinationFileSystem: dfs,
		Part2Length:           13,
		UnknownChar:           byte('1'),
		LoadMemoryLength:      loadMemoryLength,
		McCodeLength:          mcCodeLength,
	}
}

func StartDownloadParameterFromBytes(bytes []byte) (*StartDownloadParameter, error) {
	if len(bytes) < common.StartDownloadParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "StartDownloadParameter", common.StartDownloadParameterLen)
	}
	return &StartDownloadParameter{
		FunctionCode:          common.FunctionCode(bytes[0]),
		MoreDataFollowing:     util.GetBoolAt(bytes[1], 0),
		ErrorStatus:           util.GetBoolAt(bytes[1], 1),
		ErrorCode:             bytes[2:4],
		Id:                    binary.BigEndian.Uint32(bytes[4:8]),
		FileNameLength:        bytes[8],
		FileId:                bytes[9],
		BlockType:             common.BlockType(binary.BigEndian.Uint16(bytes[10:])),
		BlockNumber:           cast.ToInt(string(bytes[12:17])),
		DestinationFileSystem: common.DestinationFileSystem(bytes[17]),
		Part2Length:           bytes[18],
		UnknownChar:           bytes[19],
		LoadMemoryLength:      cast.ToInt(string(bytes[20:26])),
		McCodeLength:          cast.ToInt(string(bytes[26:32])),
	}, nil
}

func (s *StartDownloadParameter) Len() int {
	return common.StartDownloadParameterLen
}

func (s *StartDownloadParameter) ToBytes() []byte {
	res := make([]byte, 0, s.Len())
	res = append(res, byte(s.FunctionCode), util.SetBoolAt(0x00, 0, s.MoreDataFollowing)|util.SetBoolAt(0x00, 1, s.ErrorStatus))
	res = append(res, s.ErrorCode...)
	res = append(res, util.NumberToBytes(s.Id)...)
	res = append(res, s.FileNameLength, s.FileId)
	res = append(res, util.NumberToBytes(uint16(s.BlockType))...)
	res = append(res, []byte(fmt.Sprintf("%05d", s.BlockNumber))...)
	res = append(res, byte(s.DestinationFileSystem))
	res = append(res, s.Part2Length, s.UnknownChar)
	res = append(res, []byte(fmt.Sprintf("%06d", s.LoadMemoryLength))...)
	res = append(res, []byte(fmt.Sprintf("%06d", s.McCodeLength))...)
	return res
}

type EndDownloadParameter struct {
	// FunctionCode 功能码
	// 字节大小：1
	// 字节序数：0
	FunctionCode common.FunctionCode
	// MoreDataFollowing + ErrorStatus
	// 字节大小：1
	// 字节序数：1
	// MoreDataFollowing 后续是否还有更多数据
	MoreDataFollowing bool
	// ErrorStatus 错误状态
	ErrorStatus bool
	// ErrorCode 未知字节
	// 字节大小：2
	// 字节序数：2-3
	ErrorCode []byte
	// Id 下载的Id，4个字节
	// 字节大小：4
	// 字节序数：4-7
	Id uint32
	// FileNameLength 文件名长度
	// 字节大小：1
	// 字节序数：8
	FileNameLength uint8
	// FileId 文件id
	// 字节大小：1
	// 字节序数：9
	FileId uint8
	// BlockType 数据块类型
	// 字节大小：2
	// 字节序数：10-11
	BlockType common.BlockType
	// BlockNumber 数据块编号，范围00000-99999
	// 字节大小：5
	// 字节序数：12-16
	BlockNumber int
	// DestinationFileSystem 目标文件系统
	// 字节大小：1
	// 字节序数：17
	DestinationFileSystem common.DestinationFileSystem
}

func NewEndDownloadParameter(bt common.BlockType, dfs common.DestinationFileSystem, blockNumber int) *EndDownloadParameter {
	return &EndDownloadParameter{
		FunctionCode:          common.FcEndDownload,
		MoreDataFollowing:     false,
		ErrorStatus:           false,
		ErrorCode:             []byte{0x01, 0x00},
		Id:                    0,
		FileNameLength:        9,
		FileId:                byte('_'),
		BlockType:             bt,
		BlockNumber:           blockNumber,
		DestinationFileSystem: dfs,
	}
}

func EndDownloadParameterFromBytes(bytes []byte) (*EndDownloadParameter, error) {
	if len(bytes) < common.EndDownloadParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "EndDownloadParameter", common.EndDownloadParameterLen)
	}
	return &EndDownloadParameter{
		FunctionCode:          common.FunctionCode(bytes[0]),
		MoreDataFollowing:     util.GetBoolAt(bytes[1], 0),
		ErrorStatus:           util.GetBoolAt(bytes[1], 1),
		ErrorCode:             bytes[2:4],
		Id:                    binary.BigEndian.Uint32(bytes[4:8]),
		FileNameLength:        bytes[8],
		FileId:                bytes[9],
		BlockType:             common.BlockType(binary.BigEndian.Uint16(bytes[10:])),
		BlockNumber:           cast.ToInt(string(bytes[12:17])),
		DestinationFileSystem: common.DestinationFileSystem(bytes[17]),
	}, nil
}

func (e *EndDownloadParameter) Len() int {
	return common.EndDownloadParameterLen
}

func (e *EndDownloadParameter) ToBytes() []byte {
	res := make([]byte, 0, e.Len())
	res = append(res, byte(e.FunctionCode), util.SetBoolAt(0x00, 0, e.MoreDataFollowing)|util.SetBoolAt(0x00, 1, e.ErrorStatus))
	res = append(res, e.ErrorCode...)
	res = append(res, util.NumberToBytes(e.Id)...)
	res = append(res, e.FileNameLength, e.FileId)
	res = append(res, util.NumberToBytes(uint16(e.BlockType))...)
	res = append(res, []byte(fmt.Sprintf("%05d", e.BlockNumber))...)
	res = append(res, byte(e.DestinationFileSystem))
	return res
}

type UserdataParameter struct {
	// Header 参数头 固定0x000112
	// 字节大小：3
	// 字节序数：0-2
	Header []byte
	// ParameterLength 即自此之后参数长度
	// 字节大小：1
	// 字节序数：3
	ParameterLength uint8
	// Method 方法（request/response）
	// 字节大小：1
	// 字节序数：4
	Method common.Method
	// Type 类型
	// 字节大小: 1
	// 字节序数: 5
	Type common.FunctionGroup
	// SubFunction 请求子方法
	// 字节大小: 1
	// 字节序数: 6
	SubFunction byte
	// Sequence 顺序
	// 字节大小: 1
	// 字节序数: 7
	Sequence uint8
}

func NewCpuParameter(function common.CpuSubFunction) *UserdataParameter {
	return &UserdataParameter{
		Header:          []byte{0x00, 0x01, 0x12},
		ParameterLength: 4,
		Method:          common.MRequest,
		Type:            common.FgRequestCpuFunction,
		SubFunction:     byte(function),
		Sequence:        0,
	}
}

func NewBlockParameter(function common.BlockSubFunction) *UserdataParameter {
	return &UserdataParameter{
		Header:          []byte{0x00, 0x01, 0x12},
		ParameterLength: 4,
		Method:          common.MRequest,
		Type:            common.FgRequestBlockFunction,
		SubFunction:     byte(function),
		Sequence:        0,
	}
}

func NewClockParameter(function common.TimeSubFunction) *UserdataParameter {
	return &UserdataParameter{
		Header:          []byte{0x00, 0x01, 0x12},
		ParameterLength: 4,
		Method:          common.MRequest,
		Type:            common.FgRequestTimeFunction,
		SubFunction:     byte(function),
		Sequence:        0,
	}
}

func NewSecurityParameter(function common.SecuritySubFunction) *UserdataParameter {
	return &UserdataParameter{
		Header:          []byte{0x00, 0x01, 0x12},
		ParameterLength: 4,
		Method:          common.MRequest,
		Type:            common.FgRequestSecurity,
		SubFunction:     byte(function),
		Sequence:        0,
	}
}

func UserdataParameterFromBytes(bytes []byte) (*UserdataParameter, error) {
	if len(bytes) < common.UserdataParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "UserdataParameter", common.UserdataParameterLen)
	}
	return &UserdataParameter{
		Header:          bytes[:3],
		ParameterLength: bytes[3],
		Method:          common.Method(bytes[4]),
		Type:            common.FunctionGroup(bytes[5]),
		SubFunction:     bytes[6],
		Sequence:        bytes[7],
	}, nil
}

func (r *UserdataParameter) Len() int {
	return common.UserdataParameterLen
}

func (r *UserdataParameter) ToBytes() []byte {
	res := make([]byte, 0, r.Len())
	res = append(res, r.Header...)
	res = append(res, r.ParameterLength, byte(r.Method), byte(r.Type), r.SubFunction, r.Sequence)
	return res
}

type UserdataAckParameter struct {
	// Header 参数头 固定0x000112
	// 字节大小：3
	// 字节序数：0-2
	Header []byte
	// ParameterLength 即自此之后参数长度
	// 字节大小：1
	// 字节序数：3
	ParameterLength uint8
	// Method 方法（request/response）
	// 字节大小：1
	// 字节序数：4
	Method common.Method
	// Type 类型
	// 字节大小: 1
	// 字节序数: 5
	Type common.FunctionGroup
	// SubFunction 请求子方法
	// 字节大小: 1
	// 字节序数: 6
	SubFunction byte
	// Sequence 顺序
	// 字节大小: 1
	// 字节序数: 7
	Sequence uint8
	// TpduNumber TPDU编号
	// 字节大小：1，后面7位
	// 字节序数：8
	TpduNumber byte
	// LastDataUnit 是否最后一个数据单元
	// 字节大小：1，最高位，7位
	// 字节序数：9
	LastDataUnit byte
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

func UserdataAckParameterFromBytes(bytes []byte) (*UserdataAckParameter, error) {
	if len(bytes) < common.UserdataAckParameterLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "UserdataAckParameter", common.UserdataAckParameterLen)
	}
	return &UserdataAckParameter{
		Header:          bytes[:3],
		ParameterLength: bytes[3],
		Method:          common.Method(bytes[4]),
		Type:            common.FunctionGroup(bytes[5]),
		SubFunction:     bytes[6],
		Sequence:        bytes[7],
		TpduNumber:      bytes[8] & 0x7F,
		LastDataUnit:    bytes[9],
		ErrorClass:      bytes[10],
		ErrorCode:       bytes[10:12],
	}, nil
}

func (u *UserdataAckParameter) Len() int {
	return common.UserdataAckParameterLen
}

func (u *UserdataAckParameter) ToBytes() []byte {
	res := make([]byte, 0, u.Len())
	res = append(res, u.Header...)
	res = append(res, u.ParameterLength, byte(u.Method), byte(u.Type), u.SubFunction, u.Sequence, u.TpduNumber, u.LastDataUnit)
	res = append(res, u.ErrorCode...)
	return res
}
