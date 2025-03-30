package models

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
)

type TODO struct{}


type Stat[T any] struct {
	Hp T
	Attack T
	Defense T
	SpeAttack T
	SpeDefense T
	Speed T
}

func (s Stat[T]) String() string {
	ret :=             "|   Hp  | Attck | Dfnse | SpAtk | SpDef | Speed |\n"
	ret +=             "-------------------------------------------------\n"
	ret += fmt.Sprintf("|  %3d  | %3d   | %3d   | %3d   | %3d   | %3d   |\n", 
					   s.Hp, s.Attack, s.Defense, s.SpeAttack, s.SpeDefense, s.Speed)
	ret +=             "-------------------------------------------------\n"

	return ret
}

type Move struct {
	Id uint16
	Name string
	MaxPoints uint8
}

type Pokemon struct {
	Name      string
	Level     uint8
	Exp       uint32
	PokedexId uint16
	Nature    enums.Nature
	Ability   string
	HeldItem  string
	Moves     []Move
	Base      Stat[uint8]
	Battle    Stat[uint16]
	EV        Stat[uint8]
	IV        Stat[uint8]
	Gender    enums.Gender
	Form      TODO
}

func (p *Pokemon) String() string {
	ret := "-----\n"
	ret += fmt.Sprintf("Name:        '%s'\n", p.Name)
	ret += fmt.Sprintf("Level:       %d\n", p.Level)
	ret += fmt.Sprintf("Exp. points: %d\n", p.Exp)
	ret += fmt.Sprintf("Pokedex ID:  %d\n", p.PokedexId)
	ret += fmt.Sprintf("Nature:      %s\n", p.Nature) // turn this into string for logging purposes
	ret += fmt.Sprintf("Ability:     '%s'\n", p.Ability)
	ret += fmt.Sprintf("Held Item:   '%s'\n", p.HeldItem)
	ret += fmt.Sprintf("Gender:      %s\n", p.Gender)
	ret += fmt.Sprintf("----------------------- EV ----------------------\n%s\n", p.EV)
	ret += fmt.Sprintf("----------------------- IV ----------------------\n%s\n", p.IV)
	ret += fmt.Sprintf("------------------ Battle Stats -----------------\n%s\n", p.Battle)
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
