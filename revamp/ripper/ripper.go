package ripper

import (
	"fmt"
	"os"
	"unicode/utf16"

	"github.com/dingdongg/pkmn-rom-parser/v7/char"
	"github.com/dingdongg/pkmn-rom-parser/v7/data"
	"github.com/dingdongg/pkmn-rom-parser/v7/path_resolver"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/dsa"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/ripper/narc"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/walker"
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
	ret += fmt.Sprintf("Growth type:      %d\n", pm.GrowthType)
	t1, t2 := enums.PokemonType(pm.Type1), enums.PokemonType(pm.Type2)
	if t1 == t2 {
		ret += fmt.Sprintf("Type:\t\t  %s\n", t1)
	} else {
		ret += fmt.Sprintf("Types:\t\t  %s | %s\n", enums.PokemonType(pm.Type1), enums.PokemonType(pm.Type2))
	}
	a1, _ := data.GetAbility(uint(pm.Ability1))
	a2, _ := data.GetAbility(uint(pm.Ability2))
	ret += fmt.Sprintf("Ability 1:        '%s'\n", a1)
	ret += fmt.Sprintf("Ability 2:        '%s'\n\n", a2)
	ret += pm.Base.Print("Base Stats")
	ret += pm.EVYield.Print("EV Yield")
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

func buf2D[T any](x uint16, y uint16) [][]T {
	output := make([][]T, x)
	for i := range output {
		output[i] = make([]T, y)
	}

	return output
}

func RipMoveNamesGen5() []string {
	path := path_resolver.GetRoot() + "/roms/white.nds"
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x03471C00)
	moveFileMetadata := narcFile.FrameFATB.Data.Entry(203)

	size := moveFileMetadata.End - moveFileMetadata.Start + 1
	fmt.Printf("moves name offset: 0x%08X, size=0x%08X\n", moveFileMetadata.Start, size)

	buf := narcFile.FrameFIMG.Data.Data[moveFileMetadata.Start : moveFileMetadata.End]

	decryptFile := func(buffer []byte) []string {
		w := walker.NewWalker(buffer)
		numBlocks, numEntries := w.U16(), w.U16()
		// filesize, zero := w.U32(), w.U32()
		w.U32() // filesize, unused
		w.U32() // zero, unused

		blockOffsets := make([]uint32, numBlocks)
		tableOffsets := buf2D[uint32](numBlocks, numEntries)
		charCounts := buf2D[uint16](numBlocks, numEntries)
		textFlags := buf2D[uint16](numBlocks, numEntries)

		texts := make([][]string, numBlocks)
		for i := range texts {
			texts[i] = make([]string, numEntries) // technically this should be of length `numEntries`
		}

		for i := uint16(0); i < numBlocks; i++ {
			blockOffsets[i] = w.U32()
		}

		for i := uint16(0); i < numBlocks; i++ {
			w.Seek(int(blockOffsets[i]))

			// blockSize := w.U32()
			w.U32() // blockSize, unused
			for j := uint16(0); j < numEntries; j++ {
				tableOffsets[i][j] = w.U32()
				charCounts[i][j] = w.U16()
				textFlags[i][j] = w.U16()
			}

			for j := uint16(0); j < numEntries; j++ {
				encChars := dsa.NewSliceStack[uint16]()
				decChars := dsa.NewSliceStack[uint16]()
				// string := texts[i][j]

				w.Seek(int(blockOffsets[i]) + int(tableOffsets[i][j]))
				for k := uint16(0); k < charCounts[i][j]; k++ {
					encChars.Push(w.U16())
				}

				key := encChars.Peek()
				for !encChars.Empty() {
					val := ^(encChars.Pop() ^ key)
					decChars.Push(val)
					key = ((key >> 3) | (key << 13)) & 0xFFFF
				}

				charBuf := make([]uint16, 1)
				for !decChars.Empty() {
					charBuf[0] = decChars.Pop()
					char := charBuf[0]
					if char == 0xFFFF {
						break // continue, like gen 4?
					} else if char == 0xFFFE {
						texts[i][j] += "\n"
					} else if char == 0xF000 {
						fmt.Println("NEED TO APPEND SPECIAL CHAR")
						texts[i][j] += "😎"
					} else {
						res := string(utf16.Decode(charBuf))
						texts[i][j] += res
					}
				}
			}
		}

		// fmt.Println("output: ", texts)
		// return make([]string, 0) // stub
		return texts[0]
	}

	return decryptFile(buf)
}

// different for gen 5 ?
func RipMoveNames() []string {
	path := path_resolver.GetRoot() + "/roms/pkmn-pt.nds"
	f, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x0162DE00)
	moveFileMetadata := narcFile.FrameFATB.Data.Entry(647)

	// offsets are relative to start of FIMG buffer
	size := moveFileMetadata.End - moveFileMetadata.Start + 1
	fmt.Printf("moves name offset: 0x%08X, size=0x%08X\n", moveFileMetadata.Start, size)
	// yer := narcFile.FrameFIMG.Data.Data
	buf := narcFile.FrameFIMG.Data.Data[moveFileMetadata.Start : moveFileMetadata.End]
	// fmt.Println(buf)

	decryptFile := func(buffer []byte) []string {
		w := walker.NewWalker(buffer)
		num, seed := w.U16(), w.U16()

		offsets := make([]uint32, num)
		sizes := make([]uint32, num)

		// num * len(sizes) == num * num?
		// https://projectpokemon.org/rawdb/platinum/formats/msg.php
		binaryStrings := make([][]uint16, 0)
		for range num {
			binaryStrings = append(binaryStrings, make([]uint16, 0))
		}

		texts := make([]string, num)

		// generate offsets & sizes
		for i := uint16(1); i <= num; i++ {
			seedMult := seed * i
			key := uint32(((seedMult*0x02FD) & 0xFFFF)) | ((uint32(seedMult)*0x02FD0000) & 0xFFFF0000)
			offsets[i - 1] = w.U32() ^ key
			sizes[i - 1] = w.U32() ^ key
		}

		for i := uint16(1); i <= num; i++ {
			off := &offsets[i - 1]
			sz := &sizes[i - 1]
			bString := binaryStrings[i - 1]
			key := (uint32(0x91BD3)*uint32(i)) & 0x0000FFFF
			txt := &texts[i - 1]

			w.Seek(int(*off))

			for j := uint32(1); j <= *sz; j++ {
				bString = append(bString, w.U16() ^ uint16(key))
				key = (key+0x493D) & 0xFFFF
			}

			if bString[0] == 0xF100 {
				fmt.Println("de-compressing")
				// decompress from 9-bit strings to 16-bits
				newString := make([]uint16, 1)
				newString[0] = 0x0000
				bString = bString[:len(bString)-1] // pop()
				container, bit := uint16(0), uint16(0)

				for len(bString) != 0 {
					lastChar := bString[len(bString)-1]
					bString = bString[:len(bString)-1]
					container |= lastChar << bit

					for bit >= 9 {
						bit -= 9
						newString = append(newString, container & 0x01FF)
						container >>= 9
					}
				}
				binaryStrings[i - 1] = newString
				*sz = uint32(len(newString))
			}

			*txt = ""
			textStack := make([]string, 0)
			// TODO: instead of iterating from the back, we 
			// could just iterate normally...
			for len(bString) != 0 {
				lastChar := bString[len(bString)-1]
				bString = bString[:len(bString)-1] // pop()
				
				if lastChar == 0xFFFF {
					// break <-- will discard every string we look at
					continue
				} else if lastChar == 0xFFFE {
					c := bString[len(bString)-1]
					bString = bString[:len(bString)-1]
					args := []uint16{ 0x0000 }
					for k := uint16(1); k <= c; k++ {
						args = append(args, bString[len(bString)-1])
						bString = bString[:len(bString)-1]
					}

					for _, a := range args {
						converted, err := char.Char(a)
						if err != nil {
							fmt.Println("unrecognized character")
							os.Exit(1)
						}
						textStack = append(textStack, converted)
					}
				} else {
					// fmt.Printf("lastChar: 0x%04X\n", lastChar)
					c, err := char.Char(lastChar)
					if err != nil {
						fmt.Printf("unrecognized character!!!!! 0x%04X\n", lastChar)
					}
					textStack = append(textStack, c)
				}
			}

			// reverse stack and push into txt buffer
			// can be removed once iteration direction
			// of above loop is reversed
			k := len(textStack) - 1
			for k > -1 {
				*txt += textStack[k]
				k -= 1
			}
		}

		// fmt.Println("output:\n", texts)
		return texts
	}

	return decryptFile(buf)
}