// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"github.com/shiyuecamus/gs7/common"
)

type Catalog struct {
	OrderCode string
	Version   string
}

func (c Catalog) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// UnitInfo unit info
type UnitInfo struct {
	ModuleTypeName string
	SerialNumber   string
	ASName         string
	Copyright      string
	ModuleName     string
}

func (c UnitInfo) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// CommunicationInfo communication info
type CommunicationInfo struct {
	MaxPduLength   int
	MaxConnections int
	MaxMpiRate     int
	MaxBusRate     int
}

func (c CommunicationInfo) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// ProtectionInfo protection info
type ProtectionInfo struct {
	Level           uint16
	ParameterLevel  common.ParameterProtectionLevel
	CpuLevel        common.CpuProtectionLevel
	SelectorSetting common.SelectorSetting
	StartupSwitch   common.StartupSwitch
}

func (c ProtectionInfo) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

type PlcStatus byte

const (
	PsUnknown PlcStatus = 0x00
	PsRun               = 0x08
	PsStop              = 0x04
)

func (c PlcStatus) String() string {
	switch c {
	case PsRun:
		return "RUN"
	case PsStop:
		return "STOP"
	default:
		return "Unknown"
	}
}
