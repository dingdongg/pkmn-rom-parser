package validator

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts"
	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
)

/*
POKEMON PLATINUM
1. savfile size = 2^19 bytes
2. checksums are used for validating (part pokemon data) (per-pokemon basis)


2 types of validation
- validating an entire big/small block as a whole
- validating each party pokemon data individually

*/

// a "chunk" denotes a pair of small + big block that are adjacent in memory
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

const savefileSize int = 1 << 19
const footerSize uint = 0x14
const secondChunkOffset uint = 0x40000

const PLAT_SB_END uint = uint(0xCF2C) // non-inclusive
const PLAT_BB_START uint = PLAT_SB_END + 0x0
const PLAT_BB_END uint = PLAT_BB_START + 0x121E4 // non-inclusive

const HGSS_SB_END uint = uint(0xF628) // non-inclusive
const HGSS_BB_START uint = HGSS_SB_END + 0xD8 // padding included in hgss
const HGSS_BB_END uint = HGSS_BB_START + 0x12310 // non-inclusive

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
	return Block{ data, getFooter(footer), startAddr }
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

func GetChunk(savefile []byte, offset uint) Chunk {
	smallBlockFooterAddr := uint(0x0CF18) + offset
	bigBlockFooterAddr := uint(0x1F0FC) + offset
	bigBlockStart := uint(0xCF2C) + offset

	smallBlock := Block{
		savefile[offset:smallBlockFooterAddr],
		getFooter(savefile[smallBlockFooterAddr : smallBlockFooterAddr+footerSize]),
		offset,
	}

	bigBlock := Block{
		savefile[bigBlockStart:bigBlockFooterAddr],
		getFooter(savefile[bigBlockFooterAddr : bigBlockFooterAddr+footerSize]),
		bigBlockStart,
	}

	return Chunk{smallBlock, bigBlock}
}

// validates the given .sav file
func Validate(savefile []byte) error {
	if len(savefile) != savefileSize {
		return errors.New("invalid savefile")
	}

	// identify the game version here
	gameVersion, err := identifyGameVersion(savefile)
	if err != nil {
		return errors.New("unrecognized game file")
	}

	var res error = nil

	// choose validation strategy based on game version
	// ^ STRATEGY PATTERN????? (especially when i introduce suoppport for gen 5)
	switch (gameVersion) {
	case gamever.PLAT:
		res = validatePLAT(savefile)
	case gamever.HGSS:
		res = validateHGSS(savefile)
	default:
		return errors.New("unrecognized/unsupported game file")
	}

	return res
}

func identifyGameVersion(savefile []byte) (gamever.GameVer, error) {
	// gen 4 games start writing to the 0x40000-offset address space,
	// check there for the existence of a valid footer

	if isPLAT(savefile, secondChunkOffset, footerSize, PLAT_SB_END, PLAT_BB_END) {
		fmt.Println("Game: pokemon PLAT")
		return gamever.PLAT, nil
	} else if isHGSS(savefile, secondChunkOffset, footerSize, HGSS_SB_END, HGSS_BB_END) {
		fmt.Println("Game: pokemon HGSS")
		return gamever.HGSS, nil
	}

	return -1, errors.New("unrecognized game file")
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

func (c Chunk) log() {
	fmt.Println(c.SmallBlock)
	fmt.Println(c.BigBlock)
	fmt.Println("-------------")
}

func validatePLAT(savefile []byte) error {
	firstChunk := GetChunk(savefile, 0)
	secondChunk := GetChunk(savefile, secondChunkOffset)

	firstChunk.log()
	secondChunk.log()

	if !firstChunk.IsValid() {
		fmt.Println("First chunk invalid")
		return errors.New("invalid savefile")
	}

	if !secondChunk.IsValid() {
		fmt.Println("Second chunk invalid")
		return errors.New("invalid savefile")
	}

	return nil
}

func validateHGSS(savefile []byte) error {
	firstChunk := GetChunkHGSS(savefile, 0)
	secondChunk := GetChunkHGSS(savefile, secondChunkOffset)

	firstChunk.log()

	// it's possible that this first chunk hasn't been initialized yet,
	// (if the player just started and only saved once)
	// TODO: figure out a way around this
	if !firstChunk.IsValid() {
		fmt.Println("First chunk invalid")
		return errors.New("invalid savefile")
	}

	secondChunk.log()

	if !secondChunk.IsValid() {
		fmt.Println("Second chunk invalid")
		return errors.New("invalid savefile")
	}

	return nil
}

func GetChunkHGSS(savefile []byte, offset uint) Chunk {
	// fmt.Printf("HGSS: getting chunk at offset 0x%x\n", offset)
	smallBlockFooterAddr := HGSS_SB_END + offset - footerSize
	bigBlockFooterAddr := HGSS_BB_END + offset - footerSize

	smallBlock := Block{
		savefile[offset : smallBlockFooterAddr+0x4], // for some reason, footer identifier is included in the checksum calculations
		getFooter(savefile[smallBlockFooterAddr : smallBlockFooterAddr+footerSize]),
		offset,
	}

	bigBlock := Block{
		savefile[offset+HGSS_BB_START : bigBlockFooterAddr+0x4], // for some reason, footer identifier is included in the checksum calculations
		getFooter(savefile[bigBlockFooterAddr : bigBlockFooterAddr+footerSize]),
		offset+HGSS_BB_START,
	}

	return Chunk{smallBlock, bigBlock}
}