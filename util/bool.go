// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package util

// SetBoolAt sets a boolean (bit) within a byte at bit position
// without changing the other bits
// it returns the resulted byte
func SetBoolAt(b byte, bitPos uint, data bool) byte {
	if data {
		return b | (1 << bitPos)
	}
	return b &^ (1 << bitPos)
}

// GetBoolAt gets a boolean (bit) from a byte at position
func GetBoolAt(b byte, pos uint) bool {
	return b&(1<<pos) != 0
}

// EncodeBools Convert a boolean list to a byte array
func EncodeBools(in []bool) (out []byte) {
	var byteCount uint
	var i uint

	byteCount = uint(len(in)) / 8
	if len(in)%8 != 0 {
		byteCount++
	}

	out = make([]byte, byteCount)
	for i = 0; i < uint(len(in)); i++ {
		if in[i] {
			out[i/8] |= 0x01 << (i % 8)
		}
	}

	return
}

// DecodeBools Extract a specified number of Boolean values
func DecodeBools(quantity uint16, in []byte) (out []bool) {
	var i uint
	for i = 0; i < uint(quantity); i++ {
		out = append(out, ((in[i/8]>>(i%8))&0x01) == 0x01)
	}

	return
}
