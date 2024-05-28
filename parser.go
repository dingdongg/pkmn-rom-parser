package parser

import (
	"errors"

	"github.com/dingdongg/pkmn-rom-parser/v3/rom_reader"
	"github.com/dingdongg/pkmn-rom-parser/v3/validator"
)

const PERSONALITY_OFFSET = 0xA0

func Parse(savefile []byte) ([]rom_reader.Pokemon, error) {
	valid := validator.Validate(savefile)
	var res []rom_reader.Pokemon

	if !valid {
		return res, errors.New("invalid file")
	}

	// TODO: only read from/edit the most recent savefiel
	for i := uint(0); i < 6; i++ {
		res = append(res, rom_reader.GetPokemon(savefile[PERSONALITY_OFFSET:], i))
	}

	return res, nil
}
