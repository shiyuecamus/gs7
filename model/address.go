// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"gs7/common"
	"regexp"
	"strconv"
	"strings"
)

func ParseAddress(address string) (requestItem common.RequestItem, err error) {
	if address == "" {
		err = common.ErrorWithCode(common.ErrAddressEmpty)
		return
	}
	address = strings.ToUpper(address)
	address = strings.Replace(address, " ", "", -1)

	split := strings.Split(address, ".")
	var area common.AreaType
	var variableType common.ParamVariableType
	var dbNumber, byteAddress, bitAddress int
	area, err = parseArea(split)
	if err != nil {
		return
	}
	variableType, err = parseVariableType(split)
	if err != nil {
		return
	}
	dbNumber, err = parseDbNumber(split)
	if err != nil {
		return
	}
	byteAddress, err = parseByteAddress(split)
	if err != nil {
		return
	}
	bitAddress, err = parseBitAddress(split, variableType)
	if err != nil {
		return
	}

	requestItem = NewStandardRequestItem(area, dbNumber, variableType, byteAddress, bitAddress, 1)
	return
}

func parseVariableType(split []string) (variableType common.ParamVariableType, err error) {
	switch one := split[0][:1]; one {
	case "T":
		variableType = common.PvtTimer
		return
	case "C":
		variableType = common.PvtCounter
		return
	case "D":
		if len(split) < 2 {
			err = common.ErrorWithCode(common.ErrAddressInvalid)
			return
		}
		return extractVariableType(split[1], true)
	default:
		return extractVariableType(split[0], false)
	}
}

func parseBitAddress(split []string, variableType common.ParamVariableType) (int, error) {
	switch one := split[0][:1]; one {
	case "D":
		if len(split) == 3 && variableType == common.PvtBit {
			return extractNumber(split[2])
		}
		return 0, nil
	default:
		if len(split) == 2 && variableType == common.PvtBit {
			return extractNumber(split[2])
		}
		return 0, nil
	}
}

func parseByteAddress(split []string) (int, error) {
	switch one := split[0][:1]; one {
	case "D":
		if len(split) >= 2 {
			return extractNumber(split[1])
		}
		return 0, common.ErrorWithCode(common.ErrAddressInvalid)
	default:
		return 0, nil
	}
}

func parseDbNumber(split []string) (int, error) {
	switch one := split[0][:1]; one {
	case "D":
		return extractNumber(split[0])
	case "V":
		return 1, nil
	default:
		return 0, nil
	}
}

func parseArea(split []string) (common.AreaType, error) {
	switch one := split[0][:1]; one {
	case "I":
		return common.AtInputs, nil
	case "Q":
		return common.AtOutputs, nil
	case "M":
		return common.AtFlags, nil
	case "V", "D":
		return common.AtDataBlocks, nil
	case "T":
		return common.AtTimers, nil
	case "C":
		return common.AtCounters, nil
	default:
		return 0, common.ErrorWithCode(common.ErrAddressInvalid)
	}
}

func extractVariableType(src string, isDb bool) (variableType common.ParamVariableType, err error) {
	re, err := regexp.Compile("\\D")
	if err != nil {
		return 0, err
	}
	all := re.FindAllString(src, -1)
	t := strings.Join(all, "")
	if !isDb {
		if len(t) < 2 {
			return 0, common.ErrorWithCode(common.ErrAddressInvalid)
		}
		t = t[1:]
	} else if len(t) < 1 {
		return 0, common.ErrorWithCode(common.ErrAddressInvalid)
	}
	switch t {
	case "X", "BIT":
		variableType = common.PvtBit
		return
	case "B", "BYTE":
		variableType = common.PvtByte
		return
	case "C", "CHAR":
		variableType = common.PvtChar
		return
	case "DW", "DWORD":
		variableType = common.PvtDWord
		return
	case "W", "WORD":
		variableType = common.PvtWord
		return
	case "DI", "DINT":
		variableType = common.PvtDInt
		return
	case "I", "INT":
		variableType = common.PvtInt
		return
	case "D", "DATE":
		variableType = common.PvtDate
		return
	case "DT", "DATETIME":
		variableType = common.PvtDateTime
		return
	case "DTL", "DATETIMELONG":
		variableType = common.PvtDTL
		return
	case "T", "TIME":
		variableType = common.PvtTime
		return
	case "ST", "STIME":
		variableType = common.PvtS5Time
		return
	case "TOD", "TIMEOFDAY":
		variableType = common.PvtTimeOfDay
		return
	case "R", "REAL":
		variableType = common.PvtReal
		return
	case "S", "STRING":
		variableType = common.PvtString
		return
	case "WS", "WSTRING":
		variableType = common.PvtWString
		return
	default:
		err = common.ErrorWithCode(common.ErrAddressInvalid)
	}
	return
}

func extractNumber(src string) (int, error) {
	re, err := regexp.Compile("\\D")
	if err != nil {
		return 0, err
	}
	number := re.ReplaceAllString(src, "")
	return strconv.Atoi(strings.TrimSpace(number))
}
