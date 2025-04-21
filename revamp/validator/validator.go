package validator

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
)

type bytes struct {
	buf []byte
}

// should this also determine which part fo the save to return? latest vs. backup?
// cuz right now were always reading whatever is at the beginning of the savefile,
// regardless of whether it might be a backup or not
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
	sbRange, bbRange := utils.NewRange[uint](0x0, 0xC0FF), utils.NewRange[uint](0xC100, 0x1E2DF)
	return validateGen4(b.buf, sbRange, bbRange)
}

func (b bytes) checkPlatinum() error {
	sbRange, bbRange := utils.NewRange[uint](0x0, 0xCF2B), utils.NewRange[uint](0xCF2C, 0x1F10F)
	return validateGen4(b.buf, sbRange, bbRange)
}

func (b bytes) checkHeartGoldSoulSilver() error {
	sbRange, bbRange := utils.NewRange[uint](0x0, 0xF627), utils.NewRange[uint](0xF700, 0x21A0F)
	return validateGen4(b.buf, sbRange, bbRange)
}

func (b bytes) checkBlackWhite() error {
	// if I want to be super thorough,

	// I should check the mirror checksums against their originals
	// - would require a mapping from OG -> mirror addresses
	// THEN I can perform CRC16 CCITT on the checksum block to see if valid

	// let's only do the second step, for now

	cbOffset := 0x23F00
	cbSize := 0x8C
	checksumAddr := 0x23F9A

	actual := crypt.CRC16_CCITT(b.buf[cbOffset : cbOffset+cbSize])
	expected := binary.LittleEndian.Uint16(b.buf[checksumAddr : checksumAddr+0x2])

	if actual == expected {
		return nil
	}

	cbOffset += 0x24000
	checksumAddr += 0x24000
	actual = crypt.CRC16_CCITT(b.buf[cbOffset : cbOffset+cbSize])
	expected = binary.LittleEndian.Uint16(b.buf[checksumAddr : checksumAddr+0x2])

	if actual == expected {
		return nil
	}

	return fmt.Errorf("both savefile sections invalid")
}

func (b bytes) checkBlack2White2() error {
	cbOffset := 0x25F00
	cbSize := 0x93
	checksumAddr := 0x25FA2

	actual := crypt.CRC16_CCITT(b.buf[cbOffset : cbOffset+cbSize])
	expected := binary.LittleEndian.Uint16(b.buf[checksumAddr : checksumAddr+0x2])

	if actual == expected {
		return nil
	}

	cbOffset += 0x26000
	checksumAddr += 0x26000
	actual = crypt.CRC16_CCITT(b.buf[cbOffset : cbOffset+cbSize])
	expected = binary.LittleEndian.Uint16(b.buf[checksumAddr : checksumAddr+0x2])

	if actual == expected {
		return nil
	}

	return fmt.Errorf("both savefile sections invalid")
}