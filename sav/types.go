package sav

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts"
	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
)

const PLAT_SB_END uint = uint(0xCF2C) // non-inclusive
const PLAT_BB_START uint = PLAT_SB_END + 0x0
const PLAT_BB_END uint = PLAT_BB_START + 0x121E4 // non-inclusive

const HGSS_SB_END uint = uint(0xF628)            // non-inclusive
const HGSS_BB_START uint = HGSS_SB_END + 0xD8    // padding included in hgss
const HGSS_BB_END uint = HGSS_BB_START + 0x12310 // non-inclusive

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

func identifyGameVersion(savefile []byte) (ISave, error) {
	// gen 4 games start writing to the 0x40000-offset address space,
	// check there for the existence of a valid footer
	chunkTwoOffset := uint(0x40000)
	footerSize := uint(0x14)

	if isPLAT(savefile, chunkTwoOffset, footerSize, PLAT_SB_END, PLAT_BB_END) {
		fmt.Println("Game: pokemon PLAT")
		return NewSavPLAT(savefile), nil
	} else if isHGSS(savefile, chunkTwoOffset, footerSize, HGSS_SB_END, HGSS_BB_END) {
		fmt.Println("Game: pokemon HGSS")
		return NewSavHGSS(savefile), nil
	}

	return nil, fmt.Errorf("unrecognized game file")
}

func isPLAT(savefile []byte, offset uint, footerSize uint, smallBlockEnd uint, bigBlockEnd uint) bool {
	smallFooter := getFooter(savefile[offset+smallBlockEnd-footerSize : offset+smallBlockEnd])
	bigFooter := getFooter(savefile[offset+bigBlockEnd-footerSize : offset+bigBlockEnd])

	if smallFooter.BlockSize != uint32(0xCF2C) {
		return false
	}

	if smallFooter.K != consts.MAGIC_TIMESTAMP_JP_INTL && smallFooter.K != consts.MAGIC_TIMESTAMP_KR {
		return false
	}

	if bigFooter.BlockSize != uint32(0x121E4) {
		return false
	}

	if bigFooter.K != consts.MAGIC_TIMESTAMP_JP_INTL && bigFooter.K != consts.MAGIC_TIMESTAMP_KR {
		return false
	}

	return true
}

func isHGSS(savefile []byte, offset uint, footerSize uint, smallBlockEnd uint, bigBlockEnd uint) bool {
	smallFooter := getFooter(savefile[offset+smallBlockEnd-footerSize : offset+smallBlockEnd])
	bigFooter := getFooter(savefile[offset+bigBlockEnd-footerSize : offset+bigBlockEnd])

	if smallFooter.BlockSize != uint32(0xF628) {
		return false
	}

	if smallFooter.K != consts.MAGIC_TIMESTAMP_JP_INTL && smallFooter.K != consts.MAGIC_TIMESTAMP_KR {
		return false
	}

	if bigFooter.BlockSize != uint32(0x12310) {
		return false
	}

	if bigFooter.K != consts.MAGIC_TIMESTAMP_JP_INTL && bigFooter.K != consts.MAGIC_TIMESTAMP_KR {
		return false
	}

	return true
}
