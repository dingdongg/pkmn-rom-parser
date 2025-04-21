package utils

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func PrintBuffer(buf []byte, w int) {
	output := ""
	buffer := make([]string, 0)

	update := func() {
		buffer = append(buffer, "\n")
		output += strings.Join(buffer, " ")
		buffer = make([]string, 0)
	}

	for i, b := range buf {
		buffer = append(buffer, fmt.Sprintf("%02X", b))
		if i % w == w - 1 {
			update()
		}
	}
	if len(buffer) != 0 {
		update()
	}
	fmt.Println(output)
}

func U8(buf []byte, index int) uint8 {
	return uint8(buf[index])
}

func U16(buf []byte, index int) uint16 {
	return binary.LittleEndian.Uint16(buf[index : index+2])
}

func U32(buf []byte, index int) uint32 {
	return binary.LittleEndian.Uint32(buf[index : index+4])
}

func U64(buf []byte, index int) uint64 {
	return binary.LittleEndian.Uint64(buf[index : index+8])
}

func WriteU8(buf []byte, index int, val uint8) {
	buf[index] = val
}

func WriteU16(buf []byte, index int, val uint16) {
	binary.LittleEndian.PutUint16(buf[index : index+2], val)
}

func WriteU32(buf []byte, index int, val uint32) {
	binary.LittleEndian.PutUint32(buf[index : index+4], val)
}

func WriteU64(buf []byte, index int, val uint64) {
	binary.LittleEndian.PutUint64(buf[index : index+8], val)
}
