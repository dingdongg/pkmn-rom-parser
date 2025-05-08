package savefile

import (
	"fmt"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v7/char"
	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/data"

	"os"

	"github.com/dingdongg/pkmn-rom-parser/v7/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/path_resolver"
	"github.com/dingdongg/pkmn-rom-parser/v7/ripper"
	"github.com/dingdongg/pkmn-rom-parser/v7/shuffler"
	"github.com/dingdongg/pkmn-rom-parser/v7/utils"
	"github.com/dingdongg/pkmn-rom-parser/v7/validator"
)

func NewPlatSavefile(bytes []byte) *PlatSavefile {
	// find the latest savefile offset here and persist it as member variable
	// assume that `bytes` is the full savefile (includes backup)
	// include a new field that points to just the recent portion of the savefile
	// to do this, we need to figure out which block is more "recent"
	sbStart, sbEnd := uint(0x0), uint(0xCF2B)
	latestBlock, err := validator.LatestSmallBlock(bytes, utils.NewRange(sbStart, sbEnd))
	if err != nil {
		log.Fatalln("invalid savefile")
	}

	return &PlatSavefile{
		rawBytes:        bytes,       // encrypted
		latestSave:      latestBlock, // encrypted
		partyPokemon:    make([]*models.Pokemon, 0),
		moveNames:       ripper.RipMoveNames(),
		rawParty:        make([]byte, 0), // decrypted
		expTable:        ripper.RipExpTableGen4(),
		pokemonMetadata: ripper.RipPokemonData(),
	}
}

func (pt *PlatSavefile) parsePokemon(index int) models.Pokemon {
	offset := 0xA0 + index*236
	savefile := pt.latestSave.Data()
	raw := crypt.DecryptPokemon(savefile[offset : offset+236])
	pt.rawParty = append(pt.rawParty, raw...)

	blocks := shuffler.GetPokemonBlocks(raw)
	a, b, c := blocks[0], blocks[1], blocks[2]

	rawName := c[:0x16]
	name := ""
	for i := 0; i < len(rawName); i += 2 {
		code := utils.U16(rawName, i)
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

	ability, _ := data.GetAbility(uint(utils.U8(a, 0xD)))
	item, _ := data.GetItem(utils.U16(a, 0x2))

	ivBuffer := utils.U32(b, 0x10)
	getIv := func(statIndex int) uint8 {
		val := (ivBuffer >> (5 * statIndex)) & 0x1F
		return uint8(val)
	}

	genderByte := utils.U8(b, 0x18)
	gender := enums.Male
	if genderByte&0x2 != 0 {
		gender = enums.Female
	} else if genderByte&0x4 != 0 {
		gender = enums.Unknown
	}

	form := (genderByte >> 3) & 0x1F

	moves := make([]models.Move, 0)
	for i := range 0x4 {
		id := utils.U16(b, i*0x2)
		move := models.Move{
			Id:   id,
			Name: pt.moveNames[id],
		}
		moves = append(moves, move)
	}

	/*
		missing:
		- base stats
	*/
	return models.Pokemon{
		Name:      name,
		PokedexId: utils.U16(a, 0x0),
		Exp:       utils.U32(a, 0x8),
		Level:     utils.U8(battleStats, 0x4),
		Nature:    enums.Nature(utils.U32(raw, 0) % 25),
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
	}
}

func (pt *PlatSavefile) PartyPokemon() []*models.Pokemon {
	partySize := int(utils.U32(pt.latestSave.Data(), 0xA0-0x4))
	for i := range partySize {
		pkmn := pt.parsePokemon(i)
		pt.partyPokemon = append(pt.partyPokemon, &pkmn)
	}

	return pt.partyPokemon
}

func (pt *PlatSavefile) validatePokemon(p *models.Pokemon) error {
	newError := func(msg string, params ...any) error {
		return fmt.Errorf(fmt.Sprint("VALIDATION ERR: ", msg), params)
	}

	if p.PokedexId > 493 {
		return newError("invalid pokedex id %d", p.PokedexId)
	}
	itemMap := data.GenerateItemMap()
	item, ok := itemMap[p.HeldItem]
	if !ok {
		return newError("invalid item '%s'", p.HeldItem)
	} else if item.Exclusivity == "HGSS" {
		return newError("item '%s' not supported in PLAT", p.HeldItem)
	}

	abilityMap := data.GenerateAbilityMap()
	ability, ok := abilityMap[p.Ability]
	if !ok {
		return newError("invalid ability '%s'", p.Ability)
	} else if ability > 123 {
		return newError("ability '%s' not supported in PLAT", p.Ability)
	}

	// EV validation - nothign to do
	if err := p.EV.AssertBound(255); err != nil {
		return newError(err.Error()) // theoretically should never hit this due to static uint8 bound checks
	}

	// EXP validation
	growthType := pt.pokemonMetadata[p.PokedexId].GrowthType
	// contains total EXP required to reach each level (has an entry for lvl 0 for some reason)
	expTable := pt.expTable[growthType]
	if p.Exp < expTable[0] || p.Exp > expTable[100] {
		return newError("Experience points out of bounds for pokemon #%d", p.PokedexId)
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

	// moveset ID validation
	for _, move := range p.Moves {
		if move.Id > 467 {
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

	// form validation
	altFormPokemons := []utils.Pair[uint16, uint8]{
		utils.NewPair[uint16, uint8](201, 28), utils.NewPair[uint16, uint8](386, 4),
		utils.NewPair[uint16, uint8](412, 3), utils.NewPair[uint16, uint8](413, 3),
		utils.NewPair[uint16, uint8](422, 2), utils.NewPair[uint16, uint8](423, 2),
		utils.NewPair[uint16, uint8](479, 6), utils.NewPair[uint16, uint8](487, 2),
		utils.NewPair[uint16, uint8](492, 2), utils.NewPair[uint16, uint8](493, 18),
		// no pichu, only in HGSS
	}
	for _, pair := range altFormPokemons {
		if pair.First != p.PokedexId {
			continue
		}

		if p.Form > pair.Second {
			return newError("Form ID for Pokemon #%d must not exceed %d", p.PokedexId, pair.Second)
		}
	}

	// name validation - length should be 10 characters max (excluding end-of-string terminator)
	if len(p.Name) == 0 || len(p.Name) > 10 {
		return newError("name must be between 1-10 characters long")
	}

	// level validation ?
	// when both exp and level are modified, one has to be overridden... which one?
	// have level overwritten if EXP is also modified
	if p.Level > 100 {
		return newError("level %d is too big; cannot exceed 100", p.Level)
	}

	// battle stats validation
	if err := p.Battle.AssertBound(999); err != nil {
		return newError(err.Error())
	}

	return nil
}

func (pt *PlatSavefile) validate() error {
	if len(pt.partyPokemon) > 6 {
		return fmt.Errorf("VALIDATION ERR: party cannot hold more than 6 pokemon")
	}

	for _, p := range pt.partyPokemon {
		if err := pt.validatePokemon(p); err != nil {
			return err
		}
	}

	return nil
}

func (pt *PlatSavefile) Version() enums.GameVersion {
	return enums.PLAT
}

func (pt *PlatSavefile) updatePokemon(index int, p *models.Pokemon) {
	start := index * 236
	buf := pt.rawParty[start : start+236]
	pid := utils.U32(buf, 0)

	a, _ := shuffler.GetPokemonBlockLocation(shuffler.A, pid)
	b, _ := shuffler.GetPokemonBlockLocation(shuffler.B, pid)
	c, _ := shuffler.GetPokemonBlockLocation(shuffler.C, pid)

	A, B, C := buf[a:a+32], buf[b:b+32], buf[c:c+32]

	utils.WriteU16(A, 0x0, p.PokedexId)
	itemMap := data.GenerateItemMap()
	if itemId, ok := itemMap[p.HeldItem]; ok {
		utils.WriteU16(A, 0x2, uint16(itemId.Index))
	}

	utils.WriteU32(A, 0x8, p.Exp)

	abilityMap := data.GenerateAbilityMap()
	if ability, ok := abilityMap[p.Ability]; ok {
		utils.WriteU16(A, 0xD, uint16(ability))
	}

	A[0x10] = p.EV.Hp
	A[0x11] = p.EV.Attack
	A[0x12] = p.EV.Defense
	A[0x13] = p.EV.Speed
	A[0x14] = p.EV.SpeAttack
	A[0x15] = p.EV.SpeDefense

	for i, move := range p.Moves {
		utils.WriteU16(B, i*0x2, move.Id)
	}

	// pack IVs
	var packedIv uint32 = utils.U32(B, 0x10) & 0xC0_00_00_00
	packedIv |= uint32(p.IV.Hp & 0x1F)
	packedIv |= (uint32(p.IV.Attack&0x1F) << 5)
	packedIv |= (uint32(p.IV.Defense&0x1F) << 10)
	packedIv |= (uint32(p.IV.Speed&0x1F) << 15)
	packedIv |= (uint32(p.IV.SpeAttack&0x1F) << 20)
	packedIv |= (uint32(p.IV.SpeDefense&0x1F) << 25)
	utils.WriteU32(B, 0x10, packedIv)

	genderByte := B[0x18] & 0xF9
	switch p.Gender {
	case enums.Female:
		genderByte |= 0x02
	case enums.Unknown:
		genderByte |= 0x04
	}
	B[0x18] = genderByte

	// alternate forms. Ignored for pokemon without forms
	B[0x18] &= 0x1F
	B[0x18] |= (p.Form & 0x1F) << 3

	length := min(len(p.Name), 10)
	utils.Memset(C, 0x0, 0x16, 0xFF)
	for i := range length {
		letter := string(p.Name[i])
		code, err := char.Index(letter)
		if err != nil {
			code = char.END_OF_STRING // append null terminator instead and stop
		}
		utils.WriteU16(C, i*0x2, code)
	}

	battleStatBuf := buf[0x88:]
	battleStatBuf[0x4] = p.Level

	utils.WriteU16(battleStatBuf, 0x8, p.Battle.Hp)
	utils.WriteU16(battleStatBuf, 0xA, p.Battle.Attack)
	utils.WriteU16(battleStatBuf, 0xC, p.Battle.Defense)
	utils.WriteU16(battleStatBuf, 0xE, p.Battle.Speed)
	utils.WriteU16(battleStatBuf, 0x10, p.Battle.SpeAttack)
	utils.WriteU16(battleStatBuf, 0x12, p.Battle.SpeDefense)
}

func (pt *PlatSavefile) Flush() error {
	if err := pt.validate(); err != nil {
		return err
	}

	encryptedBuffer := make([]byte, 0)
	for i, p := range pt.partyPokemon {
		pt.updatePokemon(i, p)
		offset := i * 236
		ciphertext := crypt.EncryptPokemon(pt.rawParty[offset : offset+236])
		encryptedBuffer = append(encryptedBuffer, ciphertext...)
	}

	copy(pt.latestSave.Data()[0xA0:], encryptedBuffer)
	// update footer checksum
	newChecksum := crypt.CRC16_CCITT(pt.latestSave.Data())
	pt.latestSave.Footer.Checksum = newChecksum

	copy(pt.rawBytes[pt.latestSave.Offset():], pt.latestSave.Bytes())

	fullpath := path_resolver.GetRoot() + "/new-plat.sav"
	os.WriteFile(fullpath, pt.rawBytes, os.ModePerm)
	return nil
}
