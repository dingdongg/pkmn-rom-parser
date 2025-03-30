package validator

import "github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"

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
	return validateGen4(b.buf, 0xC0EC, 0x1E2CC)
}

func (b bytes) checkPlatinum() error {
	return validateGen4(b.buf, 0xCF18, 0x1F0FC)
}

func (b bytes) checkHeartGoldSoulSilver() error {
	return validateGen4(b.buf, 0xF614, 0x219FC)
}

func (b bytes) checkBlackWhite() error {
	return nil
}

func (b bytes) checkBlack2White2() error {
	return nil
}