package parser

import (
	"github.com/dingdongg/pkmn-rom-parser/v7/consts"
	"github.com/dingdongg/pkmn-rom-parser/v7/rom_reader"
	"github.com/dingdongg/pkmn-rom-parser/v7/rom_writer"
	"github.com/dingdongg/pkmn-rom-parser/v7/rom_writer/req"
	"github.com/dingdongg/pkmn-rom-parser/v7/validator"
	"github.com/dingdongg/pkmn-rom-parser/v7/validator/locator"
)

func Parse(savefile []byte) ([]rom_reader.Pokemon, error) {
	if err := validator.Validate(savefile); err != nil {
		return []rom_reader.Pokemon{}, err
	}
	// TODO: make ISave implement this method,
	// then call this method regardless of specific game version
	chunk := locator.GetLatestSaveChunk(savefile)
	// fmt.Println(*chunk)
	partyData := chunk.SmallBlock.BlockData[consts.PERSONALITY_OFFSET_HGSS:]

	return rom_reader.GetPartyPokemon(partyData), nil
}

func Write(savefile []byte, newBytes []req.WriteRequest) ([]byte, error) {
	if err := validator.Validate(savefile); err != nil {
		return []byte{}, err
	}

	chunk := locator.GetLatestSaveChunk(savefile)
	return rom_writer.UpdatePartyPokemon(savefile, *chunk, newBytes)
}
