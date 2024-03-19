// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"gs7/common"
	"time"
)

type ListBlockInfo struct {
	Type  common.BlockType
	Count uint16
}

type ListBlockTypeInfo struct {
	Number   uint16
	Flags    uint8
	Language uint8
}

// BlockInfo Managed Block Info
type BlockInfo struct {
	BlockType        int
	BlockNumber      int
	Language         int
	Flags            int
	MC7CodeLength    int
	LengthLoadMemory int
	LocalDataLength  int
	SBBLength        int
	CheckSum         int
	Version          int
	CodeDate         time.Time
	InterfaceDate    time.Time
	Author           string
	Family           string
	Header           string
}
