package savefile

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
)


func NewDpSavefile(bytes []byte) *DpSavefile {
	return &DpSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]*models.Pokemon, 0),
	}
}

func (dp *DpSavefile) PartyPokemon() []*models.Pokemon {
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
