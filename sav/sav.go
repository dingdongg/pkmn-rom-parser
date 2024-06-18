package sav

import "github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"

// tODO: include important offsets as fields
type savPLAT struct {
	Version gamever.GameVer
	Data []byte
}

// tODO: include important offsets as fields
type savHGSS struct {
	Version gamever.GameVer
	Data []byte
}

func NewSavPLAT(savefile []byte) *savPLAT {
	return &savPLAT{
		gamever.PLAT,
		savefile,
	}
}

func NewSavHGSS(savefile []byte) *savHGSS {
	return &savHGSS{
		gamever.HGSS,
		savefile,
	}
}