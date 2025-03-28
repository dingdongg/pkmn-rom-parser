package savefile

import "github.com/dingdongg/pkmn-rom-parser/v7/types"

func NewBwSavefile(bytes []byte) *BwSavefile {
	return &BwSavefile{
		rawBytes: bytes,
		partyPokemon: make([]PokemonPtr, 0),
	}
}

func (bw *BwSavefile) PartyPokemon() []PokemonPtr {
	return bw.partyPokemon
}

func (bw *BwSavefile) validate() error {
	return nil
}

func (bw *BwSavefile) Version() types.GameVersion {
	return types.BW
}

func (bw *BwSavefile) Flush() error {
	return nil
}