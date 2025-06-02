package ripper

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf16"

	"github.com/dingdongg/pkmn-rom-parser/v7/char"
	"github.com/dingdongg/pkmn-rom-parser/v7/data"
	"github.com/dingdongg/pkmn-rom-parser/v7/dsa"
	"github.com/dingdongg/pkmn-rom-parser/v7/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/path_resolver"
	"github.com/dingdongg/pkmn-rom-parser/v7/ripper/narc"
	"github.com/dingdongg/pkmn-rom-parser/v7/utils"
	"github.com/dingdongg/pkmn-rom-parser/v7/walker"
)

type PokemonMetadata struct {
	Base            models.Stat[uint8]
	Type1           uint8
	Type2           uint8
	CatchRate       uint8
	ExpYield        uint8
	EVYield         models.Stat[uint8]
	Item1           uint16
	Item2           uint16
	GenderThreshold uint8
	EggCycles       uint8
	BaseFriendship  uint8
	GrowthType      uint8
	EggGroup1       uint8
	EggGroup2       uint8
	Ability1        uint8
	Ability2        uint8
	SafariZoneRate  uint8
	Color           uint8
	Padding1        uint16
	MoveFlags       []byte // length 13
	Padding2        []byte // length 3
}

type PokemonMetadataGen5 struct {
	Base            models.Stat[uint8]
	Type1           uint8
	Type2           uint8
	CatchRate       uint8
	Stage           uint8
	EVYield         models.Stat[uint8]
	Item1           uint16
	Item2           uint16
	Item3           uint16
	GenderThreshold uint8
	EggCycles       uint8
	BaseFriendship  uint8
	GrowthType      uint8
	EggGroup1       uint8
	EggGroup2       uint8
	Ability1        uint8
	Ability2        uint8
	Ability3        uint8
	Flee            uint8 // SafariZoneRate?
	FormId          uint16
	Form            uint16
	NumForms        uint8
	Color           uint8
	BaseExp         uint16
	Height          uint16
	Weight          uint16
}

type ExperienceTable = []uint32

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

func NewPokemonGen4(buffer []byte, offset int) PokemonMetadata {
	obj := buffer[offset : offset+44]

	rawEV := utils.U16(obj, 10)
	getEV := func(index int) uint8 {
		return uint8((rawEV >> (index * 2)) & 0b11)
	}

	return PokemonMetadata{
		Base: models.Stat[uint8]{
			Hp:         utils.U8(obj, 0),
			Attack:     utils.U8(obj, 1),
			Defense:    utils.U8(obj, 2),
			SpeAttack:  utils.U8(obj, 4),
			SpeDefense: utils.U8(obj, 5),
			Speed:      utils.U8(obj, 3),
		},
		Type1:     utils.U8(obj, 6),
		Type2:     utils.U8(obj, 7),
		CatchRate: utils.U8(obj, 8),
		ExpYield:  utils.U8(obj, 9),
		EVYield: models.Stat[uint8]{
			Hp:         getEV(0),
			Attack:     getEV(1),
			Defense:    getEV(2),
			SpeAttack:  getEV(4),
			SpeDefense: getEV(5),
			Speed:      getEV(3),
		},
		Item1:           utils.U16(obj, 12),
		Item2:           utils.U16(obj, 14),
		GenderThreshold: utils.U8(obj, 16),
		EggCycles:       utils.U8(obj, 17),
		BaseFriendship:  utils.U8(obj, 18),
		GrowthType:      utils.U8(obj, 19),
		EggGroup1:       utils.U8(obj, 20),
		EggGroup2:       utils.U8(obj, 21),
		Ability1:        utils.U8(obj, 22),
		Ability2:        utils.U8(obj, 23),
		SafariZoneRate:  utils.U8(obj, 24),
		Color:           utils.U8(obj, 25),
		Padding1:        utils.U16(obj, 26),
		MoveFlags:       obj[28:41],
		Padding2:        obj[41:44],
	}
}

func NewPokemonGen5(data []byte) PokemonMetadataGen5 {
	w := walker.NewWalker(data)
	w.Seek(6) // skip base stats

	pokemon := PokemonMetadataGen5{
		Base: models.Stat[uint8]{
			Hp:         utils.U8(data, 0),
			Attack:     utils.U8(data, 1),
			Defense:    utils.U8(data, 2),
			SpeAttack:  utils.U8(data, 4),
			SpeDefense: utils.U8(data, 5),
			Speed:      utils.U8(data, 3),
		},
		Type1:     w.U8(),
		Type2:     w.U8(),
		CatchRate: w.U8(),
		Stage:     w.U8(),
	}

	rawEV := w.U16()
	getEV := func(index int) uint8 {
		return uint8((rawEV >> (index * 2)) & 0b11)
	}
	pokemon.EVYield = models.Stat[uint8]{
		Hp:         getEV(0),
		Attack:     getEV(1),
		Defense:    getEV(2),
		SpeAttack:  getEV(4),
		SpeDefense: getEV(5),
		Speed:      getEV(3),
	}

	pokemon.Item1 = w.U16()
	pokemon.Item2 = w.U16()
	pokemon.Item3 = w.U16()
	pokemon.GenderThreshold = w.U8()
	pokemon.EggCycles = w.U8()
	pokemon.BaseFriendship = w.U8()
	pokemon.GrowthType = w.U8()
	pokemon.EggGroup1 = w.U8()
	pokemon.EggGroup2 = w.U8()
	pokemon.Ability1 = w.U8()
	pokemon.Ability2 = w.U8()
	pokemon.Ability3 = w.U8()
	pokemon.Flee = w.U8()
	pokemon.FormId = w.U16()
	pokemon.Form = w.U16()
	pokemon.NumForms = w.U8()
	pokemon.Color = w.U8()
	pokemon.BaseExp = w.U16()
	pokemon.Height = w.U16()
	pokemon.Weight = w.U16()

	return pokemon
}

func RipPokemonDataGen4() []PokemonMetadata {
	path := path_resolver.GetRoot() + "/roms/pkmn-pt.nds"
	f, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x0370A400)
	buffer := narcFile.FrameFIMG.Data.Data
	ret := make([]PokemonMetadata, 0)

	numEntries := len(narcFile.FrameFATB.Data.Entries)
	for i := range numEntries {
		ret = append(ret, NewPokemonGen4(buffer, i*44))
	}

	return ret
}

func RipPokemonDataGen5() []PokemonMetadataGen5 {
	path := path_resolver.GetRoot() + "/roms/white.nds"
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x0694D400)
	data := narcFile.FrameFIMG.Data.Data
	ret := make([]PokemonMetadataGen5, 0)
	entries := narcFile.FrameFATB.Data.Entries

	for _, e := range entries {
		pokemon := NewPokemonGen5(data[e.Start:e.End])
		ret = append(ret, pokemon)
	}

	return ret
}

func buf2D[T any](x uint16, y uint16) [][]T {
	output := make([][]T, x)
	for i := range output {
		output[i] = make([]T, y)
	}

	return output
}

func RipExpTableGen5() []ExperienceTable {
	path := path_resolver.GetRoot() + "/roms/white.nds"
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x06958C00)
	entries := narcFile.FrameFATB.Data.Entries
	ret := make([]ExperienceTable, 0)
	data := narcFile.FrameFIMG.Data.Data

	for _, e := range entries {
		table := newGrowthTableGen4(data[e.Start:e.End]) // same format as gen 4 games
		ret = append(ret, table)
	}

	return ret
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

	buf := narcFile.FrameFIMG.Data.Data[moveFileMetadata.Start:moveFileMetadata.End]

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
					val := ^(encChars.Pop() ^ key) // have to negate the resulting value for some reason
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

// TODO: rip from a B2W2 file instead, since this one is 
// just a subset of the B2W2 item names
func RipItemNamesGen5() []string {
	path := path_resolver.GetRoot() + "/roms/white.nds"
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x03471C00)
	itemNameFileRange := narcFile.FrameFATB.Data.Entry(54)

	buf := narcFile.FrameFIMG.Data.Data[itemNameFileRange.Start:itemNameFileRange.End]

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
					val := ^(encChars.Pop() ^ key) // have to negate the resulting value for some reason
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

func newGrowthTableGen4(file []byte) ExperienceTable {
	// 101 entries of uint32s
	table := make(ExperienceTable, 0)

	for i := 0x0; i < len(file); i += 4 {
		table = append(table, utils.U32(file, i))
	}

	return table
}

func RipExpTableGen4() []ExperienceTable {
	path := path_resolver.GetRoot() + "/roms/pkmn-pt.nds"
	f, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x03718200)
	entries := narcFile.FrameFATB.Data.Entries
	ret := make([]ExperienceTable, 0)
	data := narcFile.FrameFIMG.Data.Data

	for _, e := range entries {
		tbl := newGrowthTableGen4(data[e.Start:e.End])
		ret = append(ret, tbl)
	}

	return ret
}

func parseMessageFileGen4(buffer []byte) []string {
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
		key := uint32(((seedMult * 0x02FD) & 0xFFFF)) | ((uint32(seedMult) * 0x02FD0000) & 0xFFFF0000)
		offsets[i-1] = w.U32() ^ key
		sizes[i-1] = w.U32() ^ key
	}

	for i := uint16(1); i <= num; i++ {
		off := &offsets[i-1]
		sz := &sizes[i-1]
		bString := binaryStrings[i-1]
		key := (uint32(0x91BD3) * uint32(i)) & 0x0000FFFF
		txt := &texts[i-1]

		w.Seek(int(*off))

		for j := uint32(1); j <= *sz; j++ {
			bString = append(bString, w.U16()^uint16(key))
			key = (key + 0x493D) & 0xFFFF
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
					newString = append(newString, container&0x01FF)
					container >>= 9
				}
			}
			binaryStrings[i-1] = newString
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
				args := []uint16{0x0000}
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
	buf := narcFile.FrameFIMG.Data.Data[moveFileMetadata.Start:moveFileMetadata.End]
	// fmt.Println(buf)

	return parseMessageFileGen4(buf)
}

func RipPokemonNamesGen4() []string {
	path := path_resolver.GetRoot() + "/roms/pkmn-pt.nds"
	f, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	narcFile := narc.NewNarcFile(f, 0x0162DE00)
	namesFileMetadata := narcFile.FrameFATB.Data.Entry(712)

	buf := narcFile.FrameFIMG.Data.Data[namesFileMetadata.Start:namesFileMetadata.End]
	names := parseMessageFileGen4(buf)
	for i := range names {
		names[i] = names[i][5:]
	}
	return names
}

type TrainerPokemon struct {
	Difficulty uint16
	Level      uint16 // technicall u8 + u8 of padding, but same thing since this is little-endian
	Species    uint16
	Seal       uint16
	ItemId     *uint16   // optional
	MoveIds    []*uint16 // optional
}

var pokemonNames []string = RipPokemonNamesGen4()
var moveNames []string = RipMoveNames()

func (pkmn TrainerPokemon) String() string {
	tokens := make([]string, 0)
	tokens = append(tokens, fmt.Sprintf("  difficulty: %d", pkmn.Difficulty))
	tokens = append(tokens, fmt.Sprintf("  level:      %d", pkmn.Level))

	form, pokedexId := (pkmn.Species&0x0C00)>>10, pkmn.Species&0x3FF
	pokemonName := pokemonNames[pokedexId]
	tokens = append(tokens, "  Species:")
	tokens = append(tokens, fmt.Sprintf("    form:       %d", form))
	tokens = append(tokens, fmt.Sprintf("    pokedex id: %d (%s)", pokedexId, pokemonName))
	tokens = append(tokens, fmt.Sprintf("  seal:       %d", pkmn.Seal))

	if pkmn.ItemId != nil {
		tokens = append(tokens, fmt.Sprintf("  item ID:    %d", *pkmn.ItemId))
	}
	for i, m := range pkmn.MoveIds {
		if m != nil {
			tokens = append(tokens, fmt.Sprintf("  move ID #%d: %d (%s)", i+1, *m, moveNames[*m]))
		}
	}

	return "\n" + strings.Join(tokens, "\n") + "\n"
}

type PlatTrainer struct {
	Flags        uint8
	Class        uint8
	BattleType   uint8
	NumPokemon   uint8
	Item1        uint16
	Item2        uint16
	Item3        uint16
	Item4        uint16
	AiType       uint32
	BattleType2  uint32
	PartyPokemon []TrainerPokemon
}

func (tr PlatTrainer) String() string {
	tokens := ""
	tokens += fmt.Sprintf("flags:  %8b\n", tr.Flags)
	tokens += fmt.Sprintf("trainer class: %d\n", tr.Class)
	tokens += fmt.Sprintf("battle type 1: %d\n", tr.BattleType)
	tokens += fmt.Sprintf("# pokemons:    %d\n", tr.NumPokemon)

	itemIds := [4]uint16{tr.Item1, tr.Item2, tr.Item3, tr.Item4}
	for i, itemId := range itemIds {
		if itemId != 0 {
			tokens += fmt.Sprintf("trainer item ID #%d: %d\n", i, itemId)
		}
	}

	tokens += fmt.Sprintf("AI Model:      %d\n", tr.AiType)
	tokens += fmt.Sprintf("battle type 2: %d\n", tr.BattleType2)

	if len(tr.PartyPokemon) > 0 {
		tokens += "-- Pokemons --"

		for _, p := range tr.PartyPokemon {
			tokens += fmt.Sprintf("%v", p)
		}
	}

	return tokens + "========"
}

func RipTrainerDataPlat() []PlatTrainer {
	ret := make([]PlatTrainer, 0)

	path := path_resolver.GetRoot() + "/roms/pkmn-pt.nds"
	f, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	trainerData := narc.NewNarcFile(f, 0x0371E000)
	trainerPokemon := narc.NewNarcFile(f, 0x03724600)

	trainerId := 0
	numPokemonsRead := 0

	for _, file := range trainerData.FrameFATB.Data.Entries {
		// read offsets in FIMG buffer
		trainerBuf := trainerData.FrameFIMG.Data.Data[file.Start:file.End]
		// fmt.Println(trainerBuf)
		w := walker.NewWalker(trainerBuf)

		trainer := PlatTrainer{
			Flags:        w.U8(),
			Class:        w.U8(),
			BattleType:   w.U8(),
			NumPokemon:   w.U8(),
			Item1:        w.U16(),
			Item2:        w.U16(),
			Item3:        w.U16(),
			Item4:        w.U16(),
			AiType:       w.U32(),
			BattleType2:  w.U32(),
			PartyPokemon: make([]TrainerPokemon, 0),
		}

		partyEntryBytes := trainerPokemon.FrameFIMG.Data.Data

		for range trainer.NumPokemon {
			if numPokemonsRead == len(trainerPokemon.FrameFATB.Data.Entries) {
				break
			}
			pkmnFile := trainerPokemon.FrameFATB.Data.Entry(numPokemonsRead)
			pokemonBuf := partyEntryBytes[pkmnFile.Start:pkmnFile.End]
			fmt.Printf("%03d|  % x\n", numPokemonsRead, pokemonBuf)
			pw := walker.NewWalker(pokemonBuf)

			pkmn := TrainerPokemon{
				Difficulty: pw.U16(),
				Level:      pw.U16(),
				Species:    pw.U16(),
				Seal:       pw.U16(),
				MoveIds:    make([]*uint16, 0),
			}

			switch trainer.Flags {
			case 0:
				{
					// fmt.Println("type 0, not adding any more flags")
					break
				}
			case 1:
				{
					// fmt.Println("Case 1 - custom moveset")
					for range 4 {
						moveId := pw.U16()
						if moveId != 0 {
							pkmn.MoveIds = append(pkmn.MoveIds, &moveId)
						}
					}
					break
				}
			case 2:
				{
					// fmt.Println("Case 2 - item")
					itemId := pw.U16()
					pkmn.ItemId = &itemId
					break
				}
			case 3:
				{
					// fmt.Println("Case 3 - item + custom moveset")
					itemId := pw.U16()
					for range 4 {
						moveId := pw.U16()
						if moveId != 0 {
							pkmn.MoveIds = append(pkmn.MoveIds, &moveId)
						}
					}
					pkmn.ItemId = &itemId
					break
				}
			default:
				{
					fmt.Println("uh oh....")
				}
			}

			// ivs := uint32(pkmn.Difficulty) * 31 / 255
			numPokemonsRead += 1
			trainer.PartyPokemon = append(trainer.PartyPokemon, pkmn)
		}

		trainerId += 1
		ret = append(ret, trainer)
	}

	return ret
}
