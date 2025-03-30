package savefile

import (
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
)

type Savefile interface {
	/*
		1. Decrypt party pokemon section

		2. Parse stream of bytes into pokemon data
	*/
	PartyPokemon() []*models.Pokemon
	Flush() error
	Version() enums.GameVersion
	validate() error
}

func NewSavefile(bytes []byte) Savefile {
	// savefile, err := identifyVersion(bytes)

	/*
	I don't like the idea of instantiating a savefile,
	validating it, and then failing

	flow should be:
	1. identify type of game
	2. validate that step 1 is correct
	3. if correct, instantiate savefile for that game; otherwise go next 
	*/
	ret := NewPlatSavefile(bytes)
	if err := ret.validate(); err != nil {
		log.Fatal(err)
	}
	return NewPlatSavefile(bytes)
}

type PlatSavefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}

type DpSavefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}

type HgssSavefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}

type BwSavefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}

type B2W2Savefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}
