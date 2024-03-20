// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package gs7

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/shiyuecamus/gs7/common"
	"github.com/shiyuecamus/gs7/util"
	"github.com/spf13/cast"
	"time"
	"unicode/utf16"
)

type Bit bool

func BitFromBytes(bs []byte) (b Bit, err error) {
	if len(bs) < 1 {
		err = errors.New("invalid bytes for Bit")
		return
	}
	b = Bit(util.GetBoolAt(bs[0], 0))
	return
}

func (b Bit) ToBytes() []byte {
	return []byte{util.SetBoolAt(0x00, 0, bool(b))}
}

func (b Bit) String() string {
	return "Bit[" + cast.ToString(bool(b)) + "]"
}

type Byte uint8

func ByteFromBytes(bs []byte) (b Byte, err error) {
	if len(bs) < 1 {
		err = errors.New("invalid bytes for Byte")
		return
	}
	var i uint8
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &i)
	b = Byte(i)
	return
}

func (b Byte) ToBytes() []byte {
	return []byte{uint8(b)}
}

func (b Byte) String() string {
	return "Byte[" + cast.ToString(uint8(b)) + "]"
}

type Char int8

func CharFromBytes(bs []byte) (c Char, err error) {
	if len(bs) < 1 {
		err = errors.New("invalid bytes for Char")
		return
	}
	var i int8
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &i)
	c = Char(i)
	return
}

func (c Char) ToBytes() []byte {
	return []byte{byte(c)}
}

func (c Char) String() string {
	return "Char[" + cast.ToString(int8(c)) + "]"
}

type Int int16

func IntFromBytes(bs []byte) (i Int, err error) {
	if len(bs) < 2 {
		err = errors.New("invalid bytes for Int")
		return
	}
	var i16 int16
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &i16)
	i = Int(i16)
	return
}

func (i Int) ToBytes() []byte {
	return util.NumberToBytes(int16(i))
}

func (i Int) String() string {
	return "Int[" + cast.ToString(int16(i)) + "]"
}

type Word uint16

func WordFromBytes(bs []byte) (w Word, err error) {
	if len(bs) < 2 {
		err = errors.New("invalid bytes for Int")
		return
	}
	var i uint16
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &i)
	w = Word(i)
	return
}

func (w Word) ToBytes() []byte {
	return util.NumberToBytes(uint16(w))
}

func (w Word) String() string {
	return "Word[" + cast.ToString(uint16(w)) + "]"
}

type DInt int32

func DIntFromBytes(bs []byte) (d DInt, err error) {
	if len(bs) < 4 {
		err = errors.New("invalid bytes for DInt")
		return
	}
	var i int32
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &i)
	d = DInt(i)
	return
}

func (d DInt) ToBytes() []byte {
	return util.NumberToBytes(int32(d))
}

func (d DInt) String() string {
	return "DInt[" + cast.ToString(int32(d)) + "]"
}

type DWord uint32

func DWordFromBytes(bs []byte) (d DWord, err error) {
	if len(bs) < 4 {
		err = errors.New("invalid bytes for DInt")
		return
	}
	var i uint32
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &i)
	d = DWord(i)
	return
}

func (d DWord) ToBytes() []byte {
	return util.NumberToBytes(uint32(d))
}

func (d DWord) String() string {
	return "DWord[" + cast.ToString(uint32(d)) + "]"
}

type Real float32

func RealFromBytes(bs []byte) (r Real, err error) {
	if len(bs) < 4 {
		err = errors.New("invalid bytes for Real")
		return
	}
	var f float32
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &f)
	r = Real(f)
	return
}

func (r Real) ToBytes() []byte {
	return util.NumberToBytes(float32(r))
}

func (r Real) String() string {
	return "Real[" + cast.ToString(float32(r)) + "]"
}

type Time time.Duration

func TimeFromBytes(bs []byte) (t Time, err error) {
	if len(bs) < 4 {
		err = errors.New("invalid bytes for Time")
		return
	}
	var i int32
	err = binary.Read(bytes.NewReader(bs), binary.BigEndian, &i)
	t = Time(time.Duration(i) * time.Millisecond)
	return
}

func (t Time) ToBytes() []byte {
	return util.NumberToBytes(uint32(time.Duration(t).Milliseconds()))
}

func (t Time) String() string {
	return "Time[" + cast.ToString(time.Duration(t).Milliseconds()) + "ms]"
}

type Date time.Time

func DateFromBytes(bs []byte) (d Date, err error) {
	if len(bs) < 2 {
		err = errors.New("invalid bytes for Date")
		return
	}
	d = Date(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC).
		Add(time.Hour * 24 * time.Duration(binary.BigEndian.Uint16(bs))))
	return
}

func (d Date) ToBytes() []byte {
	duration := time.Time(d).Sub(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	days := uint16(duration.Hours() / 24)
	return util.NumberToBytes(days)
}

func (d Date) String() string {
	return "Date[" + time.Time(d).Format("2006-01-02") + "]"
}

type TimeOfDay time.Time

func TimeOfDayFromBytes(bs []byte) (t TimeOfDay, err error) {
	if len(bs) < 4 {
		err = errors.New("invalid bytes for TimeOfDay")
		return
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	t = TimeOfDay(today.Add(time.Duration(binary.BigEndian.Uint32(bs)) * time.Millisecond))
	return
}

func (tod TimeOfDay) ToBytes() []byte {
	t := time.Time(tod)
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	duration := t.Sub(start)
	milliseconds := duration.Milliseconds()
	return util.NumberToBytes(uint32(milliseconds))
}

func (tod TimeOfDay) String() string {
	return "TimeOfDay[" + time.Time(tod).Format("15:04:05.000") + "]"
}

type DateTime time.Time

func DateTimeFromBytes(bs []byte) (d DateTime, err error) {
	if len(bs) < 8 {
		err = errors.New("invalid bytes for DateTime")
		return
	}
	year := DecodeBcd(bs[0])
	if year < 90 {
		year = year + 2000
	} else {
		year += 1900
	}
	month := DecodeBcd(bs[1])
	day := DecodeBcd(bs[2])
	hour := DecodeBcd(bs[3])
	minute := DecodeBcd(bs[4])
	second := DecodeBcd(bs[5])
	d = DateTime(time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC))
	return
}

func (d DateTime) ToBytes() []byte {
	t := time.Time(d)
	bs := make([]byte, 8)
	if t.Year() < 2000 {
		bs[0] = EncodeBcd(t.Year() - 1900)
	} else {
		bs[0] = EncodeBcd(t.Year() - 2000)
	}
	bs[1] = EncodeBcd(int(t.Month()))
	bs[2] = EncodeBcd(t.Day())
	bs[3] = EncodeBcd(t.Hour())
	bs[4] = EncodeBcd(t.Minute())
	bs[5] = EncodeBcd(t.Second())
	bs[6] = 0
	bs[7] = 0
	return bs
}

func (d DateTime) String() string {
	return "DateTime[" + time.Time(d).Format("2006-01-02 15:04:05") + "]"
}

type DateTimeLong time.Time

func DateTimeLongFromBytes(bs []byte) (d DateTimeLong, err error) {
	if len(bs) < 12 {
		err = errors.New("invalid bytes for DateTimeLong")
		return
	}
	year := binary.BigEndian.Uint16(bs[:2])
	month := bs[2]
	day := bs[3]
	hour := bs[5]
	minute := bs[6]
	second := bs[7]
	nano := binary.BigEndian.Uint32(bs[8:])
	d = DateTimeLong(time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), int(nano), time.UTC))
	return
}

func (d DateTimeLong) ToBytes() []byte {
	t := time.Time(d)
	bs := make([]byte, 0)
	year, month, day := t.Date()
	_, week := t.ISOWeek()
	hour, minute, second := t.Clock()
	nano := t.Nanosecond()
	bs = append(bs, util.NumberToBytes(uint16(year))...)
	bs = append(bs, uint8(month), uint8(day), uint8(week), uint8(hour), uint8(minute), uint8(second))
	bs = append(bs, util.NumberToBytes(uint32(nano))...)
	return bs
}

func (d DateTimeLong) String() string {
	return "DateTime[" + time.Time(d).Format("2006-01-02 15:04:05") + "]"
}

type S5Time time.Duration

func S5TimeFromBytes(bs []byte) (s S5Time, err error) {
	if len(bs) < 2 {
		err = errors.New("invalid bytes for S5Time")
		return
	}
	b := DecodeBcd(bs[0]&0b00001111)*100 + DecodeBcd(bs[1])
	switch bs[0] & 0b00110000 {
	case 0b00000000:
		b *= 10
	case 0b00010000:
		b *= 100
	case 0b00100000:
		b *= 1000
	case 0b00110000:
		b *= 10000
	}
	s = S5Time(time.Duration(b) * time.Millisecond)
	return
}

func (s S5Time) ToBytes() []byte {
	bs := make([]byte, 2)
	ms := time.Duration(s).Milliseconds()
	switch {
	case ms < 9990:
		bs[1] = EncodeBcd(int(ms) / 10 % 100)
		bs[0] = EncodeBcd(int(ms)/10/100) &^ 0b11110000
	case ms > 100 && ms < 99900:
		bs[1] = EncodeBcd(int(ms) / 100 % 100)
		bs[0] = EncodeBcd(int(ms)/100/100)&^0b11100000 | 0b00010000
	case ms > 1000 && ms < 999000:
		bs[1] = EncodeBcd(int(ms) / 1000 % 100)
		bs[0] = EncodeBcd(int(ms)/1000/100)&^0b11010000 | 0b00100000
	case ms > 10000 && ms < 9990000:
		bs[1] = EncodeBcd(int(ms) / 10000 % 100)
		bs[0] = EncodeBcd(int(ms)/10000/100)&^0b11000000 | 0b00110000
	}
	return bs
}

func (s S5Time) String() string {
	return "S5Time[" + cast.ToString(time.Duration(s).Milliseconds()) + "ms]"
}

type String string

func StringFromBytes(bs []byte, plcType common.PlcType) (s String, err error) {
	var minLen int
	switch plcType {
	case common.S200Smart:
		minLen = 1
		break
	default:
		minLen = 2
		break
	}
	if len(bs) < minLen {
		err = errors.New("invalid bytes for String")
		return
	}
	if len(bs) == minLen {
		s = ""
		return
	}
	var sub []byte
	if plcType == common.S200Smart {
		sub = bs[minLen:]
	} else {
		sub = bs[minLen:]
	}
	s = String(sub)
	return
}

func (s String) ToBytes(pduLength int) []byte {
	bs := make([]byte, 0)
	bs = append(bs, strMaxLength(pduLength), byte(len(string(s))))
	bs = append(bs, []byte(s)...)
	return bs
}

func (s String) String() string {
	return "String[" + string(s) + "]"
}

type WString string

func WStringFromBytes(bs []byte, plcType common.PlcType) (s WString, err error) {
	var minLen int
	var length int
	switch plcType {
	case common.S200Smart:
		minLen = 2
		length = int(bs[minLen/2])
		break
	default:
		minLen = 4
		length = int(binary.BigEndian.Uint16(bs[minLen/2:]))
		break
	}
	if len(bs) < minLen {
		err = errors.New("invalid bytes for String")
		return
	}
	if len(bs) == minLen {
		s = ""
		return
	}
	content := bs[minLen : minLen+int(length)*2]
	u16s := make([]uint16, len(content)/2)
	for i := 0; i < len(u16s); i++ {
		u16s[i] = binary.BigEndian.Uint16(bs[minLen+i*2:])
	}
	s = WString(utf16.Decode(u16s))
	return
}

func (s WString) ToBytes(pduLength int) []byte {
	bs := make([]byte, 0)
	bs = append(bs, util.NumberToBytes(uint16(strMaxLength(pduLength)))...)

	u16s := utf16.Encode([]rune(s))
	buf := &bytes.Buffer{}
	for _, u16 := range u16s {
		_ = binary.Write(buf, binary.BigEndian, u16)
	}
	bs = append(bs, util.NumberToBytes(uint16(len(buf.Bytes()))/2)...)
	bs = append(bs, buf.Bytes()...)
	return bs
}

func (s WString) String() string {
	return "String[" + string(s) + "]"
}

type Counter uint16

func CounterFromBytes(bs []byte) (c Counter, err error) {
	if len(bs) < 2 {
		err = errors.New("invalid bytes for Counter")
		return
	}
	c = Counter((bs[0] << 8) | bs[1])
	return
}

func (c Counter) ToBytes() []byte {
	bs := make([]byte, 2)
	bs[0] = byte((uint16(c) << 8) & 0xFF)
	bs[1] = byte(uint16(c) & 0xFF)
	return bs
}

func (c Counter) String() string {
	return "Counter[" + cast.ToString(uint16(c)) + "]"
}

type Timer time.Duration

func TimerFromBytes(bs []byte) (t Timer, err error) {
	if len(bs) < 2 {
		err = errors.New("invalid bytes for Timer")
		return
	}
	b := DecodeBcd(bs[0]&0b00001111)*100 + DecodeBcd(bs[1])
	switch bs[0] & 0b00110000 {
	case 0b00000000:
		b *= 10
	case 0b00010000:
		b *= 100
	case 0b00100000:
		b *= 1000
	case 0b00110000:
		b *= 10000
	}
	t = Timer(time.Duration(b) * time.Millisecond)
	return
}

func (t Timer) ToBytes() []byte {
	bs := make([]byte, 2)
	ms := time.Duration(t).Milliseconds()
	switch {
	case ms < 9990:
		bs[1] = EncodeBcd(int(ms) / 10 % 100)
		bs[0] = EncodeBcd(int(ms)/10/100) &^ 0b11110000
	case ms > 100 && ms < 99900:
		bs[1] = EncodeBcd(int(ms) / 100 % 100)
		bs[0] = EncodeBcd(int(ms)/100/100)&^0b11100000 | 0b00010000
	case ms > 1000 && ms < 999000:
		bs[1] = EncodeBcd(int(ms) / 1000 % 100)
		bs[0] = EncodeBcd(int(ms)/1000/100)&^0b11010000 | 0b00100000
	case ms > 10000 && ms < 9990000:
		bs[1] = EncodeBcd(int(ms) / 10000 % 100)
		bs[0] = EncodeBcd(int(ms)/10000/100)&^0b11000000 | 0b00110000
	}
	return bs
}

func (t Timer) String() string {
	return "S5Time[" + cast.ToString(time.Duration(t).Milliseconds()) + "ms]"
}

func strMaxLength(pduLength int) uint8 {
	if pduLength >= 480 {
		return 254
	}
	return 210
}

func EncodeBcd(value int) byte {
	return byte(((value / 10) << 4) | (value % 10))
}

func DecodeBcd(b byte) int {
	return int(((b >> 4) * 10) + (b & 0x0F))
}
