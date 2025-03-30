package savefile

import (
	"encoding/binary"
	"fmt"
	"strings"

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

func printBlock(buf []byte, w int) {
	output := ""
	buffer := make([]string, 0)

	update := func() {
		buffer = append(buffer, "\n")
		output += strings.Join(buffer, " ")
		buffer = make([]string, 0)
	}

	for i, b := range buf {
		buffer = append(buffer, fmt.Sprintf("%02X", b))
		if i % w == w - 1 {
			update()
		}
	}
	if len(buffer) != 0 {
		update()
	}
	fmt.Println(output)
}

func (pt *PlatSavefile) parsePokemon(index int) models.Pokemon {
	offset := 0xA0 + index*236
	raw := crypt.DecryptPokemon(pt.rawBytes[offset : offset+236])

	blocks := shuffler.GetPokemonBlocks(raw)
	a := blocks[0]
	printBlock(a, 16)

	battleStats := raw[0x88 : 0x88+0x64]

	ability, _ := data.GetAbility(uint(u8(a, 0xD)))
	item, _ := data.GetItem(u16(a, 0x2))

	return models.Pokemon{
		PokedexId: u16(a, 0x0),
		Exp: u32(a, 0x8),
		Level: u8(battleStats, 0x4),
		Nature: enums.Nature(u32(raw, 0) % 25),
		Ability: ability,
		HeldItem: item.Name,
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

func (pt *PlatSavefile) validate() error {
	// just validate the latest block pair for now
	/*
	checksum validation
	magic number in correct place 
	
	*/
	sfOffset := 0xCF18
	sf1, sf2 := sfOffset, sfOffset + 0x40000
	sf1Count, sf2Count := sf1 + 0x4, sf2 + 0x4

	// find latest 
	var latestOffset int
	countOne := binary.LittleEndian.Uint32(pt.rawBytes[sf1Count : sf1Count+0x4])
	countTwo := binary.LittleEndian.Uint32(pt.rawBytes[sf2Count : sf2Count+0x4])

	if countOne > countTwo {
		latestOffset = sf1
	} else if countOne < countTwo {
		latestOffset = sf2
	} else {
		return fmt.Errorf("SB save counts are same? how do I handle this?")
	}

	// check magic number
	magicNumOffset := latestOffset + 0xC
	magicNum := binary.LittleEndian.Uint32(pt.rawBytes[magicNumOffset : magicNumOffset+0x4])

	if magicNum != enums.MAGIC_TS_JP_INTL && magicNum != enums.MAGIC_TS_KR {
		return fmt.Errorf("magic numbers invalid")
	}

	// SB checksum vlaidation
	checksumOffset := latestOffset + 0x12
	expected := binary.LittleEndian.Uint16(pt.rawBytes[checksumOffset : checksumOffset+0x2])
	actual := crypt.CRC16_CCITT(pt.rawBytes[latestOffset-sfOffset : latestOffset])

	if expected != actual {
		return fmt.Errorf("[SMALL BLOCK] checksum mismatch: 0x%04X (expected), 0x%04X (actual)", expected, actual)
	}

	// big block validations
	bigBlockCount := binary.LittleEndian.Uint32(pt.rawBytes[latestOffset : latestOffset+0x4])

	bfOffset := 0x1F0FC
	bf1, bf2 := bfOffset, bfOffset + 0x40000
	var latestBigBlockAddr int
	bf1Count := binary.LittleEndian.Uint32(pt.rawBytes[bf1 : bf1+0x4])
	bf2Count := binary.LittleEndian.Uint32(pt.rawBytes[bf2 : bf2+0x4])

	if bigBlockCount == bf1Count {
		latestBigBlockAddr = bf1
	} else if bigBlockCount == bf2Count {
		latestBigBlockAddr = bf2
	} else {
		return fmt.Errorf("no big block match found")
	}

	bbSize := binary.LittleEndian.Uint32(pt.rawBytes[latestBigBlockAddr+0x8 : latestBigBlockAddr+0x8+0x4])
	bbChecksumOffset := latestBigBlockAddr + 0x12
	expected = binary.LittleEndian.Uint16(pt.rawBytes[bbChecksumOffset : bbChecksumOffset+0x2])
	actual = crypt.CRC16_CCITT(pt.rawBytes[uint32(latestBigBlockAddr)-bbSize+0x14 : latestBigBlockAddr])

	if expected != actual {
		return fmt.Errorf("[BIG BLOCK] checksum mismatch: 0x%04X (expected), 0x%04X (actual)", expected, actual)
	}

	return nil
}

func (pt *PlatSavefile) Version() enums.GameVersion {
	return enums.PLAT
}

func (pt *PlatSavefile) Flush() error {
	return nil
}
