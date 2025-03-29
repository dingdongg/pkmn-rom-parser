package savefile

import "github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"


type PokemonPtr *Pokemon

type Savefile interface {
	/*
		1. Decrypt party pokemon section

		2. Parse stream of bytes into pokemon data
	*/
	PartyPokemon() []PokemonPtr
	Flush() error
	Version() enums.GameVersion
	validate() error
}

func NewSavefile(bytes []byte) Savefile {
	// savefile, err := identifyVersion(bytes)
	return NewPlatSavefile(bytes)
}

type PlatSavefile struct {
	rawBytes     []byte
	partyPokemon []PokemonPtr
}

type DpSavefile struct {
	rawBytes     []byte
	partyPokemon []PokemonPtr
}

type HgssSavefile struct {
	rawBytes     []byte
	partyPokemon []PokemonPtr
}

type BwSavefile struct {
	rawBytes     []byte
	partyPokemon []PokemonPtr
}

type B2W2Savefile struct {
	rawBytes     []byte
	partyPokemon []PokemonPtr
}
