package savefile

import (
	"fmt"
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
		expTable: ripper.RipExpTableGen5(),
		pokemonMetadata: ripper.RipPokemonDataGen5(),
		itemTable: ripper.RipItemNamesGen5(),
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

	form := (genderByte >> 3) & 0x1F

	moves := make([]models.Move, 0)
	for i := range 0x4 {
		id := utils.U16(b, i*0x2)
		move := models.Move{
			Id:   id,
			Name: bw.moveNames[id],
		}
		moves = append(moves, move)
	}

	pokedexId := utils.U16(a, 0x0)
	metadata := bw.pokemonMetadata[pokedexId]

	return models.Pokemon{
		Name:      string(runes),
		PokedexId: pokedexId,
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
		Form: form,
		Base: metadata.Base,
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

func (bw *BwSavefile) validatePokemon(p *models.Pokemon) error {
	newError := func(msg string, params ...any) error {
		return fmt.Errorf(fmt.Sprint("VALIDATION ERR: ", msg), params)
	}

	if p.PokedexId > 649 {
		return newError("invalid pokedex id %d", p.PokedexId)
	}

	// TODO: improve item validation (currently, just based on name which is very brittle)
	itemFound := false
	idx := 0
	for _, itemName := range bw.itemTable {
		if itemName == p.HeldItem {
			itemFound = true
			break
		}
		idx += 1
	}

	if !itemFound {
		return newError("invalid item ID %d", idx)
	}

	abilityMap := data.GenerateAbilityMap()
	_, ok := abilityMap[p.Ability]
	if !ok {
		return newError("invalid ability '%s'", p.Ability)
	}

	if err := p.EV.AssertBound(255); err != nil {
		return newError(err.Error())
	}

	// EXP validation
	growthType := bw.pokemonMetadata[p.PokedexId].GrowthType
	expTable := bw.expTable[growthType]

	if p.Exp < expTable[0] || p.Exp > expTable[100] {
		return newError("EXP points out of bounds for pokemon #%d", p.PokedexId)
	}

	var binarySearch func(buf ripper.ExperienceTable, lo, hi int, target uint32) int
	binarySearch = func(buf ripper.ExperienceTable, lo, hi int, target uint32) int {
		if hi-lo == 1 {
			return lo
		}

		mid := (lo + hi) >> 1

		if buf[mid] < target {
			return binarySearch(buf, mid, hi, target)
		} else if buf[mid] > target {
			return binarySearch(buf, lo, mid, target)
		}

		return mid
	}

	p.Level = uint8(binarySearch(expTable, 0, len(expTable), p.Exp))
	fmt.Printf("adjusted level (exp=%d): %d\n", p.Exp, p.Level)

	numMoves := uint16(len(bw.moveNames))
	for _, move := range p.Moves {
		if move.Id > numMoves {
			return newError("invalid move '%s' (id=%d)", move.Name, move.Id)
		}
	}

	// IV validation
	if err := p.IV.AssertBound(31); err != nil {
		return newError(err.Error())
	}

	// gender bit validation
	if p.Gender.String() == "Unknown" {
		return newError("invalid gender: %d", p.Gender)
	}

	altFormPokemons := []utils.Pair[uint16, uint8]{
		utils.NewPair[uint16, uint8](201, 28), utils.NewPair[uint16, uint8](386, 4),
		utils.NewPair[uint16, uint8](412, 3), utils.NewPair[uint16, uint8](413, 3),
		utils.NewPair[uint16, uint8](422, 2), utils.NewPair[uint16, uint8](423, 2),
		utils.NewPair[uint16, uint8](479, 6), utils.NewPair[uint16, uint8](487, 2),
		utils.NewPair[uint16, uint8](492, 2), utils.NewPair[uint16, uint8](493, 18),
	}

	for _, pair := range altFormPokemons {
		id, numForms := pair.First, pair.Second
		if id != p.PokedexId {
			continue
		}

		if p.Form > numForms {
			return newError("Form ID for Pokemon #%d must not exceed %d", p.PokedexId, numForms)
		}
	}

	if len(p.Name) == 0 || len(p.Name) > 10 {
		return newError("name must be between 1-10 characters long")
	}

	if p.Level > 100 {
		return newError("level %d is too big; cannot exceed 100", p.Level)
	}

	// in generation 5, only the speed stat has a unique cap
	// that isn't applied to any other battle stat
	// stat = min(stat, 10000)
	// stat = min(stat, 8192) --> i'm assuming this is for speed comparisons after taking into account things like prio
	// https://bulbapedia.bulbagarden.net/wiki/Stat_modifier
	if p.Battle.Speed > 10000 {
		return newError("speed stat (%d) cannot exceed 10000", p.Battle.Speed)
	} else if p.Battle.Speed >= 8192 {
		return newError("speed stat (%d) cannot exceed 8192", p.Battle.Speed)
	}

	return nil
}

func (bw *BwSavefile) validate() error {
	if len(bw.partyPokemon) > 6 {
		return fmt.Errorf("VALIDATION ERR: party cannot hold more than 6 pokemon")
	}

	for _, p := range bw.partyPokemon {
		if err := bw.validatePokemon(p); err != nil {
			return err
		}
	}

	return nil
}

func (bw *BwSavefile) Version() enums.GameVersion {
	return enums.BW
}

func (bw *BwSavefile) Flush() error {
	return nil
}
