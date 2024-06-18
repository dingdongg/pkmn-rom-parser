package shuffler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts"
)

// from https://projectpokemon.org/home/docs/gen-4/pkm-structure-r65/
const tableValues = `00	ABCD	ABCD
01	ABDC	ABDC
02	ACBD	ACBD
03	ACDB	ADBC
04	ADBC	ACDB
05	ADCB	ADCB
06	BACD	BACD
07	BADC	BADC
08	BCAD	CABD
09	BCDA	DABC
10	BDAC	CADB
11	BDCA	DACB
12	CABD	BCAD
13	CADB	BDAC
14	CBAD	CBAD
15	CBDA	DBAC
16	CDAB	CDAB
17	CDBA	DCAB
18	DABC	BCDA
19	DACB	BDCA
20	DBAC	CBDA
21	DBCA	DBCA
22	DCAB	CDBA
23	DCBA	DCBA`

const (
	A = iota
	B
	C
	D
)

type blockOrder struct {
	ShuffledPos [4]uint
	OriginalPos [4]uint
}

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

// Unless you need the offset address to the block, you want to use this function (NOT GetPokemonBlockLocation())
func GetPokemonBlock(buf []byte, block uint, personality uint32) ([]byte, error) {
	if block >= A && block <= D {
		shiftValue := ((personality & 0x03E000) >> 0x0D) % 24
		unshuffleInfo := unshuffleTable[shiftValue]
		startAddr := unshuffleInfo.GetUnshuffledPos(block)
		blockChunk := buf[startAddr : startAddr+consts.BLOCK_SIZE_BYTES]

		return blockChunk, nil
	}

	return make([]byte, 0), errors.New("invalid block index")
}

// Used to get the absolute memory address location of the block. Mainly for writing purposes ATM
func GetPokemonBlockLocation(block uint, personality uint32) (uint, error) {
	if block >= A && block <= D {
		shiftValue := ((personality & 0x03E000) >> 0x0D) % 24
		unshuffleInfo := unshuffleTable[shiftValue]
		startAddr := unshuffleInfo.GetUnshuffledPos(block)

		return startAddr, nil
	}
	return 0, errors.New("invalid block index")
}

/*
block is one of 0, 1, 2, 3

metadata consists of a pokemon's PID & checksum
*/
func (bo blockOrder) GetUnshuffledPos(block uint) uint {
	metadataOffset := uint(0x8)
	startIndex := bo.OriginalPos[block]
	res := metadataOffset + (startIndex * 32)
	return res
}

func Get(shiftValue uint32) blockOrder {
	return unshuffleTable[shiftValue]
}

// prints out a list of shuffle mapping
func Extract() {
	rows := strings.FieldsFunc(tableValues, func(r rune) bool { return r == '\n' })
	var sequences [24]([2]string)

	for i, r := range rows {
		tokens := strings.FieldsFunc(r, func(r rune) bool { return r == '\t' })
		sequences[i] = [2]string{tokens[1], tokens[2]}

		fmt.Println(getShuffleInfo(sequences[i]))
	}
}

func getShuffleInfo(sequence [2]string) string {
	return fmt.Sprintf(
		"{ [4]uint{%c, %c, %c, %c}, [4]uint{%c, %c, %c, %c} }, // %s %s",
		sequence[0][0], sequence[0][1], sequence[0][2], sequence[0][3],
		sequence[1][0], sequence[1][1], sequence[1][2], sequence[1][3],
		sequence[0], sequence[1],
	)
}
