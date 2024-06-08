package req

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v4/char"
	"github.com/dingdongg/pkmn-rom-parser/v4/consts"
	"github.com/dingdongg/pkmn-rom-parser/v4/shuffler"
)

const (
	ITEM         = "ITEM"
	ABILITY      = "ABILITY"
	EV           = "EV"
	IV           = "IV"
	NICKNAME     = "NICKNAME"
	LEVEL        = "LEVEL"
	BATTLE_STATS = "BATTLE_STATS"
)

type Writable interface {
	Bytes() ([]byte, error)
}

// for battle stats
type WriteStats struct {
	Hp        uint
	Attack    uint
	Defense   uint
	SpAttack  uint
	SpDefense uint
	Speed     uint
}

// for IDs/level/IVs/EVs
type WriteUint struct {
	Val uint
}

// for nicknames
type WriteString struct {
	Val string
}

type NewData map[string]Writable

type WriteRequest struct {
	PartyIndex uint
	Contents   NewData
}

func NewWriteRequest(partyIndex uint) WriteRequest {
	return WriteRequest{
		partyIndex,
		make(NewData),
	}
}

func (wr WriteRequest) WriteItem(itemId uint) {
	wr.Contents[ITEM] = WriteUint{itemId}
}

func (wr WriteRequest) WriteAbility(abilityId uint) {
	wr.Contents[ABILITY] = WriteUint{abilityId}
}

func (wr WriteRequest) WriteEV(value uint) {
	wr.Contents[EV] = WriteUint{value}
}

func (wr WriteRequest) WriteIV(value uint) {
	wr.Contents[IV] = WriteUint{value}
}

func (wr WriteRequest) WriteNickname(name string) {
	wr.Contents[NICKNAME] = WriteString{name}
}

func (wr WriteRequest) WriteLevel(level uint) {
	wr.Contents[LEVEL] = WriteUint{level}
}

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

	for wu.Val != 0 {
		res = append(res, byte(wu.Val&0xFF))
		wu.Val >>= 8
	}

	return res, nil
}

func (ws WriteString) Bytes() ([]byte, error) {
	// TODO: may have to limit the max length of the incoming string
	res := make([]byte, 0)

	min := func(a int, b int) int {
		if a < b {
			return a
		}
		return b
	}

	end := min(len(ws.Val), 11)

	for i := 0; i < end; i++ {
		c := ws.Val[i]
		index, err := char.Index(string(c))
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

// blockIndex will be -1 if the request is invalid, or is a level/battle stat request
func GetWriteLocation(request string) (dataOffset int, blockIndex int, err error) {
	if request == ITEM {
		dataOffset = consts.BLOCK_A_ITEM
		blockIndex = shuffler.A
	} else if request == ABILITY {
		dataOffset = consts.BLOCK_A_ABILITY
		blockIndex = shuffler.A
	} else if request == EV {
		dataOffset = consts.BLOCK_A_EV
		blockIndex = shuffler.A
	} else if request == IV {
		dataOffset = consts.BLOCK_B_IV
		blockIndex = shuffler.B
	} else if request == NICKNAME {
		dataOffset = consts.BLOCK_C_NICKNAME
		blockIndex = shuffler.C
	} else if request == LEVEL {
		dataOffset = consts.BATTLE_STATS_LEVEL
		blockIndex = -1
	} else if request == BATTLE_STATS {
		dataOffset = consts.BATTLE_STATS_STAT
		blockIndex = -1
	} else {
		return -1, -1, fmt.Errorf("invalid write request '%s'", request)
	}

	return dataOffset, blockIndex, nil
}
