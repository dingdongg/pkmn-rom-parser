package models

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
)

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

func (p *Pokemon) String() string {
	ret := "-----\n"
	ret += fmt.Sprintf("Name:        '%s'\n", p.Name)
	ret += fmt.Sprintf("Level:       %d\n", p.Level)
	ret += fmt.Sprintf("Exp. points: %d\n", p.Exp)
	ret += fmt.Sprintf("Pokedex ID:  %d\n", p.PokedexId)
	ret += fmt.Sprintf("Nature:      %d\n", p.Nature) // turn this into string for logging purposes
	ret += fmt.Sprintf("Ability:     '%s'\n", p.Ability)
	ret += fmt.Sprintf("Held Item:   '%s'\n", p.HeldItem)
	ret += fmt.Sprintf("Gender:      %d\n", p.Gender)
	ret += "-----\n"
	return ret
}

/*

A:
- pokedex id
- held item
- ability
- EV
- EXP

B:
- moveset + PP
- IV
- gender
- forms

C:
- name

D:
- none

battle stats;
- level
- battle stats


from DB:
- base stats
*/
