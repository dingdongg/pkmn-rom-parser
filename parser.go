package parser

import (
	"github.com/dingdongg/pkmn-rom-parser/v5/consts"
	"github.com/dingdongg/pkmn-rom-parser/v5/rom_reader"
	"github.com/dingdongg/pkmn-rom-parser/v5/rom_writer"
	"github.com/dingdongg/pkmn-rom-parser/v5/rom_writer/req"
	"github.com/dingdongg/pkmn-rom-parser/v5/validator"
	"github.com/dingdongg/pkmn-rom-parser/v5/validator/locator"
)

func Parse(savefile []byte) ([]rom_reader.Pokemon, error) {
	if err := validator.Validate(savefile); err != nil {
		return []rom_reader.Pokemon{}, err
	}

	chunk := locator.GetLatestSaveChunk(savefile)
	partyData := chunk.SmallBlock.BlockData[consts.PERSONALITY_OFFSET:]

	return rom_reader.GetPartyPokemon(partyData), nil
}

func Write(savefile []byte, newBytes []req.WriteRequest) ([]byte, error) {
	if err := validator.Validate(savefile); err != nil {
		return []byte{}, err
	}

	chunk := locator.GetLatestSaveChunk(savefile)
	return rom_writer.UpdatePartyPokemon(savefile, *chunk, newBytes)
}
