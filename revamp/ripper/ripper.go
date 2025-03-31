package ripper

import (
	"fmt"
	"os"

	"github.com/dingdongg/pkmn-rom-parser/v7/data"
	"github.com/dingdongg/pkmn-rom-parser/v7/path_resolver"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/ripper/narc"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
)


type PokemonMetadata struct {
	Base models.Stat[uint8]
	Type1 uint8
	Type2 uint8
	CatchRate uint8
	ExpYield uint8
	EVYield models.Stat[uint8]
	Item1 uint16
	Item2 uint16
	GenderThreshold uint8
	EggCycles uint8
	BaseFriendship uint8
	GrowthType uint8
	EggGroup1 uint8
	EggGroup2 uint8
	Ability1 uint8
	Ability2 uint8
	SafariZoneRate uint8
	Color uint8
	Padding1 uint16
	MoveFlags []byte // length 13
	Padding2 []byte // length 3
}

func (pm PokemonMetadata) String() string {
	ret := "========================\n=== Pokemon Metadata ===\n========================\n"
	ret += fmt.Sprintf("Gender threshold: %d\n", pm.GenderThreshold)
	ret += fmt.Sprintf("Growth type: %d\n", pm.GrowthType)
	ret += fmt.Sprintf("Types: %s  |  %s\n", enums.PokemonType(pm.Type1), enums.PokemonType(pm.Type2))
	a1, _ := data.GetAbility(uint(pm.Ability1))
	a2, _ := data.GetAbility(uint(pm.Ability2))
	ret += fmt.Sprintf("Ability 1: '%s'\n", a1)
	ret += fmt.Sprintf("Ability 2: '%s'\n", a2)
	ret += fmt.Sprintf("------------- Base Stats -------------\n%s\n", pm.Base)
	ret += fmt.Sprintf("-------------- EV Yield -------------\n%s\n", pm.EVYield)
	
	return ret
}

func NewPokemon(buffer []byte, offset int) PokemonMetadata {
	obj := buffer[offset : offset+44]

	rawEV := utils.U16(obj, 10)
	getEV := func(index int) uint8 {
		return uint8((rawEV >> (index*2)) & 0b11)
	}

	return PokemonMetadata{
		Base: models.Stat[uint8]{
			Hp: utils.U8(obj, 0),
			Attack: utils.U8(obj, 1),
			Defense: utils.U8(obj, 2),
			SpeAttack: utils.U8(obj, 4),
			SpeDefense: utils.U8(obj,5),
			Speed: utils.U8(obj, 3),
		},
		Type1: utils.U8(obj, 6),
		Type2: utils.U8(obj, 7),
		CatchRate: utils.U8(obj, 8),
		ExpYield: utils.U8(obj, 9),
		EVYield: models.Stat[uint8]{
			Hp: getEV(0),
			Attack: getEV(1),
			Defense: getEV(2),
			SpeAttack: getEV(4),
			SpeDefense: getEV(5),
			Speed: getEV(3),
		},
		Item1: utils.U16(obj, 12),
		Item2: utils.U16(obj, 14),
		GenderThreshold: utils.U8(obj, 16),
		EggCycles: utils.U8(obj, 17),
		BaseFriendship: utils.U8(obj, 18),
		GrowthType: utils.U8(obj, 19),
		EggGroup1: utils.U8(obj, 20),
		EggGroup2: utils.U8(obj, 21),
		Ability1: utils.U8(obj, 22),
		Ability2: utils.U8(obj, 23),
		SafariZoneRate: utils.U8(obj, 24),
		Color: utils.U8(obj, 25),
		Padding1: utils.U16(obj, 26),
		MoveFlags: obj[28 : 41],
		Padding2: obj[41 : 44],
	}
}

func RipPokemonData() {
	path := path_resolver.GetRoot() + "/roms/pkmn-pt.nds"
	f, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x0370A400)
	yer := narcFile.FrameFIMG.Data.Data

	for i := range 508 {
		pkmn := NewPokemon(yer, i*44)
		fmt.Println(pkmn)
	}
}