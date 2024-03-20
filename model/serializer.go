// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"github.com/shiyuecamus/gs7/common"
)

func buildCotp(bytes []byte) (cotp common.COTP, err error) {
	if len(bytes) < 2 {
		err = common.ErrorWithCode(common.ErrModelFromBytes, "COTP", 2)
		return
	}
	switch pduType := common.PduType(bytes[1]); pduType {
	// TODO: The reject message structure is unknown.
	// TODO: If you understand.please modify the corresponding logic.
	// TODO: This is serialization based on the structure of the connection
	case common.PtConnectRequest, common.PtConnectConfirm, common.PtDisconnectRequest, common.PtDisconnectConfirm, common.PtReject:
		return COTPConnectionFromBytes(bytes)
	case common.PtData:
		return COTPDataFromBytes(bytes)
	default:
		err = common.ErrorWithCode(common.ErrTypeNotResolved, "pdu type", pduType)
		return
	}
}

func buildHeader(bytes []byte) (header common.Header, err error) {
	if len(bytes) < 2 {
		err = common.ErrorWithCode(common.ErrModelFromBytes, "Header", 2)
		return
	}
	messageType := common.MessageType(bytes[1])
	switch messageType {
	case common.MtJob:
		return RequestHeaderFromBytes(bytes)
	case common.MtAck, common.MtAckData:
		return AckHeaderFromBytes(bytes)
	case common.MtUserData:
		return RequestHeaderFromBytes(bytes)
	default:
		err = common.ErrorWithCode(common.ErrTypeNotResolved, "message type", messageType)
		return
	}
}

func buildParameter(bytes []byte, header common.Header) (parameter common.Parameter, err error) {
	switch mt := header.GetMessageType(); mt {
	case common.MtUserData:
		return UserdataAckParameterFromBytes(bytes)
	default:
		if len(bytes) < 1 {
			err = common.ErrorWithCode(common.ErrModelFromBytes, "Parameter", 1)
			return
		}
		functionCode := common.FunctionCode(bytes[0])
		switch functionCode {
		case common.FcCpuService:
			return
		case common.FcRead, common.FcWrite:
			return ReadWriteParameterFromBytes(bytes)
		case common.FcStartDownload:
			if mt == common.MtAckData {
				parameter = NewStandardParameter(functionCode)
				return
			}
			return StartDownloadParameterFromBytes(bytes)
		case common.FcDownload:
			if mt == common.MtAckData {
				parameter = NewStandardParameter(functionCode)
				return
			}
			return DownloadParameterFromBytes(bytes)
		case common.FcEndDownload:
			if mt == common.MtAckData {
				parameter = NewStandardParameter(functionCode)
				return
			}
			return EndDownloadParameterFromBytes(bytes)
		case common.FcStartUpload:
			if mt == common.MtAckData {
				return StartUploadAckParameterFromBytes(bytes)
			}
			return StartUploadParameterFromBytes(bytes)
		case common.FcUpload:
			if mt == common.MtAckData {
				return UploadAckParameterFromBytes(bytes)
			}
			return UploadParameterFromBytes(bytes)
		case common.FcEndUpload:
			if mt == common.MtAckData {
				parameter = NewStandardParameter(functionCode)
				return
			}
			return EndUploadParameterFromBytes(bytes)
		case common.FcControl:
			if mt == common.MtAck {
				return PlcControlAckParameterFromBytes(bytes)
			} else if mt == common.MtAckData {
				parameter = NewStandardParameter(functionCode)
				return
			}
			return PlcControlParameterFromBytes(bytes)
		case common.FcStop:
			if mt == common.MtAckData {
				parameter = NewStandardParameter(functionCode)
				return
			}
			return PlcStopParameterFromBytes(bytes)
		case common.FcSetupCom:
			return SetupComParameterFromBytes(bytes)
		default:
			err = common.ErrorWithCode(common.ErrTypeNotResolved, "function code", functionCode)
			return
		}
	}
}

func buildDatum(bytes []byte, header common.Header, code common.FunctionCode, group common.FunctionGroup, subFunc byte) (datum common.Datum, err error) {
	switch mt := header.GetMessageType(); mt {
	case common.MtUserData:
		switch group {
		case common.FgResponseCpuFunction:
			switch subFunc {
			case byte(common.CsfReadSzl):
				return ReadSzlAckDatumFromBytes(bytes)
			default:
				err = common.ErrorWithCode(common.ErrTypeNotResolved, "sub function", subFunc)
				return
			}
		case common.FgResponseBlockFunction:
			switch subFunc {
			case byte(common.BsfListBlock):
				return BlockListAckDatumFromBytes(bytes)
			case byte(common.BsfListBlockOfType):
				return BlockListTypeAckDatumFromBytes(bytes)
			case byte(common.BsfBlockInfo):
				return BlockInfoAckDatumFromBytes(bytes)
			default:
				err = common.ErrorWithCode(common.ErrTypeNotResolved, "sub function", subFunc)
				return
			}
		case common.FgResponseTimeFunction:
			switch subFunc {
			case byte(common.TsfReadClock):
				return ClockAckDatumFromBytes(bytes)
			case byte(common.TsfSetClock):
				return UserdataDatumFromBytes(bytes)
			default:
				err = common.ErrorWithCode(common.ErrTypeNotResolved, "sub function", subFunc)
				return
			}
		case common.FgResponseSecurity:
			switch subFunc {
			case byte(common.SsfSetPassword):
				return UserdataDatumFromBytes(bytes)
			case byte(common.SsfClearPassword):
				return UserdataDatumFromBytes(bytes)
			default:
				err = common.ErrorWithCode(common.ErrTypeNotResolved, "sub function", subFunc)
				return
			}
		default:
			err = common.ErrorWithCode(common.ErrTypeNotResolved, "function group", group)
			return
		}
	default:
		switch code {
		case common.FcRead, common.FcWrite:
			return ReadWriteDatumFromBytes(bytes, mt, code)
		case common.FcDownload, common.FcUpload:
			return UpDownloadDatumFromBytes(bytes)
		default:
			err = common.ErrorWithCode(common.ErrTypeNotResolved, "code", code)
			return
		}
	}
}
