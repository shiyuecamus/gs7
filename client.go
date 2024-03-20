// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/core"
	"github.com/shiyuecamus/gs7/logging"
	"github.com/shiyuecamus/gs7/util"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type client struct {
	m         *sync.RWMutex
	tcpClient *s7TcpClient
	cli       *gnet.Client
	conn      gnet.Conn
	logger    logging.Logger
	// timeout connect and read Timeout
	// default value 5s
	timeout time.Duration
	// autoReconnect reconnect on connection lost
	// default value false
	autoReconnect bool
	// reconnectInterval maximum connection reconnect interval.
	// each attempt will be multiplied by 2
	// default value 10s
	reconnectInterval time.Duration
	// maxReconnectTimes maximum reconnect times
	// if set to -1, infinite attempts will be made to reconnect
	// default value 5min
	maxReconnectTimes int
	// maxReconnectBackoff maximum reconnect wait time
	// default value 5min
	maxReconnectBackoff time.Duration
	// connectRetry automatically retry when attempting to connect failed
	// default value false
	connectRetry bool
	// retryInterval maximum connection retry interval.
	// each attempt will be multiplied by 2
	// default value 10s
	retryInterval time.Duration
	// maxRetries maximum connection retries
	// if set to -1, infinite attempts will be made to retry
	// default value 5
	maxRetries int
	// maxRetryBackoff maximum connection wait time
	// default value 5min
	maxRetryBackoff time.Duration

	plcType   common.PlcType
	host      string
	port      int
	rack      int
	slot      int
	pduLength int

	pduIndex uint32
	status   connectionStatus

	onConnected    func(c Client)
	onDisconnected func(c Client, err error)
}

func (c *client) init() *client {
	options := make([]gnet.Option, 0)
	options = append(options,
		gnet.WithLogger(c.logger),
		gnet.WithMulticore(true))
	tcpClient := newTcpClient(c.logger, c.timeout,
		c.tcpOnOpen, c.tcpOnClose, c.validate)
	cli, _ := gnet.NewClient(tcpClient, options...)
	_ = cli.Start()
	c.tcpClient = tcpClient
	c.cli = cli
	return c
}

func (c *client) ReadParsed(address string) *SingleParsedReadToken {
	token := NewToken(TtSingleParsedRead).(*SingleParsedReadToken)
	c.ReadBatchParsed([]string{address}).Async(func(v []any, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.v = v[0]
		token.flowComplete()
	})
	return token
}

func (c *client) ReadBatchParsed(addresses []string) *BatchParsedReadToken {
	token := NewToken(TtBatchParsedRead).(*BatchParsedReadToken)
	c.ReadBatchRaw(addresses).Async(func(v []RawInfo, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		res := make([]any, 0)
		var value any
		for _, info := range v {
			value, err = info.Parse()
			if err != nil {
				token.setError(err)
				return
			}
			res = append(res, value)
		}
		token.v = res
		token.flowComplete()
	})
	return token
}

func (c *client) ReadRaw(address string) *SingleRawReadToken {
	token := NewToken(TtSingleRawRead).(*SingleRawReadToken)
	c.ReadBatchRaw([]string{address}).Async(func(v []RawInfo, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.v = v[0]
		token.flowComplete()
	})
	return token
}

func (c *client) ReadBatchRaw(addresses []string) *BatchRawReadToken {
	token := NewToken(TtBatchRawRead).(*BatchRawReadToken)
	items, _, err := c.parseReadRequestItems(addresses)
	if err != nil {
		token.setError(err)
		return token
	}
	c.read(items).Async(func(v []*core.DataItem, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		res := make([]RawInfo, 0, len(v))
		for i, dataItem := range v {
			res = append(res, RawInfo{
				Value:   dataItem.Data,
				Type:    items[i].(*core.StandardRequestItem).VariableType,
				plcType: c.plcType,
			})
		}
		token.v = res
		token.flowComplete()
	})
	return token
}

func (c *client) WriteRaw(address string, data []byte) *SimpleToken {
	return c.WriteRawBatch([]string{address}, [][]byte{data})
}

func (c *client) WriteRawBatch(addresses []string, data [][]byte) *SimpleToken {
	requests, dataItems, err := c.parsesWriteRequestItems(addresses, data)
	if err != nil {
		token := NewToken(TtSimple).(*SimpleToken)
		token.setError(err)
		return token
	}
	return c.write(requests, dataItems)
}

func (c *client) BaseRead(area common.AreaType, dbNumber int, byteAddr int, bitAddr int, size int) *BaseReadToken {
	token := NewToken(TtBaseRead).(*BaseReadToken)
	item := core.NewStandardRequestItem(area, dbNumber, common.PvtByte, byteAddr, bitAddr, size)
	c.read([]common.RequestItem{item}).Async(func(v []*core.DataItem, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.v = v[0].Data
		token.flowComplete()
	})
	return token
}

func (c *client) BaseWrite(area common.AreaType, dbNumber int, byteAddr int, bitAddr int, data []byte) *SimpleToken {
	item := core.NewStandardRequestItem(area, dbNumber, common.PvtByte, byteAddr, bitAddr, len(data))
	dataItem := core.NewReqDataItem(data, item.VariableType.DataVariableType())
	return c.write([]common.RequestItem{item}, []common.ResponseItem{dataItem})
}

func (c *client) HotRestart() *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewHotRestart(c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) ColdRestart() *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewColdRestart(c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) StopPlc() *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewStopPlc(c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) CopyRamToRom() *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewCopyRamToRom(c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) Compress() *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewCompress(c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) InsertFile(bt common.BlockType, blockNumber int) *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewInsert(bt, common.DfsP, blockNumber, c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) UploadFile(bt common.BlockType, blockNumber int) *UploadToken {
	token := NewToken(TtUpload).(*UploadToken)
	c.send(core.NewStartUpload(bt, common.DfsA, blockNumber, c.GeneratePduNumber())).Async(func(v *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		parameter := v.GetParameter().(*core.StartUploadAckParameter)
		res := make([]byte, 0, parameter.BlockLength)
		ackParameter := core.NewUploadAckParameter()
		ackParameter.MoreDataFollowing = true
		var uploadAck *core.PDU
		for ackParameter.MoreDataFollowing {
			uploadToken := c.send(core.NewUpload(parameter.Id, c.GeneratePduNumber()))
			uploadAck, err = uploadToken.Wait()
			if err != nil {
				token.setError(err)
				return
			}
			ackParameter = uploadAck.GetParameter().(*core.UploadAckParameter)
			if ackParameter.ErrorStatus {
				err = common.ErrorWithCode(common.ErrCliUploadFailed)
				token.setError(err)
				return
			}
			datum := uploadAck.GetDatum().(*core.UpDownloadDatum)
			res = append(res, datum.Data...)
		}
		endToken := c.send(core.NewEndUpload(parameter.Id, c.GeneratePduNumber()))
		_, err = endToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		token.v = res
		token.flowComplete()
	})

	return token
}

func (c *client) DownloadFile(bytes []byte, bt common.BlockType, bn int, mC7CodeLength int) *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	total := len(bytes)
	c.send(core.NewStartDownload(bt, common.DfsP, bn, total, mC7CodeLength, c.GeneratePduNumber())).Async(func(v *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		sent := 0
		for sent < total {
			moreDataFollowing := total-sent > c.pduLength-32
			length := int(math.Min(float64(total-sent), float64(c.pduLength-32)))
			downloadToken := c.send(core.NewDownload(bt, common.DfsP, bn, moreDataFollowing, bytes[sent:sent+length], c.GeneratePduNumber()))
			_, err = downloadToken.Wait()
			if err != nil {
				token.setError(err)
				return
			}
			sent += length
		}
		endToken := c.send(core.NewEndDownload(bt, common.DfsP, bn, c.GeneratePduNumber()))
		_, err = endToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) GetSzlIds() *SzlIdsToken {
	token := NewToken(TtSzlIds).(*SzlIdsToken)
	go func() {
		szlToken := c.ReadSzl(0x0000, 0x0000)
		pdu, err := szlToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		datum := pdu.GetDatum().(*core.ReadSzlAckDatum)
		if datum.PartLength != 2 {
			token.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		res := make([]uint16, 0, datum.PartCount)
		for i := 0; i < int(datum.PartCount); i++ {
			res = append(res, binary.BigEndian.Uint16(datum.Parts[i]))
		}
		token.v = res
		token.flowComplete()
	}()
	return token
}

func (c *client) GetCatalog() *CatalogToken {
	token := NewToken(TtCatalog).(*CatalogToken)
	go func() {
		szlToken := c.ReadSzl(0x0011, 0x0000)
		pdu, err := szlToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		datum := pdu.GetDatum().(*core.ReadSzlAckDatum)
		if datum.PartCount < 3 || datum.PartLength != 28 {
			token.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		start := 2
		end := datum.PartLength - 6
		token.v = core.Catalog{
			OrderCode: strings.TrimSpace(string(datum.Parts[0][start:end])),
			Version: fmt.Sprintf("V%d.%d.%d",
				datum.Parts[2][25],
				datum.Parts[2][26],
				datum.Parts[2][27]),
		}
		token.flowComplete()
	}()
	return token
}

func (c *client) GetPlcStatus() *PlcStatusToken {
	token := NewToken(TtPlcStatus).(*PlcStatusToken)
	go func() {
		szlToken := c.ReadSzl(0x0024, 0x0000)
		pdu, err := szlToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		datum := pdu.GetDatum().(*core.ReadSzlAckDatum)
		if datum.PartCount < 0 || len(datum.Parts[0]) < 4 {
			token.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		token.v = core.PlcStatus(datum.Parts[0][3])
		token.flowComplete()
	}()
	return token
}

func (c *client) GetUnitInfo() *UnitInfoToken {
	token := NewToken(TtUnitInfo).(*UnitInfoToken)
	go func() {
		szlToken := c.ReadSzl(0x001C, 0x0000)
		pdu, err := szlToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		datum := pdu.GetDatum().(*core.ReadSzlAckDatum)
		if datum.PartCount < 6 || datum.PartLength != 34 {
			token.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		token.v = core.UnitInfo{
			ASName:         strings.TrimSpace(string(datum.Parts[0][2:26])),
			ModuleName:     strings.TrimSpace(string(datum.Parts[1][2:26])),
			Copyright:      strings.TrimSpace(string(datum.Parts[3][2:28])),
			SerialNumber:   strings.TrimSpace(string(datum.Parts[4][2:26])),
			ModuleTypeName: strings.TrimSpace(string(datum.Parts[5][2:26])),
		}
		token.flowComplete()
	}()
	return token
}

func (c *client) GetCommunicationInfo() *CommunicationInfoToken {
	token := NewToken(TtCommunicationInfo).(*CommunicationInfoToken)
	go func() {
		szlToken := c.ReadSzl(0x0131, 0x0000)
		pdu, err := szlToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		datum := pdu.GetDatum().(*core.ReadSzlAckDatum)
		if datum.PartCount < 0 || datum.PartLength != 34 {
			token.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		token.v = core.CommunicationInfo{
			MaxPduLength:   int(binary.BigEndian.Uint16(datum.Parts[0][2:])),
			MaxConnections: int(binary.BigEndian.Uint16(datum.Parts[0][4:])),
			MaxMpiRate:     int(binary.BigEndian.Uint16(datum.Parts[0][6:])),
			MaxBusRate:     int(binary.BigEndian.Uint16(datum.Parts[0][10:])),
		}
		token.flowComplete()
	}()
	return token
}

func (c *client) GetProtectionInfo() *ProtectionInfoToken {
	token := NewToken(TtProtectionInfo).(*ProtectionInfoToken)
	go func() {
		szlToken := c.ReadSzl(0x0232, 0x0004)
		pdu, err := szlToken.Wait()
		if err != nil {
			token.setError(err)
			return
		}
		datum := pdu.GetDatum().(*core.ReadSzlAckDatum)
		if datum.PartCount < 0 {
			token.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		token.v = core.ProtectionInfo{
			Level:           binary.BigEndian.Uint16(datum.Parts[0][2:]),
			ParameterLevel:  common.ParameterProtectionLevel(binary.BigEndian.Uint16(datum.Parts[0][4:])),
			CpuLevel:        common.CpuProtectionLevel(binary.BigEndian.Uint16(datum.Parts[0][6:])),
			SelectorSetting: common.SelectorSetting(binary.BigEndian.Uint16(datum.Parts[0][8:])),
			StartupSwitch:   common.StartupSwitch(binary.BigEndian.Uint16(datum.Parts[0][10:])),
		}
		token.flowComplete()
	}()
	return token
}

func (c *client) ReadSzl(szlId uint16, szlIndex uint16) *PduToken {
	t := NewToken(TtPdu).(*PduToken)
	c.send(core.NewReadSzl(szlId, szlIndex, c.GeneratePduNumber())).Async(func(v *core.PDU, err error) {
		if err != nil {
			t.setError(err)
			return
		}
		if _, ok := v.GetDatum().(*core.ReadSzlAckDatum); !ok {
			t.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		t.v = v
		t.flowComplete()
	})
	return t
}

func (c *client) BlockList() *BlockListToken {
	token := NewToken(TtBlockList).(*BlockListToken)
	c.send(core.NewBlockList(c.GeneratePduNumber())).Async(func(v *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.v = v.GetDatum().(*core.BlockListAckDatum).Blocks
		token.flowComplete()
	})
	return token
}

func (c *client) BlockListType(bt common.BlockType) *BlockListTypeToken {
	token := NewToken(TtBlockListType).(*BlockListTypeToken)
	c.send(core.NewBlockListType(bt, c.GeneratePduNumber())).Async(func(v *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.v = v.GetDatum().(*core.BlockListTypeAckDatum).Types
		token.flowComplete()
	})
	return token
}

func (c *client) BlockInfo(bt common.BlockType, bn int) *BlockInfoToken {
	token := NewToken(TtBlockInfo).(*BlockInfoToken)
	c.send(core.NewBlockInfo(bt, common.DfsA, bn, c.GeneratePduNumber())).Async(func(v *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		datum := v.GetDatum().(*core.BlockInfoAckDatum)
		codeDate := (int64(binary.BigEndian.Uint32(datum.CodeTimestamp[:4])) << 16) +
			int64(binary.BigEndian.Uint16(datum.CodeTimestamp[4:]))
		interfaceDate := (int64(binary.BigEndian.Uint32(datum.InterfaceTimestamp[:4])) << 16) +
			int64(binary.BigEndian.Uint16(datum.InterfaceTimestamp[4:]))
		token.v = core.BlockInfo{
			BlockType:        int(datum.BlockType),
			BlockNumber:      int(datum.BlockNumber),
			Language:         int(datum.Language),
			Flags:            int(datum.Flags),
			MC7CodeLength:    int(datum.MC7CodeLength),
			LengthLoadMemory: int(datum.LengthLoadMemory),
			LocalDataLength:  int(datum.LocalDataLength),
			SBBLength:        int(datum.SBBLength),
			CheckSum:         int(datum.CheckSum),
			Version:          int(datum.Version),
			CodeDate:         siemensTimestamp(codeDate),
			InterfaceDate:    siemensTimestamp(interfaceDate),
			Author:           strings.TrimSpace(string(datum.Auth)),
			Family:           strings.TrimSpace(string(datum.Family)),
			Header:           strings.TrimSpace(string(datum.Header)),
		}
		token.flowComplete()
	})
	return token
}

func (c *client) DBFill(dbNumber int, fillByte byte) *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.BlockInfo(common.DtDb, dbNumber).Async(func(v core.BlockInfo, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		data := make([]byte, v.MC7CodeLength)
		for i := 0; i < v.MC7CodeLength; i++ {
			data[i] = fillByte
		}
		c.BaseWrite(common.AtDataBlocks, dbNumber, 0, 0, data).Async(func(err error) {
			if err != nil {
				token.setError(err)
				return
			}
			token.flowComplete()
		})
	})
	return token
}

func (c *client) DBGet(dbNumber int) *BaseReadToken {
	token := NewToken(TtBaseRead).(*BaseReadToken)
	c.BlockInfo(common.DtDb, dbNumber).Async(func(v core.BlockInfo, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		c.BaseRead(common.AtDataBlocks, dbNumber, 0, 0, v.MC7CodeLength).Async(func(v []byte, err error) {
			if err != nil {
				token.setError(err)
				return
			}
			token.v = v
			token.flowComplete()
		})
	})
	return token
}

func (c *client) ClockRead() *ClockReadToken {
	token := NewToken(TtClockRead).(*ClockReadToken)
	c.send(core.NewClockRead(c.GeneratePduNumber())).Async(func(v *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		datum := v.GetDatum().(*core.ClockAckDatum)
		year, _ := strconv.Atoi(hex.EncodeToString([]byte{datum.Year1, datum.Year2}))
		month, _ := strconv.Atoi(hex.EncodeToString([]byte{datum.Month}))
		day, _ := strconv.Atoi(hex.EncodeToString([]byte{datum.Day}))
		hour, _ := strconv.Atoi(hex.EncodeToString([]byte{datum.Hour}))
		minute, _ := strconv.Atoi(hex.EncodeToString([]byte{datum.Minute}))
		second, _ := strconv.Atoi(hex.EncodeToString([]byte{datum.Second}))
		token.v = time.Date(year,
			time.Month(month),
			day,
			hour,
			minute,
			second,
			int(datum.MilliSecond)*int(time.Millisecond),
			time.UTC)
		token.flowComplete()
	})
	return token
}

func (c *client) ClockSet(t time.Time) *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewClockSet(t, c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) SetPassword(pwd string) *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	if len(pwd) > 8 {
		token.setError(common.ErrorWithCode(common.ErrPasswordLengthInvalid, 8))
		return token
	}
	c.send(core.NewSetPassword(pwd, c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) ClearPassword() *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	c.send(core.NewClearPassword(c.GeneratePduNumber())).Async(func(_ *core.PDU, err error) {
		if err != nil {
			token.setError(err)
			return
		}
		token.flowComplete()
	})
	return token
}

func (c *client) autoConnect(endpoint string) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", endpoint, c.timeout)
	backoff := c.retryInterval

	for retries := 1; c.maxRetries == -1 || retries <= c.maxRetries; retries++ {
		conn, err = net.Dial("tcp", endpoint)
		if err == nil {
			return conn, nil
		}

		c.logger.Debugf("retrying in %v, failed to connect to %s: %v", backoff, endpoint, err)
		if backoff < c.maxRetryBackoff {
			backoff *= 2
		}
		if retries != c.maxRetries {
			time.Sleep(backoff)
		}
	}
	err = common.ErrorWithCode(common.ErrTcpConnectWithAttempts, endpoint, c.maxRetries)
	return
}

func (c *client) isoConnect() *core.PDU {
	var local uint16 = 0x0100
	var remote uint16 = 0x0300
	switch c.plcType {
	case common.S200:
		local = 0x4D57
		remote = 0x4D57
		break
	case common.S200Smart:
		local = 0x1000
		remote = 0x0300
		break
	case common.S300, common.S400, common.S1200, common.S1500:
		remote += 0x20*uint16(c.rack) + uint16(c.slot)
		break
	case common.Sinumerik828d:
		local = 0x0400
		remote = 0x0D04
		break
	}
	return core.NewConnectRequest(local, remote)
}

func (c *client) validate(tpkt common.TPKT) (err error) {
	if c.pduLength > 0 && int(tpkt.GetLength()) > c.pduLength+isoHeaderSize || int(tpkt.GetLength()) < minPduSize {
		err = common.ErrorWithCode(common.ErrCliResponseInvalid)
	}
	return
}

func checkReqAck(req *core.PDU, ack *core.PDU) (err error) {
	if ack.GetHeader() == nil {
		return nil
	}

	if ackHeader, ok := ack.GetHeader().(*core.AckHeader); ok && ackHeader.ErrorClass != 0x00 {
		err = common.ErrorWithCode(common.ErrCliResponseExceptional,
			common.ErrorClassDescOrDefault(ackHeader.ErrorClass, "UnKnown"),
			common.ErrorCodeDescOrDefault(ackHeader.ErrorCode, "UnKnown"))
		return
	}

	if ack.GetHeader().GetPduReference() != req.GetHeader().GetPduReference() {
		err = common.ErrorWithCode(common.ErrCliPduReferenceMismatch)
		return
	}

	if parameter, ok := ack.GetParameter().(*core.UserdataAckParameter); ok && parameter.ErrorClass != 0x00 {
		err = common.ErrorWithCode(common.ErrCliResponseExceptional,
			common.ErrorClassDescOrDefault(parameter.ErrorClass, "UnKnown"),
			common.ErrorCodeDescOrDefault(parameter.ErrorCode, "UnKnown"))
		return
	}

	if ack.GetDatum() == nil {
		return
	}

	if datum, ok := ack.GetDatum().(*core.ReadWriteDatum); !ok {
		return
	} else if readWriteParameter, ok := req.GetParameter().(*core.ReadWriteParameter); ok {
		if len(datum.ReturnItems) != int(readWriteParameter.ItemCount) {
			err = common.ErrorWithCode(common.ErrCliResponseLengthMismatch)
			return
		}
		for i := 0; i < len(datum.ReturnItems); i++ {
			item := datum.ReturnItems[i]
			if item.GetReturnCode() != common.RcSuccess {
				err = common.ErrorWithCode(common.ErrCliResponseExceptional,
					"UnKnown", common.ReturnCodeDescOrDefault(item.GetReturnCode(), "UnKnown"))
				return
			}
		}
	}
	return
}

func (c *client) read(requests []common.RequestItem) *ReadToken {
	token := NewToken(TtRead).(*ReadToken)

	if len(requests) == 0 {
		token.setError(common.ErrorWithCode(common.ErrCliRequestDataEmpty))
		return token
	}

	go func() {
		rawNumbers := make([]uint16, 0, len(requests))
		result := make([]*core.DataItem, 0, len(requests))
		for _, request := range requests {
			request := request.(*core.StandardRequestItem)
			rawNumbers = append(rawNumbers, request.Count)
			result = append(result, core.NewReqDataItem(make([]byte, int(request.VariableType.Size()*request.Count)), request.VariableType.DataVariableType()))
		}

		groups := util.ReadRecombination(rawNumbers, c.pduLength-14, 5, 12)
		for _, group := range groups {
			newRequestItems := make([]common.RequestItem, 0)
			for i := 0; i < len(group.Items); i++ {
				item := group.Items[i]
				requestItem := *(requests[item.Index].(*core.StandardRequestItem))
				requestItem.Count = uint16(item.RipeSize)
				requestItem.ByteAddress += item.SplitOffset
				newRequestItems = append(newRequestItems, &requestItem)
			}
			request := core.NewReadRequest(newRequestItems, c.GeneratePduNumber())
			pduToken := c.send(request)
			ack, err := pduToken.Wait()
			if err != nil {
				token.setError(err)
				return
			}
			datum := ack.GetDatum().(*core.ReadWriteDatum)
			for i := 0; i < len(group.Items); i++ {
				item := group.Items[i]
				copy(result[item.Index].Data[item.SplitOffset:], datum.ReturnItems[i].(*core.DataItem).Data)
			}
		}
		token.v = result
		token.flowComplete()
	}()

	return token
}

func (c *client) write(requests []common.RequestItem, dataItems []common.ResponseItem) *SimpleToken {
	token := NewToken(TtSimple).(*SimpleToken)
	if len(requests) == 0 || len(dataItems) == 0 {
		token.setError(common.ErrorWithCode(common.ErrCliRequestDataEmpty))
		return token
	}
	if len(requests) != len(dataItems) {
		token.setError(common.ErrorWithCode(common.ErrCliRequestDataDifferent))
		return token
	}

	go func() {
		rawNumbers := make([]uint16, 0, len(requests))
		for _, request := range requests {
			request := request.(*core.StandardRequestItem)
			rawNumbers = append(rawNumbers, request.Count)
		}

		groups := util.WriteRecombination(rawNumbers, c.pduLength-12, 17)
		for _, group := range groups {
			items := group.Items
			newRequestItems := make([]common.RequestItem, 0)
			newDataItems := make([]common.ResponseItem, 0)
			for i := 0; i < len(items); i++ {
				item := items[i]
				requestItem := *(requests[item.Index].(*core.StandardRequestItem))
				requestItem.Count = uint16(item.RipeSize)
				requestItem.ByteAddress += item.SplitOffset
				newRequestItems = append(newRequestItems, &requestItem)

				dataItem := *(dataItems[item.Index].(*core.DataItem))
				dataItem.Count = uint16(item.RipeSize)
				dataItem.Data = dataItem.Data[item.SplitOffset : item.SplitOffset+item.RipeSize]
				newDataItems = append(newDataItems, &dataItem)
			}

			request := core.NewWriteRequest(newRequestItems, newDataItems, c.GeneratePduNumber())
			pduToken := c.send(request)
			_, err := pduToken.Wait()
			if err != nil {
				token.setError(err)
				return
			}
		}
		token.flowComplete()
	}()
	return token
}

func (c *client) send(request *core.PDU) *PduToken {
	p := NewToken(TtPdu).(*PduToken)
	if c.GetConn() == nil {
		p.setError(common.ErrorWithCode(common.ErrCliConnectionNil, c.host, c.port))
		return p
	}

	status := c.status.ConnectionStatus()
	if !(status == connected ||
		((status == connecting || status == reconnecting) && request.GetCOTP().GetPduType() == common.PtConnectRequest)) {
		_, ok := request.GetParameter().(*core.SetupComParameter)
		if (status != connecting && status != reconnecting) && !ok {
			p.setError(common.ErrorWithCode(common.ErrCliConnectionInactive, c.host, c.port))
			return p
		}
	}

	go func() {
		var (
			ctx RequestContext
			err error
		)
		switch request.GetCOTP().GetPduType() {
		case common.PtDisconnectRequest:
			ctx = &ConnectRequestContext{
				Request:  request,
				Response: make(chan *core.PDU),
				Error:    make(chan error),
			}
			err = c.tcpClient.handleDisconnectRequestContext(ctx)
			if err != nil {
				p.setError(err)
				return
			}
		case common.PtConnectRequest:
			ctx = &ConnectRequestContext{
				Request:  request,
				Response: make(chan *core.PDU),
				Error:    make(chan error),
			}
			err = c.tcpClient.handleConnectRequestContext(ctx)
			if err != nil {
				p.setError(err)
				return
			}
		default:
			ctx = &StandardRequestContext{
				RequestId: request.GetHeader().GetPduReference(),
				Request:   request,
				Response:  make(chan *core.PDU),
				Error:     make(chan error),
			}
			err = c.tcpClient.handleRequestContext(ctx)
			if err != nil {
				p.setError(err)
				return
			}
		}
		pdu := ctx.GetRequest().ToBytes()
		_, err = c.conn.Write(pdu)
		if err != nil {
			p.setError(err)
			return
		}
		c.logger.Debugf("S7 client sending: % x", pdu)
		ack, err := ctx.GetResponse()
		if err != nil {
			p.setError(err)
			return
		}
		p.err = checkReqAck(request, ack)
		p.v = ack
		p.flowComplete()
	}()
	return p
}

func (c *client) Connect() *ConnectToken {
	t := NewToken(TtConnect).(*ConnectToken)
	fn, err := c.status.Connecting()
	if err != nil {
		t.setError(err)
		return t
	}

	go func() {
		backoff := c.retryInterval
		var conn net.Conn
		for retries := 1; c.maxRetries == -1 || retries <= c.maxRetries; retries++ {
			conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.host, c.port), c.timeout)
			if err == nil {
				break
			}

			c.logger.Debugf("retrying in %v, failed to connect to %s: %v", backoff, fmt.Sprintf("%s:%d", c.host, c.port), err)
			if backoff < c.maxRetryBackoff {
				backoff *= 2
			}
			if retries != c.maxRetries {
				time.Sleep(backoff)
			}
		}
		if err != nil {
			_ = fn(false)
			t.setError(err)
			c.disconnectedWithError(common.ErrorWithCode(common.ErrTcpConnect, err))
			return
		}
		gc, _ := c.cli.Enroll(conn)
		c.SetConn(gc)

		c.logger.Infof("S7 client start iso connect for [%s]", fmt.Sprintf("%s:%d", c.host, c.port))
		isoToken := c.send(c.isoConnect())
		_, err = isoToken.Wait()
		if err != nil {
			_ = fn(false)
			t.setError(common.ErrorWithCode(common.ErrTcpConnect, err))
			c.disconnectedWithError(common.ErrorWithCode(common.ErrTcpConnect, err))
			return
		}
		var ack *core.PDU
		dtToken := c.send(core.NewConnectDt(uint16(c.pduLength), c.GeneratePduNumber()))
		ack, err = dtToken.Wait()
		if err != nil {
			_ = fn(false)
			t.setError(common.ErrorWithCode(common.ErrTcpConnect, err))
			c.disconnectedWithError(common.ErrorWithCode(common.ErrTcpConnect, err))
			return
		}
		if ack.GetCOTP().GetPduType() != common.PtData {
			_ = fn(false)
			t.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			c.disconnectedWithError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		if ack.GetHeader() == nil || ack.GetHeader().Len() != common.AckHeaderLen {
			_ = fn(false)
			t.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			c.disconnectedWithError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		parameter, ok := ack.GetParameter().(*core.SetupComParameter)
		if !ok {
			_ = fn(false)
			t.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			c.disconnectedWithError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		if parameter.PduLength <= 0 {
			_ = fn(false)
			t.setError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			c.disconnectedWithError(common.ErrorWithCode(common.ErrCliResponseInvalid))
			return
		}
		c.pduLength = int(parameter.PduLength)

		c.logger.Infof("S7 client for [%s] is active", fmt.Sprintf("%s:%d", c.host, c.port))
		_ = fn(true)
		t.v = c
		t.flowComplete()
		go util.Invoke(c.onConnected, []interface{}{(*Client)(nil)}, c)
	}()
	return t
}

func (c *client) tcpOnOpen(gnet.Conn) {
	return
}

func (c *client) tcpOnClose(_ gnet.Conn, err error) {
	c.disconnectedWithError(err)
	disFn, err := c.status.ConnectionLost(c.autoReconnect && c.status.ConnectionStatus() > connecting)
	if err != nil {
		return
	}

	go func() {
		c.SetConn(nil)
		if reConnFn, err := disFn(true); err == nil && reConnFn != nil {
			go c.reconnect(reConnFn)
		}
	}()
}

func (c *client) disconnectedWithError(err error) {
	c.SetConn(nil)
	c.logger.Warnf("S7 client for [%s] is disconnected with error: [%s]", fmt.Sprintf("%s:%d", c.host, c.port), err)
	go util.Invoke(c.onDisconnected, []interface{}{(*Client)(nil), (*error)(nil)}, c, err)
}

func (c *client) reconnect(connectionUp connCompletedFn) {
	c.logger.Debugf("client start reconnect")
	var (
		conn net.Conn
		err  error
	)

	backoff := c.reconnectInterval
	for retries := 1; c.maxReconnectTimes == -1 || retries <= c.maxReconnectTimes; retries++ {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.host, c.port), c.timeout)
		if err == nil {
			break
		}

		c.logger.Debugf("retrying in %v, failed to connect to %s: %v", backoff, fmt.Sprintf("%s:%d", c.host, c.port), err)
		if backoff < c.maxReconnectBackoff {
			backoff *= 2
		}
		if retries != c.maxReconnectTimes {
			time.Sleep(backoff)
		}
	}
	if err != nil {
		_ = connectionUp(false)
		return
	}
	gc, _ := c.cli.Enroll(conn)
	c.SetConn(gc)

	c.logger.Infof("S7 client start iso connect for [%s]", fmt.Sprintf("%s:%d", c.host, c.port))
	isoToken := c.send(c.isoConnect())
	_, err = isoToken.Wait()
	if err != nil {
		_ = connectionUp(false)
		return
	}
	var ack *core.PDU
	dtToken := c.send(core.NewConnectDt(uint16(c.pduLength), c.GeneratePduNumber()))
	ack, err = dtToken.Wait()
	if err != nil {
		_ = connectionUp(false)
		return
	}
	if ack.GetCOTP().GetPduType() != common.PtData {
		_ = connectionUp(false)
		return
	}
	if ack.GetHeader() == nil || ack.GetHeader().Len() != common.AckHeaderLen {
		_ = connectionUp(false)
		return
	}
	parameter, ok := ack.GetParameter().(*core.SetupComParameter)
	if !ok {
		_ = connectionUp(false)
		return
	}
	if parameter.PduLength <= 0 {
		_ = connectionUp(false)
		return
	}
	c.pduLength = int(parameter.PduLength)

	c.logger.Infof("S7 client for [%s] is active", fmt.Sprintf("%s:%d", c.host, c.port))
	_ = connectionUp(true)
	go util.Invoke(c.onConnected, []interface{}{(*Client)(nil)}, c)
}

func (c *client) Disconnect() {
	fn, err := c.status.Disconnecting()
	if err != nil {
		c.logger.Warnf("client disconnecting failed with error: %s", err.Error())
		return
	}
	c.disconnect()
	fn()
}

func (c *client) disconnect() {
	if c.cli != nil {
		_ = c.cli.Stop()
		c.cli = nil
	}
	if conn := c.GetConn(); conn != nil {
		_ = conn.Close()
		c.SetConn(nil)
	}
}

func (c *client) GetStatus() Status {
	return c.status.ConnectionStatus()
}

func (c *client) IsConnected() bool {
	s, r := c.status.ConnectionStatusRetry()
	switch {
	case s == connected:
		return true
	case c.connectRetry && s == connecting:
		return true
	case c.autoReconnect:
		return s == reconnecting || (s == disconnecting && r)
	default:
		return false
	}
}

func (c *client) IsConnectionOpen() bool {
	return c.status.ConnectionStatus() == connected
}

func siemensTimestamp(encodedDate int64) time.Time {
	return time.Date(1984, 1, 1, 0, 0, 0, 0, time.UTC).
		Add(time.Second * time.Duration(encodedDate*86400))
}

func (c *client) parseReadRequestItems(addresses []string) (items []common.RequestItem, ots []common.ParamVariableType, err error) {
	if len(addresses) == 0 {
		err = common.ErrorWithCode(common.ErrAddressEmpty)
		return
	}
	items = make([]common.RequestItem, 0)
	ots = make([]common.ParamVariableType, 0, len(addresses))
	var item common.RequestItem
	for i := 0; i < len(addresses); i++ {
		address := addresses[i]
		item, err = core.ParseAddress(address)
		if err != nil {
			return
		}
		ot := item.(*core.StandardRequestItem).VariableType
		err = c.parseRequestItem(item.(*core.StandardRequestItem))
		if err != nil {
			return
		}
		items = append(items, item)
		ots = append(ots, ot)
	}
	return
}

func (c *client) parsesWriteRequestItems(addresses []string, data [][]byte) (requests []common.RequestItem, dataItems []common.ResponseItem, err error) {
	if len(addresses) == 0 {
		err = common.ErrorWithCode(common.ErrAddressEmpty)
		return
	}
	if len(addresses) != len(data) {
		err = common.ErrorWithCode(common.ErrCliRequestDataDifferent)
		return
	}

	var item common.RequestItem
	requests = make([]common.RequestItem, 0)
	dataItems = make([]common.ResponseItem, 0)
	for i := 0; i < len(addresses); i++ {
		address := addresses[i]
		item, err = core.ParseAddress(address)
		if err != nil {
			return
		}
		if item.(*core.StandardRequestItem).VariableType == common.PvtString || item.(*core.StandardRequestItem).VariableType == common.PvtWString {
			item.(*core.StandardRequestItem).Count = item.(*core.StandardRequestItem).Count * uint16(len(data[i]))
			item.(*core.StandardRequestItem).VariableType = common.PvtByte
		} else if item.(*core.StandardRequestItem).VariableType != common.PvtBit {
			item.(*core.StandardRequestItem).Count = item.(*core.StandardRequestItem).Count * item.(*core.StandardRequestItem).VariableType.Size()
			item.(*core.StandardRequestItem).VariableType = common.PvtByte
		}
		dataItem := core.NewReqDataItem(data[i], item.(*core.StandardRequestItem).VariableType.DataVariableType())
		requests = append(requests, item)
		dataItems = append(dataItems, dataItem)
	}
	return
}

func (c *client) parseRequestItem(item *core.StandardRequestItem) (err error) {
	switch item.VariableType {
	case common.PvtString:
		item.VariableType = common.PvtByte
		count := 2
		if c.plcType == common.S200Smart {
			count = 1
		}
		item.Count = uint16(count)
		var lRes []*core.DataItem
		token := c.read([]common.RequestItem{item})
		lRes, err = token.Wait()
		if err != nil {
			return
		}
		if len(lRes) == 0 {
			err = common.ErrorWithCode(common.ErrCliResponseInvalid)
			return
		}
		item.Count = uint16(count) + uint16(lRes[0].Data[count-1])
		break
	case common.PvtWString:
		item.VariableType = common.PvtByte
		count := 4
		if c.plcType == common.S200Smart {
			count = 2
		}
		item.Count = uint16(count)
		var lRes []*core.DataItem
		token := c.read([]common.RequestItem{item})
		lRes, err = token.Wait()
		if err != nil {
			return
		}
		if len(lRes) == 0 {
			err = common.ErrorWithCode(common.ErrCliResponseInvalid)
			return
		}
		item.Count = uint16(count) + binary.BigEndian.Uint16(lRes[0].Data[count-2:])*2
		break
	case common.PvtTime, common.PvtDate, common.PvtTimeOfDay, common.PvtDateTime, common.PvtS5Time, common.PvtDTL:
		item.Count = item.VariableType.Size() * item.Count
		item.VariableType = common.PvtByte
		break
	}
	return
}

func (c *client) GeneratePduNumber() uint16 {
	index := c.IncrementAndGetPduIndex()
	if index >= 65536 {
		atomic.StoreUint32(&c.pduIndex, 0)
		index = c.IncrementAndGetPduIndex()
	}
	return uint16(index)
}

func (c *client) IncrementAndGetPduIndex() uint32 {
	return atomic.AddUint32(&c.pduIndex, 1)
}

func (c *client) SetConn(conn gnet.Conn) {
	c.m.RLock()
	defer c.m.RUnlock()
	c.conn = conn
}

func (c *client) GetConn() net.Conn {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.conn
}

func (c *client) GetPduLength() int {
	return c.pduLength
}
