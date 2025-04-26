package parser

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/rom_reader"
	"github.com/dingdongg/pkmn-rom-parser/v7/rom_writer"
	"github.com/dingdongg/pkmn-rom-parser/v7/rom_writer/req"
	"github.com/dingdongg/pkmn-rom-parser/v7/sav"
)

func Parse(savefile []byte) ([]rom_reader.Pokemon, error) {
	game, err := sav.Validate(savefile)
	if err != nil {
		return []rom_reader.Pokemon{}, err
	}

	partyData := rom_reader.GetPartyPokemon(game)
	return partyData, nil
}

func Write(savefile []byte, newBytes []req.WriteRequest) ([]byte, error) {
	game, err := sav.Validate(savefile)
	if err != nil {
		return []byte{}, err
	}

	return rom_writer.UpdatePartyPokemon(game, newBytes)
}
