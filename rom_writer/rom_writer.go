package rom_writer

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v3/char"
	"github.com/dingdongg/pkmn-rom-parser/v3/consts"
	"github.com/dingdongg/pkmn-rom-parser/v3/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v3/shuffler"
	"github.com/dingdongg/pkmn-rom-parser/v3/validator"
)

/*
tasks required at a high level:
1. overwrite (decrypted) savefile with desired changes
2. encrypt the block data
3. update the footer with the new checksum
3. return encrypted savefile

The very first thing required will be to locate the offset to the most recent savefile
Pokemon games store a second backup save version to account for memory corruption.
If the first chunk is corrupted, then the backup save will be used
Since the backup only appears in the event of memory corruptions, changes written to the savefile
should be done to the latest chunk

Currently, decryption is tightly coupled with savefile reading.
Since decryption will be needed for both reads/writes, it would be good
to pull out the encryption/decryption functionality into its own separate sub-package
-> this reduces coupling

# Once decryption is decoupled, use this subpackage to decrypt the contents before writing changes

once changes have been written, encrypt the chunk again and validate the ciphertext with the footer
- would I have to overwrite the contents in the footer as well, since the ciphertext
(and thereby the checksum) changes? <-- PROBABLY YES
*/
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

func (ws WriteStats) Bytes() ([]byte, error) {
	res := make([]byte, 0)
	stats := [6]uint{ws.Hp, ws.Attack, ws.Defense, ws.Speed, ws.SpAttack, ws.SpDefense}

	for _, v := range stats {
		// little endian (least significant byte first)
		fmt.Printf("UM: %d (0x%x)\n", v, v)
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

	res = append(res, 0xFF) // how many "0xFF"s do I insert????

	return res, nil
}

const (
	ITEM         = "ITEM"    // value will be item ID
	ABILITY      = "ABILITY" // value will be ability ID
	EV           = "EV"      // value will be
	IV           = "IV"
	NICKNAME     = "NICKNAME"
	LEVEL        = "LEVEL"
	BATTLE_STATS = "BATTLE_STATS"
)

type NewData map[string]Writable

type WriteRequests struct {
	PartyIndex uint
	Contents   NewData
}

type WriteRequestBuilder struct {
	WriteRequests
}

func NewWriteRequest(partyIndex uint) WriteRequests {
	return WriteRequests{
		partyIndex,
		make(NewData),
	}
}

func (wr WriteRequests) WriteItem(itemId uint) {
	wr.Contents[ITEM] = WriteUint{itemId}
}

func (wr WriteRequests) WriteAbility(abilityId uint) {
	wr.Contents[ABILITY] = WriteUint{abilityId}
}

func (wr WriteRequests) WriteEV(value uint) {
	wr.Contents[EV] = WriteUint{value}
}

func (wr WriteRequests) WriteIV(value uint) {
	wr.Contents[IV] = WriteUint{value}
}

func (wr WriteRequests) WriteNickname(name string) {
	wr.Contents[NICKNAME] = WriteString{name}
}

func (wr WriteRequests) WriteLevel(level uint) {
	wr.Contents[LEVEL] = WriteUint{level}
}

func (wr WriteRequests) WriteBattleStats(hp, atk, def, spa, spd, spe uint) {
	wr.Contents[BATTLE_STATS] = WriteStats{hp, atk, def, spa, spd, spe}
}

type AbsAddress uint

type StagingBuffer struct {
	Address AbsAddress
	Updates []byte
}

type StagingMap map[uint]StagingBuffer

func UpdatePartyPokemon(savefile []byte, chunk validator.Chunk, newData WriteRequests) ([]byte, error) {
	updatedPokemonIndexes := make(map[uint]bool, 0)
	base := chunk.SmallBlock.Address + consts.PERSONALITY_OFFSET
	changes := make(StagingMap)

	for req, data := range newData.Contents {
		bytes, err := data.Bytes()
		if err != nil {
			return []byte{}, err
		}
		fmt.Printf("%s: % x\n", req, bytes)

		offset := base + newData.PartyIndex*consts.PARTY_POKEMON_SIZE
		personality := binary.LittleEndian.Uint32(savefile[offset : offset+4])

		// fmt.Printf("% x\n", savefile[offset:offset+36])

		var dataOffset int
		var blockIndex int

		if req == ITEM {
			dataOffset = consts.BLOCK_A_ITEM
			blockIndex = shuffler.A
		} else if req == ABILITY {
			dataOffset = consts.BLOCK_A_ABILITY
			blockIndex = shuffler.A
		} else if req == EV {
			dataOffset = consts.BLOCK_A_EV
			blockIndex = shuffler.A
		} else if req == IV {
			dataOffset = consts.BLOCK_B_IV
			blockIndex = shuffler.B
		} else if req == NICKNAME {
			dataOffset = consts.BLOCK_C_NICKNAME
			blockIndex = shuffler.C
		} else if req == LEVEL {
			dataOffset = consts.BATTLE_STATS_LEVEL
			blockIndex = -1
		} else if req == BATTLE_STATS {
			dataOffset = consts.BATTLE_STATS_STAT
			blockIndex = -1
		}

		var addr AbsAddress

		if blockIndex != -1 {
			blockAddress, err := shuffler.GetPokemonBlockLocation(uint(blockIndex), personality)
			if err != nil {
				return []byte{}, err
			}

			addr = AbsAddress(offset + blockAddress + uint(dataOffset))

			stagingBuf, ok := changes[newData.PartyIndex]
			if !ok {
				changes[newData.PartyIndex] = StagingBuffer{
					addr,
					crypt.DecryptPokemon(savefile[offset:]), // entire single party pokemon (236B)
				}

				copy(changes[newData.PartyIndex].Updates[blockAddress+uint(dataOffset):], bytes)
			} else {
				copy(stagingBuf.Updates[blockAddress+uint(dataOffset):], bytes)
			}
		} else {
			addr = AbsAddress(offset + 0x88 + uint(dataOffset))

			stagingBuf, ok := changes[newData.PartyIndex]
			if !ok {
				changes[newData.PartyIndex] = StagingBuffer{
					addr,
					crypt.DecryptPokemon(savefile[offset:]), // entire single party pokemon (236B)
				}

				copy(changes[newData.PartyIndex].Updates[0x88+uint(dataOffset):], bytes)
			} else {
				copy(stagingBuf.Updates[0x88+uint(dataOffset):], bytes)
			}
		}

		fmt.Printf(
			"uipdates! (len=%d): % x\n", 
			len(changes[newData.PartyIndex].Updates), 
			changes[newData.PartyIndex].Updates,
		)

		// fmt.Printf("Before: % x\n", savefile[addr : addr+20])
		// copy(savefile[addr:], bytes)
		// fmt.Printf("After:  % x\n", savefile[addr : addr+20])
		// fmt.Print("-------\n\n")

		if _, seen := updatedPokemonIndexes[newData.PartyIndex]; !seen {
			updatedPokemonIndexes[newData.PartyIndex] = true
		}
	}

	for i := range updatedPokemonIndexes {
		pokemonOffset := base + i*consts.PARTY_POKEMON_SIZE
		fmt.Printf("changes for partyPokemon[%d]: % x\n", i, changes[i].Updates)
		encrypted := crypt.EncryptPokemon(changes[i].Updates)

		copy(savefile[AbsAddress(pokemonOffset):], encrypted)
	}

	start := chunk.SmallBlock.Address
	end := chunk.SmallBlock.Address + uint(chunk.SmallBlock.Footer.BlockSize) - 0x14
	newChecksum := crypt.CRC16_CCITT(savefile[start:end])

	fmt.Printf("Before (checksum): % x\n", savefile[end:end+20])
	binary.LittleEndian.PutUint16(savefile[end+18:], newChecksum)
	fmt.Printf("After (checksum):  % x\n", savefile[end:end+20])
	fmt.Print("-------\n\n")

	return savefile, nil
}
