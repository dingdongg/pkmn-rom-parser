package savefile

import "github.com/dingdongg/pkmn-rom-parser/v7/types"

func NewHgssSavefile(bytes []byte) *HgssSavefile {
	return &HgssSavefile{
		rawBytes: bytes,
		partyPokemon: make([]PokemonPtr, 0),
	}
}

func (hgss *HgssSavefile) PartyPokemon() []PokemonPtr {
	return hgss.partyPokemon
}

func (hgss *HgssSavefile) validate() error {
	return nil
}

func (hgss *HgssSavefile) Version() types.GameVersion {
	return types.HGSS
}

func (hgss *HgssSavefile) Flush() error {
	return nil
}