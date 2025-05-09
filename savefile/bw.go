package savefile

import (
	"unicode/utf16"

	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/data"
	"github.com/dingdongg/pkmn-rom-parser/v7/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/ripper"
	"github.com/dingdongg/pkmn-rom-parser/v7/shuffler"
	"github.com/dingdongg/pkmn-rom-parser/v7/utils"
)

func NewBwSavefile(bytes []byte) *BwSavefile {
	return &BwSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]*models.Pokemon, 0),
		moveNames: ripper.RipMoveNamesGen5(),
		rawParty: make([]byte, 0),
		// TODO: implement ripper for gen 5 exp table
		// TODO: implement ripper for gen 5 pokemon metadata
	}
}

func (bw *BwSavefile) parsePokemon(index int) models.Pokemon {
	offset := 0x18E08 + index*220
	raw := crypt.DecryptPokemon(bw.rawBytes[offset : offset+220])

	blocks := shuffler.GetPokemonBlocks(raw)
	a, b, c := blocks[0], blocks[1], blocks[2]

	rawName := c[:0x16]
	yer := make([]uint16, 0)
	for i := 0; i < len(rawName); i += 2 {
		code := utils.U16(rawName, i)
		if code == 0xFFFF {
			break
		}
		yer = append(yer, code)
	}
	runes := utf16.Decode(yer)

	battleStats := raw[0x88 : 0x88+0x64]

	ability, _ := data.GetAbility(uint(utils.U8(a, 0xD)))
	item, _ := data.GetItem(utils.U16(a, 0x2))

	ivBuffer := utils.U32(b, 0x10)
	getIv := func(statIndex int) uint8 {
		val := (ivBuffer >> (5 * statIndex)) & 0x1F
		return uint8(val)
	}

	genderByte := utils.U8(b, 0x18)
	gender := enums.Male
	if genderByte&0b10 != 0 {
		gender = enums.Female
	} else if genderByte&0b100 != 0 {
		gender = enums.Unknown
	}

	moves := make([]models.Move, 0)
	for i := range 0x4 {
		id := utils.U16(b, i*0x2)
		move := models.Move{
			Id:   id,
			Name: bw.moveNames[id],
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
		Name:      string(runes),
		PokedexId: utils.U16(a, 0x0),
		Exp:       utils.U32(a, 0x8),
		Level:     utils.U8(battleStats, 0x4),
		Nature:    enums.Nature(utils.U8(b, 0x19)),
		Ability:   ability,
		HeldItem:  item.Name,
		Gender:    gender,
		Moves:     moves,
		EV: models.Stat[uint8]{
			Hp:         utils.U8(a, 0x10),
			Attack:     utils.U8(a, 0x11),
			Defense:    utils.U8(a, 0x12),
			SpeAttack:  utils.U8(a, 0x14),
			SpeDefense: utils.U8(a, 0x15),
			Speed:      utils.U8(a, 0x13),
		},
		IV: models.Stat[uint8]{
			Hp:         getIv(0),
			Attack:     getIv(1),
			Defense:    getIv(2),
			SpeAttack:  getIv(4),
			SpeDefense: getIv(5),
			Speed:      getIv(3),
		},
		Battle: models.Stat[uint16]{
			Hp:         utils.U16(battleStats, 0x8),
			Attack:     utils.U16(battleStats, 0xA),
			Defense:    utils.U16(battleStats, 0xC),
			SpeAttack:  utils.U16(battleStats, 0x10),
			SpeDefense: utils.U16(battleStats, 0x12),
			Speed:      utils.U16(battleStats, 0xE),
		},
	}
}

func (bw *BwSavefile) PartyPokemon() []*models.Pokemon {
	partySize := int(utils.U32(bw.rawBytes, 0x18E04))
	for i := range partySize {
		pkmn := bw.parsePokemon(i)
		bw.partyPokemon = append(bw.partyPokemon, &pkmn)
	}

	return bw.partyPokemon
}

func (bw *BwSavefile) validate() error {
	return nil
}

func (bw *BwSavefile) Version() enums.GameVersion {
	return enums.BW
}

func (bw *BwSavefile) Flush() error {
	return nil
}
