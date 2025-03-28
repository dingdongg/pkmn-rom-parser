package savefile

import "github.com/dingdongg/pkmn-rom-parser/v7/types"


type Pokemon struct {
	Level int
	Name string
}

type PokemonPtr *Pokemon

type Savefile interface {
	PartyPokemon() []PokemonPtr
	Flush() error
	Version() types.GameVersion
	validate() error
}

func NewSavefile(bytes []byte) Savefile {
	// savefile, err := identifyVersion(bytes)
	return NewPlatSavefile(bytes)
}

type PlatSavefile struct {
	rawBytes []byte
	partyPokemon []PokemonPtr
}

type DpSavefile struct {
	rawBytes []byte
	partyPokemon []PokemonPtr
}

type HgssSavefile struct {
	rawBytes []byte
	partyPokemon []PokemonPtr
}

type BwSavefile struct {
	rawBytes []byte
	partyPokemon []PokemonPtr
}

type B2W2Savefile struct {
	rawBytes []byte
	partyPokemon []PokemonPtr
}