package sav

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
)

type ISave interface {
	GetChunk(offset uint) Chunk
	Validate() error
	GetLatestData() *Chunk
	GetPartySection() []byte
	GetPartySize() uint32
}

type gen4Savefile struct {
	version        gamever.GameVer
	data           []byte
	smallBlockSize uint
	bigBlockSize   uint
	partyOffset    uint
}

// tODO: include important offsets as fields
type savPLAT gen4Savefile
type savHGSS gen4Savefile

func Validate(savefile []byte) (ISave, error) {
	game, err := identifyGameVersion(savefile)
	if err != nil {
		return nil, err
	}

	if err = game.Validate(); err != nil {
		return nil, err
	}

	return game, nil
}
