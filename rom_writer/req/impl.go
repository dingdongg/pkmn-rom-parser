package req

import (
	"fmt"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v5/char"
	"github.com/dingdongg/pkmn-rom-parser/v5/data"
)

func (wr WriteRequest) WriteItem(itemName string) {
	itemMap := data.GenerateItemMap()
	item, ok := itemMap[itemName]
	if !ok {
		log.Fatalf("failed to write item '%s'; doesn't exist\n", itemName)
	}
	wr.Contents[ITEM] = WriteUint{item.Index, 2, 8}
}

func (wr WriteRequest) WriteAbility(ability string) {
	abilityMap := data.GenerateAbilityMap()
	abilityId, ok := abilityMap[ability]
	if !ok {
		log.Fatalf("failed to write ability '%s'; doesn't exist\n", ability)
	}
	wr.Contents[ABILITY] = WriteUint{abilityId, 1, 8}
}

func (wr WriteRequest) WriteEV(value CompressibleStat) {
	wr.Contents[EV] = WriteUint{value.Compress(8), 6, 8}
}

func (wr WriteRequest) WriteIV(value CompressibleStat) {
	fmt.Printf("0b %b\n", value.Compress(5))
	// reverse the bits, since memory is stored in little endian
	// doesn't have to be reversed for EVs/battles stats,
	// since each item (hp, atk, def, etc.) is stored in little endian
	// but preserves order (ie. hp-atk-def is NOT written as def-atk-hp)
	b := byte(value.Compress(5))
	b = ((b & 0xF0) >> 4) | ((b & 0x0F) << 4)
	b = ((b & 0xCC) >> 2) | ((b & 0x33) << 2)
	b = ((b & 0xAA) >> 1) | ((b & 0x55) << 1)

	wr.Contents[IV] = WriteUint{uint(b), 6, 5}
}

func (wr WriteRequest) WriteNickname(name string) {
	wr.Contents[NICKNAME] = WriteString{name}
}

func (wr WriteRequest) WriteLevel(level uint) {
	wr.Contents[LEVEL] = WriteUint{level, 1, 8}
}

// TODO improve signature. since golang doesn't does support struct spreading like JS, it's
// clunky to use like this
func (wr WriteRequest) WriteBattleStats(hp, atk, def, spa, spd, spe uint) {
	wr.Contents[BATTLE_STATS] = WriteStats{hp, atk, def, spa, spd, spe}
}

func (ws WriteStats) Bytes() ([]byte, error) {
	res := make([]byte, 0)
	stats := [6]uint{ws.Hp, ws.Attack, ws.Defense, ws.Speed, ws.SpAttack, ws.SpDefense}

	for _, v := range stats {
		// little endian (least significant byte first)
		little := byte(v & 0xFF)
		big := byte((v >> 8) & 0xFF)
		res = append(res, little, big)
	}

	return res, nil
}

func (wu WriteUint) Bytes() ([]byte, error) {
	res := make([]byte, 0)
	bitmask := uint(1 << wu.ItemBits) - 1

	for i := uint(0); i < wu.NumItems; i++ {
		item := wu.Val >> (i * wu.ItemBits)
		b := byte(item & bitmask)
		res = append(res, b)
	}

	return res, nil
}

func (ws WriteString) Bytes() ([]byte, error) {
	res := make([]byte, 0)

	min := func(a int, b int) int {
		if a < b {
			return a
		}
		return b
	}

	// gen. 4 NDS games allow up to 10 characters max (so 11 with null terminator)
	prunedString := ws.Val[:min(len(ws.Val), 11)]

	for _, r := range prunedString {
		index, err := char.Index(string(r))
		if err != nil {
			return []byte{}, err
		}

		little := byte(index & 0xFF)
		big := byte((index >> 8) & 0xFF)
		res = append(res, little, big)
	}

	res = append(res, 0xFF, 0xFF)

	// need to fill the (22 - len(res)) elements with 0s
	for len(res) != 22 {
		res = append(res, 0x0)
	}

	return res, nil
}