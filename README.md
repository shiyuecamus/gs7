# gs7

![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
![CopyRight-shiyuecamus](https://img.shields.io/badge/CopyRight-shiyuecamus-yellow)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

Implementation of Siemens S7 protocol in golang

Implementing underlying communication based on [gnet](https://github.com/panjf2000/gnet).

# 🍯 Overview

`gs7` is an high-performance and lightweight Siemens s7 protocol communication framework

`gs7` all read and write operations support both synchronous and asynchronous

# 🍇 Features

* Read DB, I, Q, M, V, Timer, Counter
* Synchronous and asynchronous read/write
* Single and batch raw read/write
* Large amounts of data read/write(exceeding the maximum limit of PLC: PduLength, the request will be automatically
  divided into multiple requests by the algorithm)
* Address single read/write
* Address batch read/write of multiple addresses with discontinuous addresses or addresses not in the same area
* Convert the read raw bytes to the type in golang
* Connection retry and automatic reconnection after connection lose
* Read SZL(System Status List)

# 🍆 Supported communication

* TCP

# 🍓 Quick start

Installation is as easy as:

```
go get github.com/shiyuecamus/gs7
```

# 🍉 Usage

```go
func main() {
  const (
    host = "192.168.0.1"
    port = 102
    rack = 0
    slot = 1
  )
  logger := logging.GetDefaultLogger()

  c := gs7.NewClientBuilder().
  PlcType(common.S1500).
  Host(host).
  Port(port).
  Port(102).
  Rack(rack).
  Slot(slot).
  Build()
  
  _, err := c.Connect().Wait()
  if err != nil {
    logger.Errorf("Failed to connect PLC, host: %s, port: %d, error: %s", host, port, err)
    return
  }
  
  // read
  res, err := c.ReadParsed("DB1.X0.0").Wait()
  if err != nil {
    logger.Errorf("Failed to read bit, error: %s", err)
    return
  }
  logger.Infof("Read bit success with value: %s", res)
  
  // write
  err = c.WriteRaw("DB1.X0.0", gs7.Bit(true).ToBytes()).Wait()
  if err != nil {
    logger.Errorf("Failed to read bit, error: %s", err)
    return
  }
  
  // check
  res, err = c.ReadParsed("DB1.X0.0").Wait()
  if err != nil {
    logger.Errorf("Failed to read bit, error: %s", err)
    return
  }
  logger.Infof("Read bit success with value: %s", res)
}

```

[more examples](_examples)

# 🍏 Knowledge

## Data address, region, type, and length mapping table

| abbreviation               | Area | DB Number | Byte index | Bit index | PLC Data type | Go Data Type  | ByteLength | PLC      |
|:---------------------------|------|-----------|:----------:|:---------:|:--------------|:--------------|:-----------|:---------|
| DB1.X0.1/DB1.BIT0.1        | DB   | 1         |     0      |     1     | Bit           | bool          | 1/8        | S1200    |
| QX1.6/QBIT1.6              | Q    | 0         |     1      |     6     | Bit           | bool          | 1/8        | S1200    |
| QB1/QBYTE1                 | Q    | 0         |     1      |     0     | Byte          | uint8         | 1          | S1200    |
| IX2.5/IBIT2.5              | I    | 0         |     2      |     5     | Bit           | bool          | 1/8        | S1200    |
| IW2/IWORD2                 | I    | 0         |     2      |     0     | Word          | uint16        | 2          | S1200    |
| MX3.2/MBIT3.2              | M    | 0         |     3      |     2     | Bit           | bool          | 1/8        | S1200    |
| MI3/MINT3                  | M    | 0         |     3      |     0     | Int           | int16         | 2          | S1200    |
| VX1.3/VBIT1.3              | V    | 1         |     3      |     0     | Bit           | bool          | 1/8        | 200Smart |
| VI4/VINT4                  | V    | 1         |     4      |     0     | Int           | int16         | 2          | 200Smart |
| T0                         | T    | 0         |     0      |     0     | Timer         | time.Duration | 2          | S1200    |
| C2                         | C    | 0         |     2      |     0     | Counter       | uint16        | 2          | S1200    |
| DB4.B4/DB4.BYTE4           | DB   | 4         |     4      |     0     | Byte          | uint8         | 1          | S1200    |
| DB4.C0/DB4.CHAR0           | DB   | 4         |     0      |     0     | Char          | int8          | 1          | S1200    |
| DB3.S2/DB3.STRING2         | DB   | 3         |     2      |     0     | String        | String        | N          | S1200    |
| DB1.I8/DB1.INT8            | DB   | 1         |     8      |     0     | Int           | int16         | 2          | S1200    |
| DB2.W12/DB1.WORD12         | DB   | 2         |     12     |     0     | Word          | uint16        | 2          | S1200    |
| DB1.DI0/DB1.DINT0          | DB   | 1         |     0      |     0     | DInt          | int32         | 4          | S1200    |
| DB1.DW0/DB1.DWAORD0        | DB   | 1         |     0      |     0     | DWord         | uint32        | 4          | S1200    |
| DB3.R2/DB3.REAL2           | DB   | 3         |     2      |     0     | Real          | float32       | 4          | S1200    |
| DB1.T2/DB1.TIME2           | DB   | 1         |     2      |     0     | Time          | time.Duration | 4          | S1200    |
| DB1.ST2/DB1.STIME2         | DB   | 1         |     2      |     0     | S5Time        | time.Duration | 2          | S1200    |
| DB1.D2/DB1.DATE2           | DB   | 1         |     2      |     0     | Date          | time.Time     | 2          | S1200    |
| DB1.DT6/DB1.DATETIME6      | DB   | 1         |     6      |     0     | DateTime      | time.Time     | 8          | S1200    |
| DB1.DTL2/DB1.DATETIMELONG2 | DB   | 1         |     2      |     0     | DateTimeLong  | time.Time     | 12         | S1200    |
| DB1.TOD2/DB1.TIMEOFDAY2    | DB   | 1         |     2      |     0     | TimeOfDay     | time.Time     | 4          | S1200    |

# 🌽 License

Distributed under the MIT License. See [`LICENSE`](./LICENSE) for more information.<br>
@2024 - 2099 shiyuecamus, All Rights Reserved. <br>

❗❗❗ **Please strictly abide by the MIT agreement and add the author's copyright license notice when using.**

# 🍠 Dependencies

The dependencies used in this project are as follows:

| Number | Dependency                                | Version |  License   |     Date     | Copyright         |
|:------:|:------------------------------------------|---------|:----------:|:------------:|:------------------|
|   1    | [gnet](https://github.com/panjf2000/gnet) | 2.3.5   | Apache-2.0 | 2019-present | Andy Pan          |
|   2    | [cast](https://github.com/spf13/cast)     | 1.6.0   |    MIT     |     2014     | Steve Francia     |
|   3    | [zap](https://github.com/uber-go/zap)     | 1.27.0  |    MIT     |  2016-2017   | Uber Technologies |

## Sponsor

Buy me a cup of coffee. <br>
**WeChat** (Please note the problem or purpose you encountered

![wechat](https://i.postimg.cc/c1pfY9MT/20240315190932.jpg)
![wechat](https://i.postimg.cc/x867pXGy/20240315192946.jpg)
