package req

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts"
	"github.com/dingdongg/pkmn-rom-parser/v7/shuffler"
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

func NewWriteRequest(partyIndex uint) WriteRequest {
	return WriteRequest{
		partyIndex,
		make(NewData),
	}
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
