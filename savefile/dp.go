package savefile

import "github.com/dingdongg/pkmn-rom-parser/v7/types"

func NewDpSavefile(bytes []byte) *DpSavefile {
	return &DpSavefile{
		rawBytes: bytes,
		partyPokemon: make([]PokemonPtr, 0),
	}
}

func (dp *DpSavefile) PartyPokemon() []PokemonPtr {
	return dp.partyPokemon
}

func (dp *DpSavefile) validate() error {
	return nil
}

func (dp *DpSavefile) Version() types.GameVersion {
	return types.DP
}

func (dp *DpSavefile) Flush() error {
	return nil
}