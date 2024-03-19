// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package model

import (
	"gs7/common"
	"time"
)

type PDU struct {
	TPKT      common.TPKT
	COTP      common.COTP
	Header    common.Header
	Parameter common.Parameter
	Datum     common.Datum
}

// NewConnectRequest 创建连接请求
func NewConnectRequest(local uint16, remote uint16) *PDU {
	d := &PDU{
		TPKT: NewTPKT(),
		COTP: NewCOTPConnectionForRequest(local, remote),
	}
	d.SelfCheck()
	return d
}

// NewConnectDt 创建连接setup
func NewConnectDt(pduLength uint16, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewSetupComParameter(pduLength),
	}
	d.SelfCheck()
	return d
}

// NewReadRequest 创建默认读对象
func NewReadRequest(items []common.RequestItem, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewReqReadWriteParameter(common.FcRead, items),
	}
	d.SelfCheck()
	return d
}

// NewWriteRequest 创建默认写对象
func NewWriteRequest(reqItems []common.RequestItem, dateItems []common.ResponseItem, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewReqReadWriteParameter(common.FcWrite, reqItems),
		Datum:     NewReadWriteDatum(dateItems),
	}
	d.SelfCheck()
	return d
}

// NewHotRestart 创建热启动
func NewHotRestart(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewHotRestartPlcControlParameter(),
	}
	d.SelfCheck()
	return d
}

// NewColdRestart 创建冷启动
func NewColdRestart(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewColdRestartPlcControlParameter(),
	}
	d.SelfCheck()
	return d
}

// NewStopPlc 创建PLC停止命令
func NewStopPlc(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewPlcStopParameter(),
	}
	d.SelfCheck()
	return d
}

// NewCopyRamToRom 创建复制Ram到Rom的命令
func NewCopyRamToRom(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewCopyRamToRomPlcControlParameter(),
	}
	d.SelfCheck()
	return d
}

// NewCompress 创建压缩命令
func NewCompress(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewCompressPlcControlParameter(),
	}
	d.SelfCheck()
	return d
}

// NewInsert 创建插入文件指令
func NewInsert(bt common.BlockType, dfs common.DestinationFileSystem, bn int, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewInsertPlcControlParameter(bt, bn, dfs),
	}
	d.SelfCheck()
	return d
}

// NewStartDownload 创建开始下载
func NewStartDownload(bt common.BlockType, dfs common.DestinationFileSystem,
	bn int, loadMemoryLength int, mcCodeLength int, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewStartDownloadParameter(bt, dfs, bn, loadMemoryLength, mcCodeLength),
	}
	d.SelfCheck()
	return d
}

// NewDownload 创建下载中
func NewDownload(bt common.BlockType, dfs common.DestinationFileSystem,
	bn int, moreDataFollowing bool, bytes []byte, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewDownloadParameter(bt, dfs, bn, moreDataFollowing),
		Datum:     NewUpDownloadDatum(bytes),
	}
	d.SelfCheck()
	return d
}

// NewEndDownload 创建结束下载
func NewEndDownload(bt common.BlockType, dfs common.DestinationFileSystem, bn int, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewEndDownloadParameter(bt, dfs, bn),
	}
	d.SelfCheck()
	return d
}

// NewStartUpload 创建开始上传
func NewStartUpload(bt common.BlockType, dfs common.DestinationFileSystem, bn int, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewStartUploadParameter(bt, dfs, bn),
	}
	d.SelfCheck()
	return d
}

// NewUpload 创建上传中
func NewUpload(id uint32, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewUploadParameter(id),
	}
	d.SelfCheck()
	return d
}

// NewEndUpload 创建结束上传
func NewEndUpload(id uint32, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewRequestHeader(requestId),
		Parameter: NewEndUploadParameter(id),
	}
	d.SelfCheck()
	return d
}

// NewReadSzl 读取szl
func NewReadSzl(szlId uint16, szlIndex uint16, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewCpuParameter(common.CsfReadSzl),
		Datum:     NewReadSzlDatum(szlId, szlIndex),
	}
	d.SelfCheck()
	return d
}

// NewBlockList 查询块列表
func NewBlockList(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewBlockParameter(common.BsfListBlock),
		Datum:     NewUserdataDatum(),
	}
	d.SelfCheck()
	return d
}

// NewBlockListType 查询块列表类型
func NewBlockListType(bt common.BlockType, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewBlockParameter(common.BsfListBlockOfType),
		Datum:     NewBlockListTypeDatum(bt),
	}
	d.SelfCheck()
	return d
}

// NewBlockInfo 查询块信息
func NewBlockInfo(bt common.BlockType, dfs common.DestinationFileSystem, bn int, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewBlockParameter(common.BsfBlockInfo),
		Datum:     NewBlockInfoDatum(bt, dfs, bn),
	}
	d.SelfCheck()
	return d
}

// NewClockRead 读取时间
func NewClockRead(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewClockParameter(common.TsfReadClock),
		Datum:     NewUserdataDatum(),
	}
	d.SelfCheck()
	return d
}

// NewClockSet 设置时间
func NewClockSet(t time.Time, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewClockParameter(common.TsfSetClock),
		Datum:     NewClockAckDatum(t),
	}
	d.SelfCheck()
	return d
}

// NewSetPassword 设置密码
func NewSetPassword(pwd string, requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewSecurityParameter(common.SsfSetPassword),
		Datum:     NewSetPasswordDatum(pwd),
	}
	d.SelfCheck()
	return d
}

// NewClearPassword 清除密码
func NewClearPassword(requestId uint16) *PDU {
	d := &PDU{
		TPKT:      NewTPKT(),
		COTP:      NewCOTPData(),
		Header:    NewUserDataHeader(requestId),
		Parameter: NewSecurityParameter(common.SsfClearPassword),
		Datum:     NewUserdataDatum(),
	}
	d.SelfCheck()
	return d
}

func (d *PDU) Len() int {
	l := 0
	if d.TPKT != nil {
		l += d.TPKT.Len()
	}
	if d.COTP != nil {
		l += d.COTP.Len()
	}
	if d.Header != nil {
		l += d.Header.Len()
	}
	if d.Parameter != nil {
		l += d.Parameter.Len()
	}
	if d.Datum != nil {
		l += d.Datum.Len()
	}
	return l
}

func (d *PDU) ToBytes() []byte {
	res := make([]byte, 0, d.Len())
	if d.TPKT != nil {
		res = append(res, d.TPKT.ToBytes()...)
	}
	if d.COTP != nil {
		res = append(res, d.COTP.ToBytes()...)
	}
	if d.Header != nil {
		res = append(res, d.Header.ToBytes()...)
	}
	if d.Parameter != nil {
		res = append(res, d.Parameter.ToBytes()...)
	}
	if d.Datum != nil {
		res = append(res, d.Datum.ToBytes()...)
	}
	return res
}

func (d *PDU) SelfCheck() {
	if d.Header != nil {
		d.Header.SetDataLength(0)
		d.Header.SetParameterLength(0)
	}
	if d.Parameter != nil && d.Header != nil {
		d.Header.SetParameterLength(uint16(d.Parameter.Len()))
	}
	if d.Datum != nil && d.Header != nil {
		d.Header.SetDataLength(uint16(d.Datum.Len()))
	}
	if d.TPKT != nil {
		d.TPKT.SetLength(uint16(d.Len()))
	}
}

func DataFromBytes(bytes []byte) (d *PDU, err error) {
	d = &PDU{}
	// tpkt
	if len(bytes) < common.TpktLen {
		err = common.ErrorWithCode(common.ErrModelFromBytes, "TPKT", common.TpktLen)
		return
	}

	var tpkt common.TPKT
	tpkt, err = TPKTFromBytes(bytes[:common.TpktLen])
	if err != nil {
		return
	}
	d.TPKT = tpkt
	// cotp
	remain := bytes[common.TpktLen:]

	var cotp common.COTP
	cotp, err = buildCotp(remain)
	if err != nil {
		return
	}
	d.COTP = cotp

	if len(remain) == d.COTP.Len() {
		return d, nil
	}

	// header
	remain = remain[d.COTP.Len():]
	var header common.Header
	header, err = buildHeader(remain)
	if err != nil {
		return
	}
	d.Header = header

	var fc = common.FunctionCode(255)
	var fg = common.FunctionGroup(255)
	var subFunc = byte(255)
	// parameter
	if d.Header.GetParameterLength() > 0 {
		parameterBs := remain[d.Header.Len() : d.Header.Len()+int(d.Header.GetParameterLength())]
		fc = common.FunctionCode(parameterBs[0])

		// userdata handle
		if len(parameterBs) >= common.UserdataParameterLen {
			fg = common.FunctionGroup(parameterBs[5])
			subFunc = parameterBs[6]
		}
		var parameter common.Parameter
		parameter, err = buildParameter(parameterBs, d.Header)
		if err != nil {
			return
		}
		d.Parameter = parameter
	}

	if d.Header.GetDataLength() > 0 {
		dataBs := remain[d.Header.Len()+int(d.Header.GetParameterLength()):]
		var datum common.Datum
		datum, err = buildDatum(dataBs, d.Header, fc, fg, subFunc)
		if err != nil {
			return
		}
		d.Datum = datum
	}
	return d, nil
}

func (d *PDU) GetTPKT() common.TPKT {
	return d.TPKT
}

func (d *PDU) GetCOTP() common.COTP {
	return d.COTP
}

func (d *PDU) GetHeader() common.Header {
	return d.Header
}

func (d *PDU) GetParameter() common.Parameter {
	return d.Parameter
}

func (d *PDU) GetDatum() common.Datum {
	return d.Datum
}
