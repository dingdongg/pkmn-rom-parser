package sav

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
)

type Chunk struct {
	SmallBlock Block
	BigBlock   Block
}

type Block struct {
	BlockData []byte
	Footer    Footer
	Address   uint
}

type Footer struct {
	Identifier uint32
	SaveNumber uint32
	BlockSize  uint32
	K          uint32
	T          uint16
	Checksum   uint16
}

func getFooter(buf []byte) Footer {
	return Footer{
		binary.LittleEndian.Uint32(buf[:0x4]),
		binary.LittleEndian.Uint32(buf[0x4:0x8]),
		binary.LittleEndian.Uint32(buf[0x8:0xC]),
		binary.LittleEndian.Uint32(buf[0xC:0x10]),
		binary.LittleEndian.Uint16(buf[0x10:0x12]),
		binary.LittleEndian.Uint16(buf[0x12:0x14]),
	}
}

func NewBlock(data []byte, footer []byte, startAddr uint) Block {
	return Block{data, getFooter(footer), startAddr}
}

func (b Block) String() string {
	return fmt.Sprintf(`
	Block {
		data: % +x...,
		    %s,
		address: 0x%x,
	}`, b.BlockData[0:0x10], b.Footer, b.Address)
}

// footer format specifier
func (f Footer) String() string {
	return fmt.Sprintf(`
	footer {
		identifier = 0x%x,
		saveNumber = 0x%x,
		blockSize = 0x%x,
		K = 0x%x,
		T = 0x%x,
		checksum = 0x%x,
	}`, f.Identifier, f.SaveNumber, f.BlockSize, f.K, f.T, f.Checksum)
}

func (c Chunk) IsValid() bool {
	smallChecksum := crypt.CRC16_CCITT(c.SmallBlock.BlockData)
	fmt.Println("checking small block checksum")
	if smallChecksum != c.SmallBlock.Footer.Checksum {
		return false
	}
	fmt.Println("checking big block checksum")

	bigChecksum := crypt.CRC16_CCITT(c.BigBlock.BlockData)
	fmt.Println("success")
	return bigChecksum == c.BigBlock.Footer.Checksum
}