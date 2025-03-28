package savefile

import "github.com/dingdongg/pkmn-rom-parser/v7/types"

func NewB2W2Savefile(bytes []byte) *B2W2Savefile {
	return &B2W2Savefile{
		rawBytes: bytes,
		partyPokemon: make([]PokemonPtr, 0),
	}
}

func (bw *B2W2Savefile) PartyPokemon() []PokemonPtr {
	return bw.partyPokemon
}

func (bw *B2W2Savefile) validate() error {
	return nil
}

func (bw *B2W2Savefile) Version() types.GameVersion {
	return types.B2W2
}

func (bw *B2W2Savefile) Flush() error {
	return nil
}