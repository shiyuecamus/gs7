// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package common

type ObjectBytes interface {
	Len() int
	ToBytes() []byte
}

type S7BaseData interface {
	ObjectBytes
	GetTPKT() TPKT
	GetCOTP() COTP
	GetHeader() Header
	GetParameter() Parameter
	GetDatum() Datum
	SelfCheck()
}

type TPKT interface {
	ObjectBytes
	// GetVersion 版本号，常量0x03 <br>
	// 字节大小：1
	// 字节序数：0
	GetVersion() byte
	SetVersion(byte)
	// GetReserved 预留，默认值0x00
	// 字节大小：1
	// 字节序数：1
	GetReserved() byte
	SetReserved(byte)
	// GetLength 长度，包括后面负载payload+版本号+预留+长度
	// 字节大小：2
	// 字节序数：2-3
	GetLength() uint16
	SetLength(uint16)
}

type COTP interface {
	ObjectBytes
	// GetLength 长度（但并不包含length这个字段）
	// 字节大小：1
	// 字节序数：0
	GetLength() byte
	SetLength(byte)
	// GetPduType PDU类型
	// 字节大小：1
	// 字节序数：1
	GetPduType() PduType
	SetPduType(PduType)
}

type Header interface {
	ObjectBytes
	// GetProtocolId 协议id
	// 字节大小：1
	// 字节序数：0
	GetProtocolId() byte
	SetProtocolId(byte)
	// GetMessageType pdu（协议数据单元（Protocol Data Unit））的类型
	// 字节大小：1
	// 字节序数：1
	GetMessageType() MessageType
	SetMessageType(MessageType)
	// GetReserved 保留
	// 字节大小：2
	// 字节序数：2-3
	GetReserved() []byte
	SetReserved([]byte)
	// GetPduReference pdu的参考–由主站生成，每次新传输递增，大端
	// 字节大小：2
	// 字节序数：4-5
	GetPduReference() uint16
	SetPduReference(uint16)
	// GetParameterLength 参数的长度（大端）
	// 字节大小：2
	// 字节序数：6-7
	GetParameterLength() uint16
	SetParameterLength(uint16)
	// GetDataLength 数据的长度（大端）
	// 字节大小：2
	// 字节序数：8-9
	GetDataLength() uint16
	SetDataLength(uint16)
}

type RequestItem interface {
	ObjectBytes
	// GetSpecificationType 变量规范
	// 对于读/写消息，它总是具有值0x12
	// 字节大小：1
	// 字节序数：0
	GetSpecificationType() byte
	SetSpecificationType(byte)
	// GetLengthOfFollowing 其余部分的长度规范
	// 字节大小：1
	// 字节序数：1
	GetLengthOfFollowing() byte
	SetLengthOfFollowing(byte)
	// GetSyntaxId 寻址模式和项结构其余部分的格式，它具有任意类型寻址的常量值0x10
	// 字节大小：1
	// 字节序数：2
	GetSyntaxId() SyntaxID
	SetSyntaxId(SyntaxID)
}

type Parameter interface {
	ObjectBytes
}

type UserdataParameter interface {
	Parameter
}

type PlcControlParamBlock interface {
	ObjectBytes
}

type Datum interface {
	ObjectBytes
}

type ResponseItem interface {
	ObjectBytes
	// GetReturnCode 返回码
	// 字节大小：1
	// 字节序数：0
	GetReturnCode() ReturnCode
	SetReturnCode(ReturnCode)
}
