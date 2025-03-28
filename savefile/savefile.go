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

func NewPlatSavefile(bytes []byte) *PlatSavefile {
	return &PlatSavefile{
		rawBytes: bytes,
		partyPokemon: make([]PokemonPtr, 0),
	}
}

func (pt *PlatSavefile) PartyPokemon() []PokemonPtr {
	return pt.partyPokemon
}

func (pt *PlatSavefile) validate() error {
	return nil
}

func (pt *PlatSavefile) Version() types.GameVersion {
	return types.PLAT
}

func (pt *PlatSavefile) Flush() error {
	return nil
}