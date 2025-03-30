package validator

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
)

type bytes struct {
	buf []byte
}

func IdentifyGame(buf []byte) enums.GameVersion {
	b := bytes{buf}

	if err := b.checkDiamondPearl(); err == nil {
		return enums.DP
	} else if err = b.checkPlatinum(); err == nil {
		return enums.PLAT
	} else if err = b.checkHeartGoldSoulSilver(); err == nil {
		return enums.HGSS
	} else if err = b.checkBlackWhite(); err == nil {
		return enums.BW
	} else if err = b.checkBlack2White2(); err == nil {
		return enums.B2W2
	}

	return enums.INVALID
}

func (b bytes) checkDiamondPearl() error {
	sfOffset := 0xC0EC
	sf1, sf2 := sfOffset, sfOffset + 0x40000
	sf1Count, sf2Count := sf1 + 0x4, sf2 + 0x4

	var latestOffset int
	countOne := binary.LittleEndian.Uint32(b.buf[sf1Count : sf1Count+0x4])
	countTwo := binary.LittleEndian.Uint32(b.buf[sf2Count : sf2Count+0x4])

	if countOne > countTwo {
		latestOffset = sf1
	} else if countOne < countTwo {
		latestOffset = sf2
	} else {
		return fmt.Errorf("SB save counts are same? how do I handle this?")
	}

	magicNumOffset := latestOffset + 0xC
	magicNum := binary.LittleEndian.Uint32(b.buf[magicNumOffset : magicNumOffset+0x4])

	if magicNum != enums.MAGIC_TS_JP_INTL && magicNum != enums.MAGIC_TS_KR {
		return fmt.Errorf("magic numbers invalid")
	}

	// SB checksum vlaidation
	sbSize := binary.LittleEndian.Uint32(b.buf[latestOffset+0x8 : latestOffset+0x8+0x4])
	checksumOffset := latestOffset + 0x12
	expected := binary.LittleEndian.Uint16(b.buf[checksumOffset : checksumOffset+0x2])
	actual := crypt.CRC16_CCITT(b.buf[uint32(latestOffset)-sbSize+0x14 : latestOffset])

	if expected != actual {
		return fmt.Errorf("[SMALL BLOCK] checksum mismatch: 0x%04X (expected), 0x%04X (actual)", expected, actual)
	}

	// big block validations
	bigBlockCount := binary.LittleEndian.Uint32(b.buf[latestOffset : latestOffset+0x4])

	bfOffset := 0x1E2CC
	bf1, bf2 := bfOffset, bfOffset + 0x40000
	var latestBigBlockAddr int
	bf1Count := binary.LittleEndian.Uint32(b.buf[bf1 : bf1+0x4])
	bf2Count := binary.LittleEndian.Uint32(b.buf[bf2 : bf2+0x4])

	if bigBlockCount == bf1Count {
		latestBigBlockAddr = bf1
	} else if bigBlockCount == bf2Count {
		latestBigBlockAddr = bf2
	} else {
		return fmt.Errorf("no big block match found")
	}

	bbSize := binary.LittleEndian.Uint32(b.buf[latestBigBlockAddr+0x8 : latestBigBlockAddr+0x8+0x4])
	bbChecksumOffset := latestBigBlockAddr + 0x12
	expected = binary.LittleEndian.Uint16(b.buf[bbChecksumOffset : bbChecksumOffset+0x2])
	actual = crypt.CRC16_CCITT(b.buf[uint32(latestBigBlockAddr)-bbSize+0x14 : latestBigBlockAddr])

	if expected != actual {
		return fmt.Errorf("[BIG BLOCK] checksum mismatch: 0x%04X (expected), 0x%04X (actual)", expected, actual)
	}

	return nil
}

func (b bytes) checkPlatinum() error {
		// just validate the latest block pair for now
	/*
	checksum validation
	magic number in correct place 
	
	*/
	sfOffset := 0xCF18
	sf1, sf2 := sfOffset, sfOffset + 0x40000
	sf1Count, sf2Count := sf1 + 0x4, sf2 + 0x4

	// find latest 
	var latestOffset int
	countOne := binary.LittleEndian.Uint32(b.buf[sf1Count : sf1Count+0x4])
	countTwo := binary.LittleEndian.Uint32(b.buf[sf2Count : sf2Count+0x4])

	if countOne > countTwo {
		latestOffset = sf1
	} else if countOne < countTwo {
		latestOffset = sf2
	} else {
		return fmt.Errorf("SB save counts are same? how do I handle this?")
	}

	// check magic number
	magicNumOffset := latestOffset + 0xC
	magicNum := binary.LittleEndian.Uint32(b.buf[magicNumOffset : magicNumOffset+0x4])

	if magicNum != enums.MAGIC_TS_JP_INTL && magicNum != enums.MAGIC_TS_KR {
		return fmt.Errorf("magic numbers invalid")
	}

	// SB checksum vlaidation
	sbSize := binary.LittleEndian.Uint32(b.buf[latestOffset+0x8 : latestOffset+0x8+0x4])
	checksumOffset := latestOffset + 0x12
	expected := binary.LittleEndian.Uint16(b.buf[checksumOffset : checksumOffset+0x2])
	actual := crypt.CRC16_CCITT(b.buf[uint32(latestOffset)-sbSize+0x14 : latestOffset])

	if expected != actual {
		return fmt.Errorf("[SMALL BLOCK] checksum mismatch: 0x%04X (expected), 0x%04X (actual)", expected, actual)
	}

	// big block validations
	bigBlockCount := binary.LittleEndian.Uint32(b.buf[latestOffset : latestOffset+0x4])

	bfOffset := 0x1F0FC
	bf1, bf2 := bfOffset, bfOffset + 0x40000
	var latestBigBlockAddr int
	bf1Count := binary.LittleEndian.Uint32(b.buf[bf1 : bf1+0x4])
	bf2Count := binary.LittleEndian.Uint32(b.buf[bf2 : bf2+0x4])

	if bigBlockCount == bf1Count {
		latestBigBlockAddr = bf1
	} else if bigBlockCount == bf2Count {
		latestBigBlockAddr = bf2
	} else {
		return fmt.Errorf("no big block match found")
	}

	bbSize := binary.LittleEndian.Uint32(b.buf[latestBigBlockAddr+0x8 : latestBigBlockAddr+0x8+0x4])
	bbChecksumOffset := latestBigBlockAddr + 0x12
	expected = binary.LittleEndian.Uint16(b.buf[bbChecksumOffset : bbChecksumOffset+0x2])
	actual = crypt.CRC16_CCITT(b.buf[uint32(latestBigBlockAddr)-bbSize+0x14 : latestBigBlockAddr])

	if expected != actual {
		return fmt.Errorf("[BIG BLOCK] checksum mismatch: 0x%04X (expected), 0x%04X (actual)", expected, actual)
	}

	return nil
}

func (b bytes) checkHeartGoldSoulSilver() error {
	return nil
}

func (b bytes) checkBlackWhite() error {
	return nil
}

func (b bytes) checkBlack2White2() error {
	return nil
}