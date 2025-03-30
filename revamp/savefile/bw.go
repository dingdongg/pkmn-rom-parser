package savefile

import (
	"encoding/binary"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v7/char"
	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/data"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/shuffler"
)


func NewBwSavefile(bytes []byte) *BwSavefile {
	return &BwSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]*models.Pokemon, 0),
	}
}

func (bw *BwSavefile) parsePokemon(index int) models.Pokemon {
	offset := 0x18E08 + index*220
	raw := crypt.DecryptPokemon(bw.rawBytes[offset : offset+220])

	blocks := shuffler.GetPokemonBlocks(raw)
	a, b, c := blocks[0], blocks[1], blocks[2]

	rawName := c[:0x16]
	name := ""
	for i := 0; i < len(rawName); i += 2 {
		code := u16(rawName, i)
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

	ability, _ := data.GetAbility(uint(u8(a, 0xD)))
	item, _ := data.GetItem(u16(a, 0x2))

	ivBuffer := u32(b, 0x10)
	getIv := func(statIndex int) uint8 {
		val := (ivBuffer >> (5*statIndex)) & 0x1F
		return uint8(val)
	}

	genderByte := u8(b, 0x18)
	gender := enums.Male
	if genderByte & 0b10 != 0 {
		gender = enums.Female
	} else if genderByte & 0b100 != 0 {
		gender = enums.Unknown
	}

	/*
	missing: 
	- base stats
	- alternate forms
	- movesets
	*/
	return models.Pokemon{
		Name: name,
		PokedexId: u16(a, 0x0),
		Exp: u32(a, 0x8),
		Level: u8(battleStats, 0x4),
		Nature: enums.Nature(u8(b, 0x19)),
		Ability: ability,
		HeldItem: item.Name,
		Gender: gender,
		EV: models.Stat[uint8]{
			Hp: u8(a, 0x10),
			Attack: u8(a, 0x11),
			Defense: u8(a, 0x12),
			SpeAttack: u8(a, 0x14),
			SpeDefense: u8(a, 0x15),
			Speed: u8(a, 0x13),
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
			Hp: u16(battleStats, 0x8),
			Attack: u16(battleStats, 0xA),
			Defense: u16(battleStats, 0xC),
			SpeAttack: u16(battleStats, 0x10),
			SpeDefense: u16(battleStats, 0x12),
			Speed: u16(battleStats, 0xE),
		},
	}
}

func (bw *BwSavefile) PartyPokemon() []*models.Pokemon {
	partySize := int(binary.LittleEndian.Uint32(bw.rawBytes[0x18E04:0x18E08]))
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
