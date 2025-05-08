package savefile

import (
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v7/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/ripper"
	"github.com/dingdongg/pkmn-rom-parser/v7/validator"
	"github.com/dingdongg/pkmn-rom-parser/v7/validator/block"
)

type Savefile interface {
	/*
		1. Decrypt party pokemon section

		2. Parse stream of bytes into pokemon data
	*/
	PartyPokemon() []*models.Pokemon
	Flush() error
	Version() enums.GameVersion
	validate() error
}

func NewSavefile(bytes []byte) Savefile {
	game := validator.IdentifyGame(bytes)

	switch game {
	case enums.DP:
		return NewDpSavefile(bytes)
	case enums.PLAT:
		return NewPlatSavefile(bytes)
	case enums.HGSS:
		return NewHgssSavefile(bytes)
	case enums.BW:
		return NewBwSavefile(bytes)
	case enums.B2W2:
		return NewB2W2Savefile(bytes)
	default:
		log.Fatal("unrecognized save file")
		return nil
	}
}

type PlatSavefile struct {
	rawBytes        []byte
	latestSave      *block.Block
	partyPokemon    []*models.Pokemon
	moveNames       []string
	rawParty        []byte
	expTable        []ripper.ExperienceTable
	pokemonMetadata []ripper.PokemonMetadata
}

type DpSavefile struct {
	rawBytes     []byte
	latestSave *block.Block
	partyPokemon []*models.Pokemon
	moveNames []string
	rawParty []byte
	expTable []ripper.ExperienceTable
	pokemonMetadata []ripper.PokemonMetadata
}

type HgssSavefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}

type BwSavefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}

type B2W2Savefile struct {
	rawBytes     []byte
	partyPokemon []*models.Pokemon
}
