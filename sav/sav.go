package sav

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
	"github.com/dingdongg/pkmn-rom-parser/v7/validator"
)

type ISave interface {
	GetChunk(offset uint) validator.Chunk
	Validate() error
}

// tODO: include important offsets as fields
type savPLAT struct {
	version gamever.GameVer
	data []byte
}

// tODO: include important offsets as fields
type savHGSS struct {
	version gamever.GameVer
	data []byte
}

func NewSavPLAT(savefile []byte) *savPLAT {
	return &savPLAT{
		gamever.PLAT,
		savefile,
	}
}

func (s *savPLAT) GetChunk(offset uint) validator.Chunk {
	return validator.Chunk{}
}

func (s *savPLAT) Validate() error {
	return nil
}

func NewSavHGSS(savefile []byte) *savHGSS {
	return &savHGSS{
		gamever.HGSS,
		savefile,
	}
}

func (s *savHGSS) GetChunk(offset uint) validator.Chunk {
	return validator.Chunk{}
}

func (s *savHGSS) Validate() error {
	return nil
}