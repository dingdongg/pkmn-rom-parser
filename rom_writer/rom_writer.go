package rom_writer

import (
	"encoding/binary"

	"github.com/dingdongg/pkmn-rom-parser/v3/char"
	"github.com/dingdongg/pkmn-rom-parser/v3/crypt"
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

Once decryption is decoupled, use this subpackage to decrypt the contents before writing changes

once changes have been written, encrypt the chunk again and validate the ciphertext with the footer
- would I have to overwrite the contents in the footer as well, since the ciphertext
(and thereby the checksum) changes? <-- PROBABLY YES


*/
type Writable interface {
	Bytes() ([]byte, error)
}

// for battle stats
type WriteStats struct {
	Hp uint
	Attack uint
	Defense uint
	SpAttack uint
	SpDefense uint
	Speed uint
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
	res := make([]byte, 12)
	stats := [6]uint{ ws.Hp, ws.Attack, ws.Defense, ws.Speed, ws.SpAttack, ws.SpDefense }

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
		res = append(res, byte(wu.Val & 0xFF))
		wu.Val >>= 8
	}

	return res, nil
}

func (ws WriteString) Bytes() ([]byte, error) {
	// TODO: may have to limit the max length of the incoming string
	res := make([]byte, 0)

	max := func(a int, b int) int {
		if a > b {
			return a
		}
		return b
	}

	end := max(len(ws.Val) - 1, 22)

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
	ITEM = "ITEM"					// value will be item ID
	ABILITY = "ABILITY"				// value will be ability ID
	EV = "EV"						// value will be 
	IV = "IV"
	NICKNAME = "NICKNAME"
	LEVEL = "LEVEL"
	BATTLE_STATS = "BATTLE_STATS"
)

type NewData map[string]Writable

type WriteRequests struct {
	PartyIndex uint
	Contents NewData
}

func UpdatePartyPokemon(ciphertext []byte, newData WriteRequests) ([]byte, error) {
	return []byte{}, nil

	// for req, data := range newData.Contents {
	// 	switch req {
	// 		case ITEM:
	// 		case ABILITY:
	// 		case IV:
	// 		case EV: {
	// 		}
		
	// 	}
	// 	offset := newData.PartyIndex * 236
	// 	bytes := ciphertext[offset:]
	// 	buffer := crypt.DecryptPokemon(bytes)
	// 	personality := binary.LittleEndian.Uint32(buffer[0:4])
	// 	battleStatsBuffer := crypt.DecryptBattleStats(ciphertext[offset + 0x88:], personality)
	// }
}