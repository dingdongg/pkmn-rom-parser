package savefile

import (
	"fmt"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v7/char"
	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/data"

	// "os"
	// "github.com/dingdongg/pkmn-rom-parser/v7/path_resolver"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/models"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/ripper"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/validator"
	"github.com/dingdongg/pkmn-rom-parser/v7/shuffler"
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
		rawBytes:     bytes,						// encrypted
		latestSave: latestBlock,					// encrypted
		partyPokemon: make([]*models.Pokemon, 0), 
		moveNames: ripper.RipMoveNames(),
		rawParty: make([]byte, 0),					// decrypted
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
		val := (ivBuffer >> (5*statIndex)) & 0x1F
		return uint8(val)
	}

	genderByte := utils.U8(b, 0x18)
	gender := enums.Male
	if genderByte & 0x2 != 0 {
		gender = enums.Female
	} else if genderByte & 0x4 != 0 {
		gender = enums.Unknown
	}

	moves := make([]models.Move, 0)
	for i := range 0x4 {
		id := utils.U16(b, i*0x2)
		move := models.Move{
			Id: id,
			Name: pt.moveNames[id],
		}
		moves = append(moves, move)
	}

	/*
	missing: 
	- base stats
	- alternate forms
	*/
	return models.Pokemon{
		Name: name,
		PokedexId: utils.U16(a, 0x0),
		Exp: utils.U32(a, 0x8),
		Level: utils.U8(battleStats, 0x4),
		Nature: enums.Nature(utils.U32(raw, 0) % 25),
		Ability: ability,
		HeldItem: item.Name,
		Gender: gender,
		Moves: moves,
		EV: models.Stat[uint8]{
			Hp: utils.U8(a, 0x10),
			Attack: utils.U8(a, 0x11),
			Defense: utils.U8(a, 0x12),
			SpeAttack: utils.U8(a, 0x14),
			SpeDefense: utils.U8(a, 0x15),
			Speed: utils.U8(a, 0x13),
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
			Hp: utils.U16(battleStats, 0x8),
			Attack: utils.U16(battleStats, 0xA),
			Defense: utils.U16(battleStats, 0xC),
			SpeAttack: utils.U16(battleStats, 0x10),
			SpeDefense: utils.U16(battleStats, 0x12),
			Speed: utils.U16(battleStats, 0xE),
		},
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
	// EXP validation - TODO
	// check that exp is not over the maximum for its growth type

	// moveset ID validation
	for _, move := range p.Moves {
		if move.Id > 467 {
			return newError("invalid move '%s' (id=%d)", move.Name, move.Id)
		}
	}

	// IV validation
	if p.IV.Hp > 31 {
		return newError("HP IV (=%d) cannot exceed 31", p.IV.Hp)
	}
	if p.IV.Attack > 31 {
		return newError("ATTACK IV (=%d) cannot exceed 31", p.IV.Attack)
	}
	if p.IV.Defense > 31 {
		return newError("DEFENSE IV (=%d) cannot exceed 31", p.IV.Defense)
	}
	if p.IV.SpeAttack > 31 {
		return newError("SPECIAL ATK IV (=%d) cannot exceed 31", p.IV.SpeAttack)
	}
	if p.IV.SpeDefense > 31 {
		return newError("SPECIAL DEF IV (=%d) cannot exceed 31", p.IV.SpeDefense)
	}
	if p.IV.Speed > 31 {
		return newError("SPEED IV (=%d) cannot exceed 31", p.IV.Speed)
	}

	// gender bit validation
	if p.Gender.String() == "Unknown" {
		return newError("invalid gender: %d", p.Gender)
	}

	// form validation - TODO after implementing support for forms

	// name validation - length should be 10 characters max (excluding end-of-string terminator)
	if len(p.Name) == 0 || len(p.Name) > 10 {
		return newError("name must be between 1-10 characters long")
	}

	// level validation ?
	if p.Level > 100 {
		return newError("level %d is too big; cannot exceed 100", p.Level)
	}

	// battle stats validation
	if p.Battle.Hp > 999 {
		return newError("HP battle stat (=%d) cannot exceed 999", p.Battle.Hp)
	}
	if p.Battle.Attack > 999 {
		return newError("ATTACK battle stat (=%d) cannot exceed 999", p.Battle.Attack)
	}
	if p.Battle.Defense > 999 {
		return newError("DEFENSE battle stat (=%d) cannot exceed 999", p.Battle.Defense)
	}
	if p.Battle.SpeAttack > 999 {
		return newError("SPECIAL ATK battle stat (=%d) cannot exceed 999", p.Battle.SpeAttack)
	}
	if p.Battle.SpeDefense > 999 {
		return newError("SPECIAL DEF battle stat (=%d) cannot exceed 999", p.Battle.SpeDefense)
	}
	if p.Battle.Speed > 999 {
		return newError("SPEED battle stat (=%d) cannot exceed 999", p.Battle.Speed)
	}

	return nil
}

/*
this functino shsould be aimed at validating the internal pokemon data before flushing. 

TODO: complete this function, and write the move parsing logic for all concrete savefiles
*/
func (pt *PlatSavefile) validate() error {
	/*
	what does it mean to "validate" data before flushing?
	- the party pokemon field will have new and old data
	- checksum validation (has to be done again after flushing, though?)
	- check that the values in the party pokemon structs 
	  conform to the numeric limits imposed by the game
	  (ie. level cannot be greater than 100, valid move IDs, etc.)
	
	"Flushing data"
	- for sake of simplicity, we can "pack" the entire party pokemon contents
	  back into the savefile format.
	- then we need to encrypt these changes and update checksums, and return the 
	  updated savefile
	
	
	*/
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
	start := index*236
	buf := pt.rawParty[start : start+236]
	pid := utils.U32(buf, 0)

	a, _ := shuffler.GetPokemonBlockLocation(shuffler.A, pid)
	b, _ := shuffler.GetPokemonBlockLocation(shuffler.B, pid)
	c, _ := shuffler.GetPokemonBlockLocation(shuffler.C, pid)

	A, B, C := buf[a : a+32], buf[b : b+32], buf[c : c+32]

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
	packedIv |= (uint32(p.IV.Attack & 0x1F) << 5) 
	packedIv |= (uint32(p.IV.Defense & 0x1F) << 10) 
	packedIv |= (uint32(p.IV.Speed & 0x1F) << 15) 
	packedIv |= (uint32(p.IV.SpeAttack & 0x1F) << 20) 
	packedIv |= (uint32(p.IV.SpeDefense & 0x1F) << 25) 
	utils.WriteU32(B, 0x10, packedIv)

	genderByte := B[0x18] & 0xF9
	switch (p.Gender) {
	case enums.Female: genderByte |= 0x02
	case enums.Unknown: genderByte |= 0x04
	}
	B[0x18] = genderByte

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
	/*
	block A:
		pokedex id
		item id
		ability id
		EVs
		EXP

	block B:
		moveset ids
		IVs
		gender bits
		form byte?
	
	block C:
		name

	block D:
		none

	battle stats:
		level
		battle stats
	*/
}

func (pt *PlatSavefile) Flush() error {
	if err := pt.validate(); err != nil {
		return err
	}

	encryptedBuffer := make([]byte, 0)
	for i, p := range pt.partyPokemon {
		pt.updatePokemon(i, p)
		offset := i*236
		ciphertext := crypt.EncryptPokemon(pt.rawParty[offset : offset+236])
		encryptedBuffer = append(encryptedBuffer, ciphertext...)
	}

	copy(pt.latestSave.Data()[0xA0:], encryptedBuffer)
	// update footer checksum
	newChecksum := crypt.CRC16_CCITT(pt.latestSave.Data())
	pt.latestSave.Footer.Checksum = newChecksum

	copy(pt.rawBytes[pt.latestSave.Offset():], pt.latestSave.Bytes())

	// fullpath := path_resolver.GetRoot() + "/new-plat.sav"
	// os.WriteFile(fullpath, pt.rawBytes, os.ModePerm)
	return nil
}
