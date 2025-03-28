package savefile

import "github.com/dingdongg/pkmn-rom-parser/v7/types"

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