// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package core

import "github.com/spf13/cast"

// PlcControlInsertParamBlock PLC控制参数块，插入功能
type PlcControlInsertParamBlock struct {
	// UnknownByte 未知字节，固定0x00
	UnknownByte byte
	// FileNames 文件名
	FileNames []string
}

func NewPlcControlInsertParamBlock(filenames []string) *PlcControlInsertParamBlock {
	return &PlcControlInsertParamBlock{
		UnknownByte: 0x00,
		FileNames:   filenames,
	}
}

func (p *PlcControlInsertParamBlock) Len() int {
	sum := 2
	for _, name := range p.FileNames {
		sum += len(name)
	}
	return sum
}

func (p *PlcControlInsertParamBlock) ToBytes() []byte {
	res := make([]byte, 0, p.Len())
	res = append(res, cast.ToUint8(len(p.FileNames)))
	for _, name := range p.FileNames {
		res = append(res, []byte(name)...)
	}
	return res
}

type PlcControlStringParamBlock struct {
	ParamBlock string
}

func NewPlcControlStringParamBlock(paramBlock string) *PlcControlStringParamBlock {
	return &PlcControlStringParamBlock{
		ParamBlock: paramBlock,
	}
}

func (p *PlcControlStringParamBlock) Len() int {
	return len(p.ParamBlock)
}

func (p *PlcControlStringParamBlock) ToBytes() []byte {
	res := make([]byte, 0, p.Len())
	res = append(res, []byte(p.ParamBlock)...)
	return res
}
