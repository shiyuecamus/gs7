// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package common

const (
	TpktLen                    int = 4
	CotpDataLen                    = 3
	CotpConnectionLen              = 18
	RequestHeaderLen               = 10
	AckHeaderLen                   = 12
	ReturnItemLen                  = 1
	NckRequestItemLen              = 10
	StandardRequestItemLen         = 12
	StandardParameterLen           = 1
	ReadWriteParameterMinLen       = 2
	PlcStopParameterMinLen         = 7
	PlcControlParameterMinLen      = 11
	PlcControlAckParameterLen      = 2
	SetupComParameterLen           = 8
	DownloadParameterLen           = 18
	EndDownloadParameterLen        = 18
	StartDownloadParameterLen      = 32
	UploadParameterLen             = 8
	UploadAckParameterLen          = 2
	EndUploadParameterLen          = 8
	StartUploadParameterLen        = 18
	StartUploadAckParameterLen     = 16
	UserdataParameterLen           = 8
	UserdataAckParameterLen        = 12
	UpDownloadDatumMinLen          = 4
	ReadSzlAckDatumMinLen          = 4
	BlockAckDatumMinLen            = 4
	DataItemMinLen                 = 4
	ReadSzlDatumLen                = 8
	BlockListTypeDatumLen          = 6
	BlockInfoDatumLen              = 12
	ClockReadAckDatumLen           = 14
	SetPasswordDatumLen            = 12
	UserdataDatumLen               = 4
)
