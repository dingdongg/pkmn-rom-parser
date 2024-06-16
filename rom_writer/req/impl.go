package req

import (
	"encoding/binary"
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
	wr.Contents[ITEM] = WriteUint{item.Index, 2}
}

func (wr WriteRequest) WriteAbility(ability string) {
	abilityMap := data.GenerateAbilityMap()
	abilityId, ok := abilityMap[ability]
	if !ok {
		log.Fatalf("failed to write ability '%s'; doesn't exist\n", ability)
	}
	wr.Contents[ABILITY] = WriteUint{abilityId, 1}
}

// TODO improve signature. since golang doesn't does support struct spreading like JS, it's
// clunky to use like this
func (wr WriteRequest) WriteBattleStats(hp, atk, def, spa, spd, spe uint) {
	wr.Contents[BATTLE_STATS] = WriteStats{hp, atk, def, spa, spd, spe}
}

func (wr WriteRequest) WriteEV(hp, atk, def, spa, spd, spe uint) {
	wr.Contents[EV] = WriteEffortValue{hp, atk, def, spa, spd, spe}
}

func (wr WriteRequest) WriteIV(hp, atk, def, spa, spd, spe uint) {
	wr.Contents[IV] = WriteIndivValue{hp, atk, def, spa, spd, spe}
}

func (wr WriteRequest) WriteNickname(name string) {
	wr.Contents[NICKNAME] = WriteString{name}
}

func (wr WriteRequest) WriteLevel(level uint) {
	wr.Contents[LEVEL] = WriteUint{level, 1}
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

func (wev WriteEffortValue) Bytes() ([]byte, error) {
	buf := make([]byte, 6)
	stats := [6]uint{wev.Hp, wev.Attack, wev.Defense, wev.Speed, wev.SpAttack, wev.SpDefense}

	for i, s := range stats {
		if s > 255 {
			return []byte{}, fmt.Errorf("effort value must be <= 255")
		}
		buf[i] = byte(s)
	}

	return buf, nil
}

func (wiv WriteIndivValue) Bytes() ([]byte, error) {
	buf := uint32(0)
	stats := [6]uint{wiv.Hp, wiv.Attack, wiv.Defense, wiv.Speed, wiv.SpAttack, wiv.SpDefense}
	
	for i, s := range stats {
		if s > 31 {
			return []byte{} , fmt.Errorf("individual value must be <= 31")
		}
		masked := s & 0b11111
		buf |= uint32(masked << (i * 5))
	}
	
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, buf)
	return res, nil
}

func (wu WriteUint) Bytes() ([]byte, error) {
	res := make([]byte, 0)

	if wu.NumBytes == 1 {
		res = append(res, byte(wu.Val))
	} else if wu.NumBytes == 2 {
		res = binary.LittleEndian.AppendUint16(res, uint16(wu.Val))
	} else if wu.NumBytes == 4 {
		res = binary.LittleEndian.AppendUint32(res, uint32(wu.Val))
	} else if wu.NumBytes == 8 {
		res = binary.LittleEndian.AppendUint64(res, uint64(wu.Val))
	} else {
		return res, fmt.Errorf("too many bytes")
	}

	return res, nil
}

func (ws WriteString) Bytes() ([]byte, error) {
	// gen. 4 NDS games allow up to 10 characters max (so 11 with null terminator)
	if len(ws.Val) > 10 {
		return []byte{}, fmt.Errorf("string can only be max 10 bytes long")
	}

	res := make([]byte, 0)

	for _, r := range ws.Val {
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