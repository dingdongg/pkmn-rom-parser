package sav

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
)

type ISave interface {
	Chunk(offset uint) Chunk
	Validate() error
	LatestData() *Chunk
	PartySection() []byte
	PartySize() uint32
	PartyOffset() uint
	Get(start uint, numBytes uint) []byte
	Data() []byte
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
