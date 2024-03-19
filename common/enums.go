// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package common

// AreaType 数据区域
type AreaType byte

// ParamVariableType Transport size (variable Type) in Item dat
type ParamVariableType byte

// PduType PDU类型
type PduType byte

// MessageType
// 消息的一般类型（有时称为ROSCTR类型）
// 消息的其余部分在很大程度上取决于Message Type和功能代码。
type MessageType byte

// FunctionCode 功能码
// Job request/Ack-Data function codes
type FunctionCode byte

// SyntaxID 寻址模式和项结构其余部分的格式
// 它具有任意类型寻址的常量值0x10
type SyntaxID byte

// NckModule NCK的模块
type NckModule byte

// ReturnCode 操作的返回值，0xff信号成功
// 在写入请求消息中，此字段始终设置为零
type ReturnCode byte

// DataVariableType 数据返回变量的类型和长度
// Transport size in data Transport size (variable Type)
type DataVariableType byte

// BlockType 文件地址块类型
type BlockType uint16

// DestinationFileSystem 目标文件系统
type DestinationFileSystem byte

// PlcType plc类型
type PlcType byte

// Method 用户数据参数方法
type Method byte

// FunctionGroup 用户数据参数类型
type FunctionGroup byte

// CpuSubFunction  Cpu子方法
type CpuSubFunction byte

// BlockSubFunction  块子方法
type BlockSubFunction byte

// TimeSubFunction  时间子方法
type TimeSubFunction byte

// SecuritySubFunction  安全子方法
type SecuritySubFunction byte

// ParameterProtectionLevel 参数保护级别
type ParameterProtectionLevel uint16

// CpuProtectionLevel CPU保护级别
type CpuProtectionLevel uint16

// SelectorSetting 选择器设置
type SelectorSetting uint16

// StartupSwitch 启动开环
type StartupSwitch uint16

const (
	// S200 S200
	S200 PlcType = 0x00
	// S200Smart S200_SMART
	S200Smart = 0x01
	// S300 S300
	S300 = 0x02
	// S400 S400
	S400 = 0x03
	// S1200 S1200
	S1200 = 0x04
	// S1500 S1500
	S1500 = 0x05
	// Sinumerik828d Sinumerik828d
	Sinumerik828d = 0x06

	// PvtBit 位
	PvtBit = 0x01
	// PvtByte 字节
	PvtByte = 0x02
	// PvtChar 字符
	PvtChar = 0x03
	// PvtWord 字
	PvtWord = 0x04
	// PvtInt INT
	PvtInt = 0x05
	// PvtDWord 双字
	PvtDWord = 0x06
	// PvtDInt DINT
	PvtDInt = 0x07
	// PvtReal 浮点
	PvtReal = 0x08
	// PvtDate 日期
	PvtDate = 0x09
	// PvtTimeOfDay TOD
	PvtTimeOfDay = 0x0A
	// PvtTime 时间
	PvtTime = 0x0B
	// PvtS5Time S5TIME
	PvtS5Time = 0x0C
	// PvtDateTime 日期和时间
	PvtDateTime = 0x0F
	// PvtDTL dtl
	PvtDTL = 0x10
	// PvtCounter 计数器
	PvtCounter = 0x1C
	// PvtTimer 定时器
	PvtTimer = 0x1D
	// PvtString 字符串
	PvtString ParamVariableType = 0x00
	// PvtWString 字符串
	PvtWString = 0xFF

	// AtSystemInfo 200系列系统信息
	AtSystemInfo AreaType = 0x03
	// AtSystemFlag 200系统标志
	AtSystemFlag = 0x05
	// AtAnalogInputs 200系列模拟量输入
	AtAnalogInputs = 0x06
	// AtAnalogOutputs 200系列模拟量输出
	AtAnalogOutputs = 0x07
	// AtDirectPeripheralAccess 直接访问外设
	AtDirectPeripheralAccess = 0x80
	// AtInputs 输入（I）
	AtInputs = 0x81
	// AtOutputs 输出（Q）
	AtOutputs = 0x82
	// AtFlags 内部标志（M）
	AtFlags = 0x83
	// AtDataBlocks 数据块（DB）
	AtDataBlocks = 0x84
	// AtInstanceDataBlocks 背景数据块（DI）
	AtInstanceDataBlocks = 0x85
	// AtLocalData 局部变量（L)
	AtLocalData = 0x86
	// AtUnknownYet 全局变量（V）
	AtUnknownYet = 0x87
	// AtCounters S7计数器（C）
	AtCounters = 0x1C
	// AtTimers S7定时器（T）
	AtTimers = 0x1D
	// AtIecCounters Iec计数器（200系列）
	AtIecCounters = 0x1E
	// AtIecTimers Iec定时器（200系列）
	AtIecTimers = 0x1F

	// PtConnectRequest 连接请求
	PtConnectRequest PduType = 0xE0
	// PtConnectConfirm 连接请求
	PtConnectConfirm PduType = 0xD0
	// PtDisconnectRequest 断开请求
	PtDisconnectRequest PduType = 0x80
	// PtDisconnectConfirm 断开确认
	PtDisconnectConfirm PduType = 0xC0
	// PtReject 拒绝
	PtReject PduType = 0x50
	// PtData 数据
	PtData PduType = 0xF0

	// MtJob 开工干活的意思，主设备通过job向从设备发出“干活”的命令
	// 具体是读取数据还是写数据由parameter决定
	MtJob MessageType = 0x01
	// MtAck 确认 确认有没有数据字段
	MtAck = 0x02
	// MtAckData 从设备回应主设备的job
	MtAckData = 0x03
	// MtUserData 原始协议的扩展，参数字段包含请求/响应id
	// 用于编程/调试，SZL读取，安全功能，时间设置，循环读取
	MtUserData = 0x07

	// FcCpuService CPU服务
	FcCpuService FunctionCode = 0x00
	// FcRead 读变量
	FcRead = 0x04
	// FcWrite 写变量
	FcWrite = 0x05
	// FcStartDownload 开始下载
	FcStartDownload = 0xFA
	// FcDownload 下载阻塞
	FcDownload = 0xFB
	// FcEndDownload 下载结束
	FcEndDownload = 0xFC
	// FcStartUpload 开始上传
	FcStartUpload = 0x1D
	// FcUpload 上传
	FcUpload = 0x1E
	// FcEndUpload 结束上传
	FcEndUpload = 0x1F
	// FcControl 控制PLC
	FcControl = 0x28
	// FcStop 停止PLC
	FcStop = 0x29
	// FcSetupCom 设置通信
	FcSetupCom = 0xF0

	// SiAny Address data S7-Any pointer-like DB1.DBX10.2
	SiAny SyntaxID = 0x10
	// SiPbcRId R_ID for PBC
	SiPbcRId = 0x13
	// SiAlarmLockFree Alarm lock/free dataset
	SiAlarmLockFree = 0x15
	// SiAlarmInd Alarm indication dataset
	SiAlarmInd = 0x16
	// SiAlarmAck Alarm acknowledge message dataset
	SiAlarmAck = 0x19
	// SiAlarmQueryReq Alarm query request dataset
	SiAlarmQueryReq = 0x1A
	// SiNotifyInd Notify indication dataset
	SiNotifyInd = 0x1C
	// SiDriveesAny DRIVEESANY seen on Drive ES Starter with routing over S7
	SiDriveesAny = 0xA2
	// SiS1200SYM Symbolic byteAddress mode of S7-1200
	SiS1200SYM = 0xB2
	// SiDbRead Kind of DB block read， seen only at an S7-400
	SiDbRead = 0xB0
	// SiNck Sinumerik NCK HMI access
	SiNck = 0x82

	// NmY Global system data
	NmY NckModule = 0x10
	// NmYNCFL NCK instruction groups
	NmYNCFL = 0x11
	// NmFU NCU global settable frames
	NmFU = 0x12
	// NmFA Active NCU global frames
	NmFA = 0x13
	// NmTO Tool data
	NmTO = 0x14
	// NmRP Arithmetic parameters
	NmRP = 0x15
	// NmSE Setting data
	NmSE = 0x16
	// NmSGUD SGUD( (byte) ,Block
	NmSGUD = 0x17
	// NmLUD Local userdata
	NmLUD = 0x18
	// NmTC Toolholder parameters
	NmTC = 0x19
	// NmM Machine data
	NmM = 0x1A
	// NmWAL Working area limitation
	NmWAL = 0x1C
	// NmDIAG Internal diagnostic data
	NmDIAG = 0x1E
	// NmCC Unknown
	NmCC = 0x1F
	// NmFE Channel( (byte) ,specific external frame
	NmFE = 0x20
	// NmTD Tool data: General data
	NmTD = 0x21
	// NmTS Tool edge data: Monitoring data
	NmTS = 0x22
	// NmTG Tool data: Grinding( (byte) ,specific data
	NmTG = 0x23
	// NmTU Tool data
	NmTU = 0x24
	// NmTUE Tool edge data, userdefined data
	NmTUE = 0x25
	// NmTV Tool data, directory
	NmTV = 0x26
	// NmTM Magazine data: General data
	NmTM = 0x27
	// NmTP Magazine data: Location data
	NmTP = 0x28
	// NmTPM Magazine data: Multiple assignment of location data
	NmTPM = 0x29
	// NmTT Magazine data: Location typ
	NmTT = 0x2A
	// NmTMV Magazine data: Directory
	NmTMV = 0x2B
	// NmTMC Magazine data: Configuration data
	NmTMC = 0x2C
	// NmMGUD MGUD( (byte) ,Block
	NmMGUD = 0x2D
	// NmUGUD UGUD( (byte) ,Block
	NmUGUD = 0x2E
	// NmGUD4 GUD4( (byte) ,Block
	NmGUD4 = 0x2F
	// NmGUD5 GUD5( (byte) ,Block
	NmGUD5 = 0x30
	// NmGUD6 GUD6( (byte) ,Block
	NmGUD6 = 0x31
	// NmGUD7 GUD7( (byte) ,Block
	NmGUD7 = 0x32
	// NmGUD8 GUD8( (byte) ,Block
	NmGUD8 = 0x33
	// NmGUD9 GUD9( (byte) ,Block
	NmGUD9 = 0x34
	// NmPA Channel( (byte) ,specific protection zones
	NmPA = 0x35
	// NmGD1 SGUD( (byte) ,Block GD1
	NmGD1 = 0x36
	// NmNIB State data: Nibbling
	NmNIB = 0x37
	// NmETP Types of events
	NmETP = 0x38
	// NmETPD Data lists for protocolling
	NmETPD = 0x39
	// NmSYNACT Channel( (byte) ,specific synchronous actions
	NmSYNACT = 0x3A
	// NmDIAGN Diagnostic data
	NmDIAGN = 0x3B
	// NmVSYN Channel( (byte) ,specific user variables for synchronous actions
	NmVSYN = 0x3C
	// NmTUS Tool data: user monitoring data
	NmTUS = 0x3D
	// NmTUM Tool data: user magazine data
	NmTUM = 0x3E
	// NmTUP Tool data: user magazine place data
	NmTUP = 0x3F
	// NmTF Parameterizing, return parameters of _N_TMGETT, _N_TSEARC
	NmTF = 0x40
	// NmFB Channel( (byte) ,specific base frames
	NmFB = 0x41
	// NmSSP2 State data: Spindle
	NmSSP2 = 0x42
	// NmPUD program global Benutzerdaten
	NmPUD = 0x43
	// NmTOS Edge( (byte) ,related location( (byte) ,dependent fine total offsets
	NmTOS = 0x44
	// NmTOST Edge( (byte) ,related location( (byte) ,dependent fine total offsets, transformed
	NmTOST = 0x45
	// NmTOE Edge( (byte) ,related coarse total offsets, setup offsets
	NmTOE = 0x46
	// NmTOET Edge( (byte) ,related coarse total offsets, transformed setup offsets
	NmTOET = 0x47
	// NmAD Adapter data
	NmAD = 0x48
	// NmTOT Edge data: Transformed offset data
	NmTOT = 0x49
	// NmAEV Working offsets: Directory
	NmAEV = 0x4A
	// NmYFAFL NCK instruction groups (Fanuc)
	NmYFAFL = 0x4B
	// NmFS System( (byte) ,Frame
	NmFS = 0x4C
	// NmSD Servo data
	NmSD = 0x4D
	// NmTAD Application( (byte) ,specific data
	NmTAD = 0x4E
	// NmTAO Application( (byte) ,specific cutting edge data
	NmTAO = 0x4F
	// NmTAS Application( (byte) ,specific monitoring data
	NmTAS = 0x50
	// NmTAM Application( (byte) ,specific magazine data
	NmTAM = 0x51
	// NmTAP Application( (byte) ,specific magazine location data
	NmTAP = 0x52
	// NmMEM Unknown
	NmMEM = 0x53
	// NmSALUC Alarm actions: List in reverse chronological order
	NmSALUC = 0x54
	// NmAUXFU Auxiliary functions
	NmAUXFU = 0x55
	// NmTDC Tool/Tools
	NmTDC = 0x56
	// NmCP Generic coupling
	NmCP = 0x57
	// NmSDME Unknown
	NmSDME = 0x6E
	// NmSPARPI Program pointer on interruption
	NmSPARPI = 0x6F
	// NmSEGA State data: Geometry axes in tool offset memory (extended)
	NmSEGA = 0x70
	// NmSEMA State data: Machine axes (extended)
	NmSEMA = 0x71
	// NmSSP State data: Spindle
	NmSSP = 0x72
	// NmSGA State data: Geometry axes in tool offset memory
	NmSGA = 0x73
	// NmSMA State data: Machine axes
	NmSMA = 0x74
	// NmSALAL Alarms: List organized according to time
	NmSALAL = 0x75
	// NmSALAP Alarms: List organized according to priority
	NmSALAP = 0x76
	// NmSALA Alarms: List organized according to time
	NmSALA = 0x77
	// NmSSYNAC Synchronous actions
	NmSSYNAC = 0x78
	// NmSPARPF Program pointers for block search and stop run
	NmSPARPF = 0x79
	// NmSPARPP Program pointer in automatic operation
	NmSPARPP = 0x7A
	// NmSNCF Active G functions
	NmSNCF = 0x7B
	// NmSPARP Part program information
	NmSPARP = 0x7D
	// NmSINF Part( (byte) ,program( (byte) ,specific status data
	NmSINF = 0x7E
	// NmS State data
	NmS = 0x7F
	// NmUNKNOWN1 State data
	NmUNKNOWN1 = 0x80
	// NmUNKNOWN2 State data
	NmUNKNOWN2 = 0x81
	// NmUNKNOWN3 State data
	NmUNKNOWN3 = 0x82
	// NmUNKNOWN4 State data
	NmUNKNOWN4 = 0x83
	// NmUNKNOWN5 State data
	NmUNKNOWN5 = 0x84
	// NmUNKNOWN6 State data
	NmUNKNOWN6 = 0x85

	// RcReserved 未定义，预留
	RcReserved ReturnCode = 0x00
	// RcHardwareError 硬件错误
	RcHardwareError = 0x01
	// RcAccessingTheObjectNotAllowed 对象不允许访问
	RcAccessingTheObjectNotAllowed = 0x03
	// RcInvalidAddress 无效地址，所需的地址超出此PLC的极限
	RcInvalidAddress = 0x05
	// RcDataTypeNotSupported 数据类型不支持
	RcDataTypeNotSupported = 0x06
	// RcDataTypeInconsistent 数据类型不一致
	RcDataTypeInconsistent = 0x07
	// RcObjectDoesNotExist 对象不存在
	RcObjectDoesNotExist = 0x0A
	// RcSuccess 成功
	RcSuccess = 0xFF

	// DvtNull 无
	DvtNull DataVariableType = 0x00
	// DvtBit bit access, len is in bits
	DvtBit = 0x03
	// DvtByteWordDword byte/word/dword access, len is in bits
	DvtByteWordDword = 0x04
	// DvtInt int access, len is in bits
	DvtInt = 0x05
	// DvtDint int access, len is in bytes
	DvtDint = 0x06
	// DvtReal real access, len is in bytes
	DvtReal = 0x07
	// DvtOctetString octet string, len is in bytes
	DvtOctetString = 0x09

	DtOb  BlockType = 0x3038
	DtDb            = 0x3041
	DtSdb           = 0x3042
	DtFc            = 0x3043
	DtSfc           = 0x3044
	DtFb            = 0x3045
	DtSfb           = 0x3046

	// DfsP （Passive (copied, but not chained) module)：被动文件系统
	DfsP DestinationFileSystem = 0x50
	// DfsA (Active embedded module)：主动文件系统
	DfsA DestinationFileSystem = 0x41
	// DfsB (Active as well as passive module)：既主既被文件系统种
	DfsB DestinationFileSystem = 0x42

	// MRequest request
	MRequest Method = 0x11
	// MResponse response
	MResponse = 0x12

	// FgRequestModeTransition request for model transition
	FgRequestModeTransition FunctionGroup = 0x40
	// FgResponseModeTransition response for model transition
	FgResponseModeTransition = 0x80
	// FgRequestProgrammerCmd request for programmer command
	FgRequestProgrammerCmd = 0x41
	// FgResponseProgrammerCmd response for programmer command
	FgResponseProgrammerCmd = 0x81
	// FgRequestCyclicData request for cyclic data read
	FgRequestCyclicData = 0x42
	// FgResponseCyclicData response for cyclic data read
	FgResponseCyclicData = 0x82
	// FgRequestBlockFunction request for block functions
	FgRequestBlockFunction = 0x43
	// FgResponseBlockFunction response for block functions
	FgResponseBlockFunction = 0x83
	// FgRequestCpuFunction request cpu functions
	FgRequestCpuFunction = 0x44
	// FgResponseCpuFunction response cpu functions
	FgResponseCpuFunction = 0x84
	// FgRequestSecurity request for block security
	FgRequestSecurity = 0x45
	// FgResponseSecurity response for security
	FgResponseSecurity = 0x85
	// FgRequestPBC request for PBC
	FgRequestPBC = 0x45
	// FgResponsePBC response for PBC
	FgResponsePBC = 0x85
	// FgRequestTimeFunction request for time functions
	FgRequestTimeFunction = 0x47
	// FgResponseTimeFunction response for time functions
	FgResponseTimeFunction = 0x87
	// FgRequestNC request for NC programming
	FgRequestNC = 0x47
	// FgResponseNc response for NC programming
	FgResponseNc = 0x87

	// CsfReadSzl read szl data
	CsfReadSzl                CpuSubFunction = 0x01
	CsfMessageService                        = 0x02
	CsfDiagnosticMessage                     = 0x03
	CsfDisplayAlarm                          = 0x05
	CsfDisplayNotify                         = 0x06
	CsfLockAlarm                             = 0x07
	CsfLockNotify                            = 0x08
	CsfDisplayScan                           = 0x09
	CsfConfirmAlarm                          = 0x0B
	CsfConfirmDisplayAlarm                   = 0x0C
	CsfLockDisplayAlarm                      = 0x0D
	CsfCancelLockDisplayAlarm                = 0x0E
	CsfDisplayAlarmSQ                        = 0x11
	CsfDisplayAlarmS                         = 0x12
	CsfQueryAlarm                            = 0x13

	// BsfListBlock list block
	BsfListBlock BlockSubFunction = 0x01
	// BsfListBlockOfType list block of type
	BsfListBlockOfType = 0x02
	// BsfBlockInfo GetBlockInfo
	BsfBlockInfo = 0x03

	// TsfReadClock read clock
	TsfReadClock TimeSubFunction = 0x01
	// TsfSetClock set clock
	TsfSetClock TimeSubFunction = 0x02

	// SsfSetPassword set session password
	SsfSetPassword SecuritySubFunction = 0x01
	// SsfClearPassword clear session password
	SsfClearPassword = 0x02

	PplNoPassword        ParameterProtectionLevel = 0x0000
	PplSelectorPassword                           = 0x0001
	PplWritePassword                              = 0x0002
	PplReadWritePassword                          = 0x0003

	CplUnknown           CpuProtectionLevel = 0x0000
	CplAccessGrant                          = 0x0001
	CplReadOnly                             = 0x0002
	CplReadWritePassword                    = 0x0003

	SpUnknown SelectorSetting = 0x0000
	SpRun                     = 0x0001
	SpRunP                    = 0x0002
	SpStop                    = 0x0003
	SpMRES                    = 0x0004

	SsUnknown StartupSwitch = 0x0000
	SsCRST                  = 0x0001
	SsWRST                  = 0x0002
)

func (w ParamVariableType) DataVariableType() DataVariableType {
	switch w {
	case PvtBit:
		return DvtBit
	case PvtCounter, PvtTimer:
		return DvtOctetString
	default:
		return DvtByteWordDword
	}
}

func (w ParamVariableType) Size() uint16 {
	switch w {
	case PvtBit:
		return 1
	case PvtByte:
		return 1
	case PvtChar:
		return 1
	case PvtWord:
		return 2
	case PvtInt:
		return 2
	case PvtDWord:
		return 4
	case PvtDInt:
		return 4
	case PvtReal:
		return 4
	case PvtDate:
		return 2
	case PvtTimeOfDay:
		return 4
	case PvtTime:
		return 4
	case PvtS5Time:
		return 2
	case PvtDateTime:
		return 8
	case PvtDTL:
		return 12
	case PvtCounter:
		return 2
	case PvtTimer:
		return 2
	default:
		return 0
	}
}
