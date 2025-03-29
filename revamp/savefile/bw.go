package savefile

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
)


func NewBwSavefile(bytes []byte) *BwSavefile {
	return &BwSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]*models.Pokemon, 0),
	}
}

func (bw *BwSavefile) PartyPokemon() []*models.Pokemon {
	return bw.partyPokemon
}

func (bw *BwSavefile) validate() error {
	return nil
}

func (bw *BwSavefile) Version() enums.GameVersion {
	return enums.BW
}

func (bw *BwSavefile) Flush() error {
	return nil
}
