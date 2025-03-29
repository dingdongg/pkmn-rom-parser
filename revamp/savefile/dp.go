package savefile
import "github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"


func NewDpSavefile(bytes []byte) *DpSavefile {
	return &DpSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]PokemonPtr, 0),
	}
}

func (dp *DpSavefile) PartyPokemon() []PokemonPtr {
	return dp.partyPokemon
}

func (dp *DpSavefile) validate() error {
	return nil
}

func (dp *DpSavefile) Version() enums.GameVersion {
	return enums.DP
}

func (dp *DpSavefile) Flush() error {
	return nil
}
