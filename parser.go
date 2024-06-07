package parser

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v3/rom_reader"
	"github.com/dingdongg/pkmn-rom-parser/v3/validator"
	"github.com/dingdongg/pkmn-rom-parser/v3/validator/locator"
)

const PERSONALITY_OFFSET = 0xA0

func Parse(savefile []byte) ([]rom_reader.Pokemon, error) {
	var res []rom_reader.Pokemon

	if err := validator.Validate(savefile); err != nil {
		return res, err
	}

	chunk := locator.GetLatestSaveChunk(savefile)
	partyData := chunk.SmallBlock.BlockData[PERSONALITY_OFFSET:]

	fmt.Println(chunk.SmallBlock.Footer)
	fmt.Println(chunk.BigBlock.Footer)

	for i := uint(0); i < 6; i++ {
		res = append(res, rom_reader.GetPokemon(partyData, i))
	}

	return res, nil
}
