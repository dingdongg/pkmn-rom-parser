package parser

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v2/rom_reader"
	"github.com/dingdongg/pkmn-rom-parser/v2/validator"
)

const PERSONALITY_OFFSET = 0xA0

func Parse(savefile []byte) []rom_reader.Pokemon {
	valid := validator.Validate(savefile)

	if valid {
		fmt.Println("SAVEFILE IS VALID")
	}

	// TODO: only read from/edit the most recent savefiel
	var res []rom_reader.Pokemon

	for i := uint(0); i < 6; i++ {
		res = append(res, rom_reader.GetPokemon(savefile[PERSONALITY_OFFSET:], i))
	}

	return res
}
