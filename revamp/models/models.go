package models

import "github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"

type TODO struct{}

type Pokemon struct {
	Name      string
	Level     uint8
	Exp       uint32
	PokedexId uint16
	Nature    enums.Nature
	Ability   string
	HeldItem  string
	Moves     TODO
	Base      TODO
	Battle    TODO
	EV        TODO
	IV        TODO
	Gender    enums.Gender
	Form      TODO
}