// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/spf13/cast"
)

type ErrorCode uint16

const (
	ErrOk                       = 0x0000
	ErrCommon                   = 0x0001
	ErrModelFromBytes           = 0x0002
	ErrTypeNotResolved          = 0x0003
	ErrVariableTypeUnrecognized = 0x0004
	ErrPasswordLengthInvalid    = 0x0005

	ErrCliConnectionNil          = 0x0101
	ErrCliConnectionConnecting   = 0x0102
	ErrCliRequestDataEmpty       = 0x0103
	ErrCliRequestDataDifferent   = 0x0104
	ErrCliUploadFailed           = 0x0105
	ErrCliResponseInvalid        = 0x0106
	ErrCliResponseExceptional    = 0x0107
	ErrCliPduReferenceMismatch   = 0x0108
	ErrCliResponseLengthMismatch = 0x0109
	ErrCliConnectionInactive     = 0x0110
	ErrCliRequestItemInvalid     = 0x0111
	ErrCliSzlPartsInvalid        = 0x0112
	ErrCliConnectionNotNil       = 0x0113

	ErrTcpRequestProcessing   = 0x1001
	ErrTcpRequestTimeout      = 0x1002
	ErrTcpRequestRejected     = 0x1003
	ErrTcpConnect             = 0x1004
	ErrTcpResponseEmpty       = 0x1005
	ErrTcpConnectWithAttempts = 0x1006

	ErrAddressEmpty   = 0x1101
	ErrAddressInvalid = 0x1102
)

func ErrorWithCode(code int, params ...any) (err error) {
	switch code {
	case ErrCommon:
		return fmt.Errorf(cast.ToString(params[0]), params[1:]...)
	case ErrModelFromBytes:
		return fmt.Errorf("[%s] format must gte [%d] bytes", params...)
	case ErrTypeNotResolved:
		return fmt.Errorf("cannot be resolved %s: %d", params...)
	case ErrVariableTypeUnrecognized:
		return fmt.Errorf("variable type is unrecognized for: [%d]", params...)
	case ErrPasswordLengthInvalid:
		return fmt.Errorf("password length must lte %d", params...)
	case ErrCliConnectionNil:
		return fmt.Errorf("connection for [%s:%d] is nil", params...)
	case ErrCliConnectionConnecting:
		return fmt.Errorf("connection for [%s:%d] is connecting", params...)
	case ErrCliRequestDataEmpty:
		return errors.New("request data is empty")
	case ErrCliRequestDataDifferent:
		return errors.New("request data length is different from addresses")
	case ErrCliUploadFailed:
		return errors.New("upload failed from response status")
	case ErrCliResponseInvalid:
		return errors.New("invalid response")
	case ErrCliResponseExceptional:
		return fmt.Errorf("response exceptional, class:[%s], reason: [%s]", params...)
	case ErrCliPduReferenceMismatch:
		return errors.New("pdu reference mismatch")
	case ErrCliResponseLengthMismatch:
		return errors.New("response data does not match the length of request data")
	case ErrCliConnectionInactive:
		return fmt.Errorf("connection for [%s:%d] is inactive", params...)
	case ErrCliRequestItemInvalid:
		return errors.New("request item invalid")
	case ErrCliSzlPartsInvalid:
		return errors.New("szl parts invalid")
	case ErrCliConnectionNotNil:
		return fmt.Errorf("connection for [%s:%d] is not nil", params...)
	case ErrTcpRequestProcessing:
		return fmt.Errorf("tcp client request for [%d] is already processing", params...)
	case ErrTcpRequestTimeout:
		return errors.New("request timeout")
	case ErrTcpConnect:
		return fmt.Errorf("tcp connection with error: %s", params...)
	case ErrTcpResponseEmpty:
		return errors.New("empty response")
	case ErrTcpConnectWithAttempts:
		return fmt.Errorf("failed to connect to [%s] after [%d] attempts", params...)
	case ErrAddressEmpty:
		return errors.New("request address is empty")
	case ErrAddressInvalid:
		return errors.New("request address is invalid")
	default:
		return
	}
}

var ErrorClassDescMap = map[byte]string{
	0x00: "没有错误",
	0x81: "应用关系",
	0x82: "对象定义",
	0x83: "没有可用资源",
	0x84: "服务处理中错误",
	0x85: "请求错误",
	0x87: "访问错误",
	0xD2: "下载错误",
}

func ErrorClassDescOrDefault(errorClass byte, defaultVal string) string {
	s, ok := ErrorClassDescMap[errorClass]
	if ok {
		return s
	}
	return defaultVal
}

var ErrorCodeDescMap = map[uint16]string{
	0x0000: "没有错误",
	0x0110: "块号无效",
	0x0111: "请求长度无效",
	0x0112: "参数无效",
	0x0113: "块类型无效",
	0x0114: "找不到块",
	0x0115: "块已存在",
	0x0116: "块被写保护",
	0x0117: "块/操作系统更新太大",
	0x0118: "块号无效",
	0x0119: "输入的密码不正确",
	0x011A: "PG资源错误",
	0x011B: "PLC资源错误",
	0x011C: "协议错误",
	0x011D: "块太多（与模块相关的限制）",
	0x011E: "不再与数据库建立连接，或者S7DOS句柄无效",
	0x011F: "结果缓冲区太小",
	0x0120: "块结束列表",
	0x0140: "可用内存不足",
	0x0141: "由于缺少资源，无法处理作业",
	0x8001: "当块处于当前状态时，无法执行请求的服务",
	0x8003: "S7协议错误：传输块时发生错误",
	0x8100: "应用程序，一般错误：远程模块未知的服务",
	0x8104: "未在模块上实现此服务或报告了帧错误",
	0x8204: "对象的类型规范不一致",
	0x8205: "复制的块已存在且未链接",
	0x8301: "模块上的内存空间或工作内存不足，或者指定的存储介质不可访问",
	0x8302: "可用资源太少或处理器资源不可用",
	0x8304: "无法进一步并行上传。存在资源瓶颈",
	0x8305: "功能不可用",
	0x8306: "工作内存不足（用于复制，链接，加载AWP）",
	0x8307: "保持性工作记忆不够（用于复制，链接，加载AWP）",
	0x8401: "S7协议错误：无效的服务序列（例如，加载或上载块）",
	0x8402: "由于寻址对象的状态，服务无法执行",
	0x8404: "S7协议：无法执行该功能",
	0x8405: "远程块处于DISABLE状态（CFB）。该功能无法执行",
	0x8500: "S7协议错误：帧错误",
	0x8503: "来自模块的警报：服务过早取消",
	0x8701: "寻址通信伙伴上的对象时出错（例如，区域长度错误）",
	0x8702: "模块不支持所请求的服务",
	0x8703: "拒绝访问对象",
	0x8704: "访问错误：对象已损坏",
	0xD001: "协议错误：非法的作业号",
	0xD002: "参数错误：非法的作业变体",
	0xD003: "参数错误：模块不支持调试功能",
	0xD004: "参数错误：作业状态非法",
	0xD005: "参数错误：作业终止非法",
	0xD006: "参数错误：非法链路断开ID",
	0xD007: "参数错误：缓冲区元素数量非法",
	0xD008: "参数错误：扫描速率非法",
	0xD009: "参数错误：执行次数非法",
	0xD00A: "参数错误：非法触发事件",
	0xD00B: "参数错误：非法触发条件",
	0xD011: "调用环境路径中的参数错误：块不存在",
	0xD012: "参数错误：块中的地址错误",
	0xD014: "参数错误：正在删除/覆盖块",
	0xD015: "参数错误：标签地址非法",
	0xD016: "参数错误：由于用户程序错误，无法测试作业",
	0xD017: "参数错误：非法触发号",
	0xD025: "参数错误：路径无效",
	0xD026: "参数错误：非法访问类型",
	0xD027: "参数错误：不允许此数据块数",
	0xD031: "内部协议错误",
	0xD032: "参数错误：结果缓冲区长度错误",
	0xD033: "协议错误：作业长度错误",
	0xD03F: "编码错误：参数部分出错（例如，保留字节不等于0）",
	0xD041: "数据错误：非法状态列表ID",
	0xD042: "数据错误：标签地址非法",
	0xD043: "数据错误：找不到引用的作业，检查作业数据",
	0xD044: "数据错误：标签值非法，检查作业数据",
	0xD045: "数据错误：HOLD中不允许退出ODIS控制",
	0xD046: "数据错误：运行时测量期间非法测量阶段",
	0xD047: "数据错误：“读取作业列表”中的非法层次结构",
	0xD048: "数据错误：“删除作业”中的非法删除ID",
	0xD049: "“替换作业”中的替换ID无效",
	0xD04A: "执行'程序状态'时出错",
	0xD05F: "编码错误：数据部分出错（例如，保留字节不等于0，...）",
	0xD061: "资源错误：没有作业的内存空间",
	0xD062: "资源错误：作业列表已满",
	0xD063: "资源错误：触发事件占用",
	0xD064: "资源错误：没有足够的内存空间用于一个结果缓冲区元素",
	0xD065: "资源错误：没有足够的内存空间用于多个结果缓冲区元素",
	0xD066: "资源错误：可用于运行时测量的计时器被另一个作业占用",
	0xD067: "资源错误：“修改标记”作业过多（特别是多处理器操作）",
	0xD081: "当前模式下不允许使用的功能",
	0xD082: "模式错误：无法退出HOLD模式",
	0xD0A1: "当前保护级别不允许使用的功能",
	0xD0A2: "目前无法运行，因为正在运行的函数会修改内存",
	0xD0A3: "I / O上活动的“修改标记”作业太多（特别是多处理器操作）",
	0xD0A4: "'强制'已经建立",
	0xD0A5: "找不到引用的作业",
	0xD0A6: "无法禁用/启用作业",
	0xD0A7: "无法删除作业，例如因为当前正在读取作业",
	0xD0A8: "无法替换作业，例如因为当前正在读取或删除作业",
	0xD0A9: "无法读取作业，例如因为当前正在删除作业",
	0xD0AA: "处理操作超出时间限制",
	0xD0AB: "进程操作中的作业参数无效",
	0xD0AC: "进程操作中的作业数据无效",
	0xD0AD: "已设置操作模式",
	0xD0AE: "作业是通过不同的连接设置的，只能通过此连接进行处理",
	0xD0C1: "访问标签时至少检测到一个错误",
	0xD0C2: "切换到STOP / HOLD模式",
	0xD0C3: "访问标记时至少检测到一个错误。模式更改为STOP / HOLD",
	0xD0C4: "运行时测量期间超时",
	0xD0C5: "块堆栈的显示不一致，因为块被删除/重新加载",
	0xD0C6: "作业已被删除，因为它所引用的作业已被删除",
	0xD0C7: "由于退出了STOP模式，因此作业被自动删除",
	0xD0C8: "由于测试作业和正在运行的程序之间不一致，“块状态”中止",
	0xD0C9: "通过复位OB90退出状态区域",
	0xD0CA: "通过在退出前重置OB90并访问错误读取标签退出状态范围",
	0xD0CB: "外设输出的输出禁用再次激活",
	0xD0CC: "调试功能的数据量受时间限制",
	0xD201: "块名称中的语法错误",
	0xD202: "函数参数中的语法错误",
	0xD205: "RAM中已存在链接块：无法进行条件复制",
	0xD206: "EPROM中已存在链接块：无法进行条件复制",
	0xD208: "超出模块的最大复制（未链接）块数",
	0xD209: "（至少）模块上找不到给定块之一",
	0xD20A: "超出了可以与一个作业链接的最大块数",
	0xD20B: "超出了一个作业可以删除的最大块数",
	0xD20C: "OB无法复制，因为关联的优先级不存在",
	0xD20D: "SDB无法解释（例如，未知数）",
	0xD20E: "没有（进一步）阻止可用",
	0xD20F: "超出模块特定的最大块大小",
	0xD210: "块号无效",
	0xD212: "标头属性不正确（与运行时相关）",
	0xD213: "SDB太多。请注意对正在使用的模块的限制",
	0xD216: "无效的用户程序 - 重置模块",
	0xD217: "不允许在模块属性中指定的保护级别",
	0xD218: "属性不正确（主动/被动）",
	0xD219: "块长度不正确（例如，第一部分或整个块的长度不正确）",
	0xD21A: "本地数据长度不正确或写保护错误",
	0xD21B: "模块无法压缩或压缩早期中断",
	0xD21D: "传输的动态项目数据量是非法的",
	0xD21E: "无法为模块（例如FM，CP）分配参数。系统数据无法链接",
	0xD220: "编程语言无效。请注意对正在使用的模块的限制",
	0xD221: "连接或路由的系统数据无效",
	0xD222: "全局数据定义的系统数据包含无效参数",
	0xD223: "通信功能块的实例数据块错误或超出最大背景数据块数",
	0xD224: "SCAN系统数据块包含无效参数",
	0xD225: "DP系统数据块包含无效参数",
	0xD226: "块中发生结构错误",
	0xD230: "块中发生结构错误",
	0xD231: "至少有一个已加载的OB无法复制，因为关联的优先级不存在",
	0xD232: "加载块的至少一个块编号是非法的",
	0xD234: "块在指定的内存介质或作业中存在两次",
	0xD235: "该块包含不正确的校验和",
	0xD236: "该块不包含校验和",
	0xD237: "您将要加载块两次，即CPU上已存在具有相同时间戳的块",
	0xD238: "指定的块中至少有一个不是DB",
	0xD239: "至少有一个指定的DB在装载存储器中不可用作链接变量",
	0xD23A: "至少有一个指定的DB与复制和链接的变体有很大不同",
	0xD240: "违反了协调规则",
	0xD241: "当前保护级别不允许该功能",
	0xD242: "处理F块时的保护冲突",
	0xD250: "更新和模块ID或版本不匹配",
	0xD251: "操作系统组件序列不正确",
	0xD252: "校验和错误",
	0xD253: "没有可用的可执行加载程序; 只能使用存储卡进行更新",
	0xD254: "操作系统中的存储错误",
	0xD280: "在S7-300 CPU中编译块时出错",
	0xD2A1: "块上的另一个块功能或触发器处于活动状态",
	0xD2A2: "块上的触发器处于活动状态。首先完成调试功能",
	0xD2A3: "块未激活（链接），块被占用或块当前被标记为删除",
	0xD2A4: "该块已被另一个块函数处理",
	0xD2A6: "无法同时保存和更改用户程序",
	0xD2A7: "块具有“未链接”属性或未处理",
	0xD2A8: "激活的调试功能阻止将参数分配给CPU",
	0xD2A9: "正在为CPU分配新参数",
	0xD2AA: "当前正在为模块分配新参数",
	0xD2AB: "当前正在更改动态配置限制",
	0xD2AC: "正在运行的激活或取消激活分配（SFC 12）暂时阻止R-KiR过程",
	0xD2B0: "在RUN（CiR）中配置时发生错误",
	0xD2C0: "已超出最大工艺对象数",
	0xD2C1: "模块上已存在相同的技术数据块",
	0xD2C2: "无法下载用户程序或下载硬件配置",
	0xD401: "信息功能不可用",
	0xD402: "信息功能不可用",
	0xD403: "服务已登录/注销（诊断/ PMC）",
	0xD404: "达到的最大节点数。不再需要登录诊断/ PMC",
	0xD405: "不支持服务或函数参数中的语法错误",
	0xD406: "当前不可用的必需信息",
	0xD407: "发生诊断错误",
	0xD408: "更新已中止",
	0xD409: "DP总线错误",
	0xD601: "函数参数中的语法错误",
	0xD602: "输入的密码不正确",
	0xD603: "连接已合法化",
	0xD604: "已启用连接",
	0xD605: "由于密码不存在，因此无法进行合法化",
	0xD801: "至少有一个标记地址无效",
	0xD802: "指定的作业不存在",
	0xD803: "非法的工作状态",
	0xD804: "非法循环时间（非法时基或多个）",
	0xD805: "不能再设置循环读取作业",
	0xD806: "引用的作业处于无法执行请求的功能的状态",
	0xD807: "功能因过载而中止，这意味着执行读取周期所需的时间比设置的扫描周期时间长",
	0xDC01: "日期和/或时间无效",
	0xE201: "CPU已经是主设备",
	0xE202: "由于闪存模块中的用户程序不同，无法进行连接和更新",
	0xE203: "由于固件不同，无法连接和更新",
	0xE204: "由于内存配置不同，无法连接和更新",
	0xE205: "由于同步错误导致连接/更新中止",
	0xE206: "由于协调违规而拒绝连接/更新",
	0xEF01: "S7协议错误：ID2错误; 工作中只允许00H",
	0xEF02: "S7协议错误：ID2错误; 资源集不存在",
}

func ErrorCodeDescOrDefault(errorCode []byte, defaultVal string) string {
	s, ok := ErrorCodeDescMap[binary.BigEndian.Uint16(errorCode)]
	if ok {
		return s
	}
	return defaultVal
}

var ReturnCodeDescMap = map[ReturnCode]string{
	RcReserved:                     "未知",
	RcSuccess:                      "成功",
	RcHardwareError:                "硬件错误",
	RcAccessingTheObjectNotAllowed: "对象不允许访问",
	RcInvalidAddress:               "无效地址，所需的地址超出此PLC的极限",
	RcDataTypeNotSupported:         "数据类型不支持",
	RcDataTypeInconsistent:         "数据类型不一致",
	RcObjectDoesNotExist:           "对象不存在",
}

func ReturnCodeDescOrDefault(returnCode ReturnCode, defaultVal string) string {
	s, ok := ReturnCodeDescMap[returnCode]
	if ok {
		return s
	}
	return defaultVal
}
