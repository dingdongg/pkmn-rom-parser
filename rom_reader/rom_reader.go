package rom_reader

import (
	"encoding/binary"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v6/char"
	"github.com/dingdongg/pkmn-rom-parser/v6/consts"
	"github.com/dingdongg/pkmn-rom-parser/v6/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v6/data"
	"github.com/dingdongg/pkmn-rom-parser/v6/shuffler"
)

type Stats struct {
	Hp        uint
	Attack    uint
	Defense   uint
	SpAttack  uint
	SpDefense uint
	Speed     uint
}

type BattleStat struct {
	Level uint
	Stats Stats
}

type Pokemon struct {
	PokedexId uint16
	Name      string
	BattleStat
	Item    string
	Nature  string
	Ability string
	EVs     Stats
	IVs     Stats
}

const (
	A uint = iota
	B uint = iota
	C uint = iota
	D uint = iota
)

func GetPartyPokemon(ciphertext []byte) []Pokemon {
	var party []Pokemon

	for i := uint(0); i < 6; i++ {
		party = append(party, parsePokemon(ciphertext, i))
	}

	return party
}

func parsePokemon(ciphertext []byte, partyIndex uint) Pokemon {
	offset := partyIndex * consts.PARTY_POKEMON_SIZE
	plaintext := crypt.DecryptPokemon(ciphertext[offset:])
	personality := binary.LittleEndian.Uint32(plaintext[0:4])

	blockA, err := shuffler.GetPokemonBlock(plaintext, A, personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block A: ", err)
	}

	blockB, err := shuffler.GetPokemonBlock(plaintext, B, personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block B: ", err)
	}

	blockC, err := shuffler.GetPokemonBlock(plaintext, C, personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block C: ", err)
	}

	ivBytes := binary.LittleEndian.Uint32(blockB[0x10:0x14])

	ivs := Stats{
		uint((ivBytes >> 0) & 0b11111),
		uint((ivBytes >> 5) & 0b11111),
		uint((ivBytes >> 10) & 0b11111),
		uint((ivBytes >> 20) & 0b11111),
		uint((ivBytes >> 25) & 0b11111),
		uint((ivBytes >> 15) & 0b11111),
	}

	// fmt.Printf("% x\n", blockB[0x10:0x14])

	dexId := binary.LittleEndian.Uint16(blockA[:2])
	heldItem, err := data.GetItem(binary.LittleEndian.Uint16(blockA[2:4]))
	if err != nil {
		log.Fatal("error while parsing item: ", err)
	}

	nature, err := data.GetNature(uint(personality % 25))
	if err != nil {
		log.Fatal("error while parsing nature: ", err)
	}

	ability, err := data.GetAbility(uint(blockA[0xD]))
	if err != nil {
		log.Fatal("error while parsing ability: ", err)
	}

	pokemonNameLength := 22
	name := ""

	for i := 0; i < pokemonNameLength; i += 2 {
		charIndex := binary.LittleEndian.Uint16(blockC[i : i+2])
		str, err := char.Char(charIndex)
		if err != nil {
			break
		}
		name += str
	}

	battleStatsPlaintext := plaintext[0x88:]
	battleStats := BattleStat{
		uint(battleStatsPlaintext[4]),
		Stats{
			uint(binary.LittleEndian.Uint16(battleStatsPlaintext[0x8:0xA])),
			uint(binary.LittleEndian.Uint16(battleStatsPlaintext[0xA:0xC])),
			uint(binary.LittleEndian.Uint16(battleStatsPlaintext[0xC:0xE])),
			uint(binary.LittleEndian.Uint16(battleStatsPlaintext[0x10:0x12])),
			uint(binary.LittleEndian.Uint16(battleStatsPlaintext[0x12:0x14])),
			uint(binary.LittleEndian.Uint16(battleStatsPlaintext[0xE:0x10])),
		},
	}

	hpEVOffset := 0x10
	attackEVOffset := 0x11
	defenseEVOffset := 0x12
	speedEVOffset := 0x13
	specialAtkEVOffset := 0x14
	specialDefEVOffset := 0x15

	// fmt.Printf(
	// 	"Stats:\n\t- HP:  %d\n\t- ATK: %d\n\t- DEF: %d\n\t- SpA: %d\n\t- SpD: %d\n\t- SPE: %d\n",
	// 	blockA[hpEVOffset], blockA[attackEVOffset],
	// 	blockA[defenseEVOffset], blockA[specialAtkEVOffset],
	// 	blockA[specialDefEVOffset], blockA[speedEVOffset],
	// )

	evSum := 0
	for i := 0; i < 6; i++ {
		evSum += int(blockA[hpEVOffset+i])
	}

	return Pokemon{
		dexId,
		name,
		battleStats,
		heldItem.Name,
		nature,
		ability,
		Stats{
			uint(blockA[hpEVOffset]),
			uint(blockA[attackEVOffset]),
			uint(blockA[defenseEVOffset]),
			uint(blockA[specialAtkEVOffset]),
			uint(blockA[specialDefEVOffset]),
			uint(blockA[speedEVOffset]),
		},
		ivs,
	}
}
