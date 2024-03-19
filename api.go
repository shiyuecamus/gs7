// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import (
	"gs7/common"
	"time"
)

type Client interface {
	// Connect plc connect
	Connect() *ConnectToken
	// GetStatus get s7 client connection status
	GetStatus() Status
	// IsConnected return plc is connected
	IsConnected() bool
	// IsConnectionOpen return plc is connection open
	IsConnectionOpen() bool
	// Disconnect plc connection close and release
	Disconnect()
	// GetPduLength return plc max pdu length
	GetPduLength() int
	// GetPlcType return plc type
	GetPlcType() common.PlcType

	// ReadRaw read raw bytes from address
	ReadRaw(address string) *SingleRawReadToken
	// ReadBatchRaw read batch raw bytes from addresses
	ReadBatchRaw(addresses []string) *BatchRawReadToken
	// ReadParsed read auto parsed data from address
	ReadParsed(address string) *SingleParsedReadToken
	// ReadBatchParsed read batch auto parsed data from address
	ReadBatchParsed(addresses []string) *BatchParsedReadToken
	// WriteRaw write raw bytes to plc address
	WriteRaw(address string, data []byte) *SimpleToken
	// WriteRawBatch write batch raw bytes to plc addresses
	WriteRawBatch(addresses []string, data [][]byte) *SimpleToken
	// DBRead data block read
	// Support exceeds the maximum pdu length.
	// If the maximum pdu length is exceeded, it will be divided into multiple requests
	// And the aggregated results will be returned after the last request
	DBRead(dbNumber int, start int, size int) *DbReadToken
	// DBWrite data block write
	// Support exceeds the maximum pdu length.
	// If the maximum pdu length is exceeded, it will be divided into multiple requests
	DBWrite(dbNumber int, start int, data []byte) *SimpleToken
	// DBGet get all data of the data block
	DBGet(dbNumber int) *DbReadToken
	// DBFill fill the data block to the specified byte
	DBFill(dbNumber int, fillByte byte) *SimpleToken

	// HotRestart Puts the CPU in run mode performing and hot start.
	HotRestart() *SimpleToken
	// ColdRestart change CPU into run mode performing and cold start
	ColdRestart() *SimpleToken
	// StopPlc change CPU to stop mode
	StopPlc() *SimpleToken
	// CopyRamToRom copy Ram to Rom
	CopyRamToRom() *SimpleToken
	// Compress compress
	Compress() *SimpleToken
	// InsertFile insert file
	InsertFile(blockType common.BlockType, blockNumber int) *SimpleToken
	// UploadFile upload file content from PLC to PC
	UploadFile(blockType common.BlockType, blockNumber int) *UploadToken
	// DownloadFile download file content from PC to PLC
	DownloadFile(bytes []byte, blockType common.BlockType, blockNumber int, mC7CodeLength int) *SimpleToken
	// ClockRead read plc clock
	ClockRead() *ClockReadToken
	// ClockSet set plc clock
	ClockSet(t time.Time) *SimpleToken

	// ReadSzl
	// ref: https://support.industry.siemens.com/cs/mdm/109755202?c=22058881035&lc=cs-CZ
	ReadSzl(szlId uint16, szlIndex uint16) *PduToken
	// GetSzlIds get szl ids
	GetSzlIds() *SzlIdsToken
	// GetCatalog get plc catalog（order code and version）
	GetCatalog() *CatalogToken
	// GetPlcStatus get plc mode running status
	GetPlcStatus() *PlcStatusToken
	// GetUnitInfo get plc unit info
	GetUnitInfo() *UnitInfoToken
	// GetCommunicationInfo get plc communication info
	GetCommunicationInfo() *CommunicationInfoToken
	// GetProtectionInfo get plc protection level info
	GetProtectionInfo() *ProtectionInfoToken

	// BlockList list blocks info（block count）
	BlockList() *BlockListToken
	// BlockListType list blocks of type
	BlockListType(blockType common.BlockType) *BlockListTypeToken
	// BlockInfo get block info
	BlockInfo(bt common.BlockType, bn int) *BlockInfoToken

	// SetPassword set session pwd
	SetPassword(pwd string) *SimpleToken
	// ClearPassword clear session pwd
	ClearPassword() *SimpleToken
}
