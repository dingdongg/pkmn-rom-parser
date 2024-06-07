package rom_reader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v3/char_encoder"
	"github.com/dingdongg/pkmn-rom-parser/v3/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v3/data"
)

type blockOrder struct {
	ShuffledPos [4]uint
	OriginalPos [4]uint
}

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
	IVs		Stats
}

const (
	A uint = iota
	B uint = iota
	C uint = iota
	D uint = iota
)

const BLOCK_SIZE_BYTES uint = 32
const PARTY_POKEMON_SIZE uint = 236

// populated with results from the shuffler package!
var unshuffleTable [24]blockOrder = [24]blockOrder{
	{[4]uint{A, B, C, D}, [4]uint{A, B, C, D}}, // ABCD ABCD
	{[4]uint{A, B, D, C}, [4]uint{A, B, D, C}}, // ABDC ABDC
	{[4]uint{A, C, B, D}, [4]uint{A, C, B, D}}, // ACBD ACBD
	{[4]uint{A, C, D, B}, [4]uint{A, D, B, C}}, // ACDB ADBC
	{[4]uint{A, D, B, C}, [4]uint{A, C, D, B}}, // ADBC ACDB
	{[4]uint{A, D, C, B}, [4]uint{A, D, C, B}}, // ADCB ADCB
	{[4]uint{B, A, C, D}, [4]uint{B, A, C, D}}, // BACD BACD
	{[4]uint{B, A, D, C}, [4]uint{B, A, D, C}}, // BADC BADC
	{[4]uint{B, C, A, D}, [4]uint{C, A, B, D}}, // BCAD CABD
	{[4]uint{B, C, D, A}, [4]uint{D, A, B, C}}, // BCDA DABC
	{[4]uint{B, D, A, C}, [4]uint{C, A, D, B}}, // BDAC CADB
	{[4]uint{B, D, C, A}, [4]uint{D, A, C, B}}, // BDCA DACB
	{[4]uint{C, A, B, D}, [4]uint{B, C, A, D}}, // CABD BCAD
	{[4]uint{C, A, D, B}, [4]uint{B, D, A, C}}, // CADB BDAC
	{[4]uint{C, B, A, D}, [4]uint{C, B, A, D}}, // CBAD CBAD
	{[4]uint{C, B, D, A}, [4]uint{D, B, A, C}}, // CBDA DBAC
	{[4]uint{C, D, A, B}, [4]uint{C, D, A, B}}, // CDAB CDAB
	{[4]uint{C, D, B, A}, [4]uint{D, C, A, B}}, // CDBA DCAB
	{[4]uint{D, A, B, C}, [4]uint{B, C, D, A}}, // DABC BCDA
	{[4]uint{D, A, C, B}, [4]uint{B, D, C, A}}, // DACB BDCA
	{[4]uint{D, B, A, C}, [4]uint{C, B, D, A}}, // DBAC CBDA
	{[4]uint{D, B, C, A}, [4]uint{D, B, C, A}}, // DBCA DBCA
	{[4]uint{D, C, A, B}, [4]uint{C, D, B, A}}, // DCAB CDBA
	{[4]uint{D, C, B, A}, [4]uint{D, C, B, A}}, // DCBA DCBA
}

func GetPartyPokemon(ciphertext []byte) []Pokemon {
	var party []Pokemon

	for i := uint(0); i < 6; i++ {
		party = append(party, parsePokemon(ciphertext, i))
	}

	return party
}

// block is one of 0, 1, 2, 3
func (bo blockOrder) getUnshuffledPos(block uint) uint {
	metadataOffset := uint(0x8)
	startIndex := bo.OriginalPos[block]
	res := metadataOffset + (startIndex * BLOCK_SIZE_BYTES)
	return res
}

func getPokemonBlock(buf []byte, block uint, personality uint32) ([]byte, error) {
	if block >= A && block <= D {
		shiftValue := ((personality & 0x03E000) >> 0x0D) % 24
		unshuffleInfo := unshuffleTable[shiftValue]
		startAddr := unshuffleInfo.getUnshuffledPos(block)
		blockChunk := buf[startAddr : startAddr+BLOCK_SIZE_BYTES]

		return blockChunk, nil
	}

	return make([]byte, 0), errors.New("invalid block index")
}

func parsePokemon(ciphertext []byte, partyIndex uint) Pokemon {
	offset := partyIndex * PARTY_POKEMON_SIZE
	plaintext := crypt.DecryptPokemon(ciphertext[offset:])
	personality := binary.LittleEndian.Uint32(plaintext[0 : 4])

	blockA, err := getPokemonBlock(plaintext, A, personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block A: ", err)
	}

	blockB, err := getPokemonBlock(plaintext, B, personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block B: ", err)
	}

	blockC, err := getPokemonBlock(plaintext, C, personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block C: ", err)
	}

	ivBytes := binary.LittleEndian.Uint32(blockB[0x10 : 0x14])

	ivs := Stats {
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
		log.Fatal(err)
	}

	nature, err := data.GetNature(uint(personality % 25))
	if err != nil {
		log.Fatal(err)
	}

	ability, err := data.GetAbility(uint(blockA[0xD]))
	if err != nil {
		log.Fatal(err)
	}

	pokemonNameLength := 22
	name := ""

	for i := 0; i < pokemonNameLength; i += 2 {
		charIndex := binary.LittleEndian.Uint16(blockC[i : i+2])
		str, err := char_encoder.Char(charIndex)
		if err != nil {
			break
		}
		name += str
	}

	fmt.Printf("Pokemon: '%s'\n", name)

	battleStatsPlaintext := crypt.DecryptBattleStats(ciphertext[offset + 0x88:], personality)
	battleStats := BattleStat{
		uint(battleStatsPlaintext[4]),
		Stats{
			uint(binary.LittleEndian.Uint16(plaintext[0x8:0xA])),
			uint(binary.LittleEndian.Uint16(plaintext[0xA:0xC])),
			uint(binary.LittleEndian.Uint16(plaintext[0xC:0xE])),
			uint(binary.LittleEndian.Uint16(plaintext[0x10:0x12])),
			uint(binary.LittleEndian.Uint16(plaintext[0x12:0x14])),
			uint(binary.LittleEndian.Uint16(plaintext[0xE:0x10])),
		},
	}

	hpEVOffset := 0x10
	attackEVOffset := 0x11
	defenseEVOffset := 0x12
	speedEVOffset := 0x13
	specialAtkEVOffset := 0x14
	specialDefEVOffset := 0x15

	fmt.Printf(
		"Stats:\n\t- HP:  %d\n\t- ATK: %d\n\t- DEF: %d\n\t- SpA: %d\n\t- SpD: %d\n\t- SPE: %d\n",
		blockA[hpEVOffset], blockA[attackEVOffset],
		blockA[defenseEVOffset], blockA[specialAtkEVOffset],
		blockA[specialDefEVOffset], blockA[speedEVOffset],
	)

	evSum := 0
	for i := 0; i < 6; i++ {
		evSum += int(blockA[hpEVOffset+i])
	}

	// fmt.Printf("Total EV Spenditure: %d / 510\n", evSum)

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
