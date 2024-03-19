// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"encoding/binary"
	"gs7/common"
	"gs7/util"
)

type TPKT struct {
	// Version 版本号，常量0x03 <br>
	// 字节大小：1
	// 字节序数：0
	Version byte
	// Reserved 预留，默认值0x00
	// 字节大小：1
	// 字节序数：1
	Reserved byte
	// Length 长度，包括后面负载payload+版本号+预留+长度
	// 字节大小：2
	// 字节序数：2-3
	Length uint16
}

func NewTPKT() *TPKT {
	return &TPKT{
		Version:  0x03,
		Reserved: 0x00,
		Length:   0,
	}
}

func (t *TPKT) Len() int {
	return common.TpktLen
}

func (t *TPKT) ToBytes() []byte {
	res := make([]byte, 0, t.Len())
	res = append(res, t.Version, t.Reserved)
	res = append(res, util.NumberToBytes(t.Length)...)
	return res
}

func TPKTFromBytes(bytes []byte) (*TPKT, error) {
	if len(bytes) < common.TpktLen {
		return nil, common.ErrorWithCode(common.ErrModelFromBytes, "TPKT", common.TpktLen)
	}
	return &TPKT{
		Version:  bytes[0],
		Reserved: bytes[1],
		Length:   binary.BigEndian.Uint16(bytes[2:]),
	}, nil
}

func (t *TPKT) GetVersion() byte {
	return t.Version
}

func (t *TPKT) GetReserved() byte {
	return t.Reserved
}

func (t *TPKT) GetLength() uint16 {
	return t.Length
}

func (t *TPKT) SetVersion(version byte) {
	t.Version = version
}

func (t *TPKT) SetReserved(reserved byte) {
	t.Reserved = reserved
}

func (t *TPKT) SetLength(length uint16) {
	t.Length = length
}
