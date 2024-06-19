package sav

import (
	"fmt"
)

const PLAT_SB_END uint = uint(0xCF2C) // non-inclusive
const PLAT_BB_START uint = PLAT_SB_END + 0x0
const PLAT_BB_END uint = PLAT_BB_START + 0x121E4 // non-inclusive

const HGSS_SB_END uint = uint(0xF628)            // non-inclusive
const HGSS_BB_START uint = HGSS_SB_END + 0xD8    // padding included in hgss
const HGSS_BB_END uint = HGSS_BB_START + 0x12310 // non-inclusive

const MAGIC_TIMESTAMP_JP_INTL = 0x20060623
const MAGIC_TIMESTAMP_KR = 0x20070903

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

	if smallFooter.K != MAGIC_TIMESTAMP_JP_INTL && smallFooter.K != MAGIC_TIMESTAMP_KR {
		return false
	}

	if bigFooter.BlockSize != uint32(0x121E4) {
		return false
	}

	if bigFooter.K != MAGIC_TIMESTAMP_JP_INTL && bigFooter.K != MAGIC_TIMESTAMP_KR {
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

	if smallFooter.K != MAGIC_TIMESTAMP_JP_INTL && smallFooter.K != MAGIC_TIMESTAMP_KR {
		return false
	}

	if bigFooter.BlockSize != uint32(0x12310) {
		return false
	}

	if bigFooter.K != MAGIC_TIMESTAMP_JP_INTL && bigFooter.K != MAGIC_TIMESTAMP_KR {
		return false
	}

	return true
}
