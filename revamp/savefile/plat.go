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

func NewPlatSavefile(bytes []byte) *PlatSavefile {
	return &PlatSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]*models.Pokemon, 0),
	}
}

func u8(buf []byte, index int) uint8 {
	return uint8(buf[index])
}

func u16(buf []byte, index int) uint16 {
	return binary.LittleEndian.Uint16(buf[index : index+2])
}

func u32(buf []byte, index int) uint32 {
	return binary.LittleEndian.Uint32(buf[index : index+4])
}

func u64(buf []byte, index int) uint64 {
	return binary.LittleEndian.Uint64(buf[index : index+8])
}

func (pt *PlatSavefile) parsePokemon(index int) models.Pokemon {
	offset := 0xA0 + index*236
	raw := crypt.DecryptPokemon(pt.rawBytes[offset : offset+236])

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
	if genderByte & 0x2 != 0 {
		gender = enums.Female
	} else if genderByte & 0x4 != 0 {
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
		Nature: enums.Nature(u32(raw, 0) % 25),
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
