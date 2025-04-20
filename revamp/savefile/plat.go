package savefile

import (
	"encoding/binary"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v7/char"
	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/data"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/ripper"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
	"github.com/dingdongg/pkmn-rom-parser/v7/shuffler"
)

func NewPlatSavefile(bytes []byte) *PlatSavefile {
	return &PlatSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]*models.Pokemon, 0),
	}
}

func (pt *PlatSavefile) parsePokemon(index int) models.Pokemon {
	offset := 0xA0 + index*236
	raw := crypt.DecryptPokemon(pt.rawBytes[offset : offset+236])

	blocks := shuffler.GetPokemonBlocks(raw)
	a, b, c := blocks[0], blocks[1], blocks[2]

	rawName := c[:0x16]
	name := ""
	for i := 0; i < len(rawName); i += 2 {
		code := utils.U16(rawName, i)
		if code == char.END_OF_STRING {
			break
		}
		chr, err := char.Char(code)
		if err != nil {
			log.Fatal(err)
		}
		name += chr
	}

	battleStats := raw[0x88 : 0x88+0x64]

	ability, _ := data.GetAbility(uint(utils.U8(a, 0xD)))
	item, _ := data.GetItem(utils.U16(a, 0x2))

	ivBuffer := utils.U32(b, 0x10)
	getIv := func(statIndex int) uint8 {
		val := (ivBuffer >> (5*statIndex)) & 0x1F
		return uint8(val)
	}

	genderByte := utils.U8(b, 0x18)
	gender := enums.Male
	if genderByte & 0x2 != 0 {
		gender = enums.Female
	} else if genderByte & 0x4 != 0 {
		gender = enums.Unknown
	}

	moves := make([]models.Move, 0)
	moveNames := ripper.RipMoveNames()
	for i := range 0x4 {
		id := utils.U16(b, i*0x2)
		move := models.Move{
			Id: id,
			Name: moveNames[id],
		}
		moves = append(moves, move)
	}

	/*
	missing: 
	- base stats
	- alternate forms
	- movesets
	*/
	return models.Pokemon{
		Name: name,
		PokedexId: utils.U16(a, 0x0),
		Exp: utils.U32(a, 0x8),
		Level: utils.U8(battleStats, 0x4),
		Nature: enums.Nature(utils.U32(raw, 0) % 25),
		Ability: ability,
		HeldItem: item.Name,
		Gender: gender,
		Moves: moves,
		EV: models.Stat[uint8]{
			Hp: utils.U8(a, 0x10),
			Attack: utils.U8(a, 0x11),
			Defense: utils.U8(a, 0x12),
			SpeAttack: utils.U8(a, 0x14),
			SpeDefense: utils.U8(a, 0x15),
			Speed: utils.U8(a, 0x13),
		},
		IV: models.Stat[uint8]{
			Hp: getIv(0),
			Attack: getIv(1),
			Defense: getIv(2),
			SpeAttack: getIv(4),
			SpeDefense: getIv(5),
			Speed: getIv(3),
		},
		Battle: models.Stat[uint16]{
			Hp: utils.U16(battleStats, 0x8),
			Attack: utils.U16(battleStats, 0xA),
			Defense: utils.U16(battleStats, 0xC),
			SpeAttack: utils.U16(battleStats, 0x10),
			SpeDefense: utils.U16(battleStats, 0x12),
			Speed: utils.U16(battleStats, 0xE),
		},
	}
}

func (pt *PlatSavefile) PartyPokemon() []*models.Pokemon {
	partySize := int(binary.LittleEndian.Uint32(pt.rawBytes[0x9C:0xA0]))
	for i := range partySize {
		pkmn := pt.parsePokemon(i)
		pt.partyPokemon = append(pt.partyPokemon, &pkmn)	
	}

	return pt.partyPokemon
}

/*
this functino shsould be aimed at validating the internal pokemon data before flushing. 

thus, validation should be an external function not tied to any concrete impl. of Savefile 

TODO: complete this function, and write the move parsing logic for all concrete savefiles
*/
func (pt *PlatSavefile) validate() error {
	return nil
}

func (pt *PlatSavefile) Version() enums.GameVersion {
	return enums.PLAT
}

func (pt *PlatSavefile) Flush() error {
	if err := pt.validate(); err != nil {
		return err
	}

	// flush that shit
	return nil
}
