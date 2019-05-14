package acinbase

import (
	"bytes"
	"encoding/binary"
)

// 字符数组转成 uint16
func BytesToUint16(bits []byte) uint16 {
	var tmp uint16
	bytesBuffer := bytes.NewBuffer(bits)
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return tmp
}

// uint16字符数组转成 
func Uint16ToBytes(digit uint16) []byte {
	 bytesBuffer := bytes.NewBuffer([]byte{})
	 binary.Write(bytesBuffer, binary.BigEndian, digit)
	 return bytesBuffer.Bytes()
}
