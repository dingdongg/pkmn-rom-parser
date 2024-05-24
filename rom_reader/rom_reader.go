package rom_reader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/char_encoder"
	"github.com/dingdongg/pkmn-rom-parser/prng"
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
	HeldItemId uint16 // just return the in-memory value for now, figure out the mapping later
	Nature     string
	AbilityId  uint
	EVs        Stats
}

const (
	A uint = iota
	B uint = iota
	C uint = iota
	D uint = iota
)

const BLOCK_SIZE_BYTES uint = 32
const PARTY_POKEMON_SIZE uint = 236

var natureTable [25]string = [25]string{
	"Hardy",
	"Lonely",
	"Brave",
	"Adamant",
	"Naughty",
	"Bold",
	"Docile",
	"Relaxed",
	"Impish",
	"Lax",
	"Timid",
	"Hasty",
	"Serious",
	"Jolly",
	"Naive",
	"Modest",
	"Mild",
	"Quiet",
	"Bashful",
	"Rash",
	"Calm",
	"Gentle",
	"Sassy",
	"Careful",
	"Quirky",
}

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

// `ciphertext` must be a slice with the first byte
// referring to the first pokemon data structure
func GetPokemon(ciphertext []byte, partyIndex uint) Pokemon {
	offset := partyIndex * PARTY_POKEMON_SIZE

	personality := binary.LittleEndian.Uint32(ciphertext[offset : offset+4])
	checksum := binary.LittleEndian.Uint16(ciphertext[offset+6 : offset+8])

	rand := prng.Init(checksum, personality)
	return decryptPokemon(rand, ciphertext[offset:])
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

// first block of ciphertext points to offset 0x88 in a party pokemon block
// TODO: needs some validation/testing
func getPokemonBattleStats(ciphertext []byte, personality uint32) BattleStat {
	bsprng := prng.InitBattleStatPRNG(personality)
	var plaintext []byte

	for i := 0; i < 0x14; i += 2 {
		decrypted := bsprng.Next() ^ binary.LittleEndian.Uint16(ciphertext[i:i+2])
		plaintext = append(plaintext, byte(decrypted&0xFF), byte((decrypted>>8)&0xFF))
	}

	stats := Stats{
		uint(binary.LittleEndian.Uint16(plaintext[0x8:0xA])),
		uint(binary.LittleEndian.Uint16(plaintext[0xA:0xC])),
		uint(binary.LittleEndian.Uint16(plaintext[0xC:0xE])),
		uint(binary.LittleEndian.Uint16(plaintext[0x10:0x12])),
		uint(binary.LittleEndian.Uint16(plaintext[0x12:0x14])),
		uint(binary.LittleEndian.Uint16(plaintext[0xE:0x10])),
	}

	return BattleStat{uint(plaintext[4]), stats}
}

func decryptPokemon(prng prng.PRNG, ciphertext []byte) Pokemon {
	plaintext_buf := ciphertext[:8]
	plaintext_sum := uint16(0)

	// 1. XOR to get plaintext words
	for i := 0x8; i < 0x87; i += 0x2 {
		word := binary.LittleEndian.Uint16(ciphertext[i : i+2])
		plaintext := word ^ prng.Next()
		plaintext_sum += plaintext
		littleByte := byte(plaintext & 0x00FF)
		bigByte := byte((plaintext >> 8) & 0x00FF)
		plaintext_buf = append(plaintext_buf, littleByte, bigByte)
	}

	if plaintext_sum == prng.Checksum {
		fmt.Println("checksum is valid!")
	} else {
		fmt.Printf("Checksum invalid. expected 0x%x, got 0x%x\n", prng.Checksum, plaintext_sum)
	}

	blockA, err := getPokemonBlock(plaintext_buf, A, prng.Personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block A: ", err)
	}

	blockC, err := getPokemonBlock(plaintext_buf, C, prng.Personality)
	if err != nil {
		log.Fatal("Unexpected error while parsing block C: ", err)
	}

	dexId := binary.LittleEndian.Uint16(blockA[:2])
	heldItem := binary.LittleEndian.Uint16(blockA[2:4])
	nature := natureTable[prng.Personality%25]
	ability := binary.LittleEndian.Uint16(blockA[0xD:0xF])

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

	battleStats := getPokemonBattleStats(ciphertext[0x88:], prng.Personality)

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

	fmt.Printf("Total EV Spenditure: %d / 510\n", evSum)

	return Pokemon{
		dexId,
		name,
		battleStats,
		heldItem,
		nature,
		uint(ability),
		Stats{
			uint(blockA[hpEVOffset]),
			uint(blockA[attackEVOffset]),
			uint(blockA[defenseEVOffset]),
			uint(blockA[specialAtkEVOffset]),
			uint(blockA[specialDefEVOffset]),
			uint(blockA[speedEVOffset]),
		},
	}
}
