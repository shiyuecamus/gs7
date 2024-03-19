package util

import (
	"bytes"
	"encoding/binary"
)

func NumberToBytes(n any) []byte {
	res, _ := NumberToBytesE(n)
	return res
}

func NumberToBytesE(n any) ([]byte, error) {
	bytesBuffer := bytes.NewBuffer([]byte{})

	if err := binary.Write(bytesBuffer, binary.BigEndian, n); err != nil {
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}
