package sav

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
)

type Savefile interface {
	PartyPokemon() any
	Flush() error
	Version() enums.GameVersion
	validate() error
}

type gen4Savefile struct {
	version        gamever.GameVer
	data           []byte
	smallBlockSize uint
	bigBlockSize   uint
	partyOffset    uint
}

type gen5Savefile struct {
	version        gamever.GameVer
	data           []byte
	smallBlockSize uint
	bigBlockSize   uint
	partyOffset    uint
}

// tODO: include important offsets as fields
type savPLAT gen4Savefile
type savHGSS gen4Savefile
type savBW gen5Savefile

func Validate(savefile []byte) (Savefile, error) {
	game, err := identifyGameVersion(savefile)
	if err != nil {
		return nil, err
	}

	if err = game.validate(); err != nil {
		return nil, err
	}

	return game, nil
}
