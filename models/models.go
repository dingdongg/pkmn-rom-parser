package models

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/enums"
)

type TODO struct{}

type StatNumber interface {
	uint | uint8 | uint16
}

type Stat[T StatNumber] struct {
	Hp         T
	Attack     T
	Defense    T
	SpeAttack  T
	SpeDefense T
	Speed      T
}

func statStringHeader(title string) string {
	statLogWidth := 49                         // 7*6 + 7 == 49
	numDashes := statLogWidth - len(title) - 2 // extra space on either side of title
	output := ""
	firstHalfWidth := numDashes >> 1
	// prepend dashes to title
	for range firstHalfWidth {
		output += "-"
	}
	output += fmt.Sprintf(" %s ", title)
	// append dashes to title
	for range numDashes - firstHalfWidth {
		output += "-"
	}

	return output
}

func (s Stat[T]) Print(title string) string {
	header := statStringHeader(title)
	ret := header
	return fmt.Sprintf("%s\n%s\n", ret, s)
}

func (s Stat[T]) String() string {
	ret := "|   Hp  | Attck | Dfnse | SpAtk | SpDef | Speed |\n"
	ret += "-------------------------------------------------\n"
	ret += fmt.Sprintf("|  %3d  | %3d   | %3d   | %3d   | %3d   | %3d   |\n",
		s.Hp, s.Attack, s.Defense, s.SpeAttack, s.SpeDefense, s.Speed)
	ret += "-------------------------------------------------\n"

	return ret
}

type Move struct {
	Id        uint16
	Name      string
	MaxPoints uint8
}

func (m Move) String() string {
	return fmt.Sprintf("'%s'", m.Name)
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
	Form      uint8
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
	ret += fmt.Sprintf("Moves:       %s\n", p.Moves)
	ret += p.EV.Print("EV")
	ret += p.IV.Print("IV")
	ret += p.Battle.Print("Battle Stats")
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
- moveset + PP ***
- IV
- gender
- forms ***

C:
- name

D:
- none

battle stats;
- level
- battle stats


from DB:
- base stats ***
*/
