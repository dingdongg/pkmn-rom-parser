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

type AbsAddress uint

func UpdatePartyPokemon(savefile []byte, chunk validator.Chunk, newData WriteRequests) ([]byte, error) {
	updatedPokemonIndexes := make([]uint, 0)
	base := chunk.SmallBlock.Address + consts.PERSONALITY_OFFSET

	for req, data := range newData.Contents {
		bytes, err := data.Bytes()
		if err != nil {
			return []byte{}, err
		}
		fmt.Printf("%s: % x\n", req, bytes)

		offset := base + newData.PartyIndex * consts.PARTY_POKEMON_SIZE
		personality := binary.LittleEndian.Uint32(savefile[offset : offset+4])

		var writeLocation int
		var blockIndex int

		if req == ITEM {
			writeLocation = consts.BLOCK_A_ITEM
			blockIndex = shuffler.A
		} else if req == ABILITY {
			writeLocation = consts.BLOCK_A_ABILITY
			blockIndex = shuffler.A
		} else if req == EV {
			writeLocation = consts.BLOCK_A_EV
			blockIndex = shuffler.A
		} else if req == IV {
			writeLocation = consts.BLOCK_B_IV
			blockIndex = shuffler.B
		} else if req == NICKNAME {
			writeLocation = consts.BLOCK_C_NICKNAME
			blockIndex = shuffler.C
		} else if req == LEVEL {
			writeLocation = consts.BATTLE_STATS_LEVEL
			blockIndex = -1
		} else if req == BATTLE_STATS {
			writeLocation = consts.BATTLE_STATS_STAT
			blockIndex = -1
		}

		var addr AbsAddress

		if blockIndex != -1 {
			blockAddress, err := shuffler.GetPokemonBlockLocation(uint(blockIndex), personality)
			if err != nil {
				return []byte{}, nil
			}

			addr = AbsAddress(offset + blockAddress + uint(writeLocation))
		} else {
			addr = AbsAddress(offset + 0x88 + uint(writeLocation))
		}

		fmt.Printf("Before: % x\n", savefile[addr : addr+20])
		copy(savefile[addr:], bytes)
		fmt.Printf("After:  % x\n", savefile[addr : addr+20])
		fmt.Print("-------\n\n")
		updatedPokemonIndexes = append(updatedPokemonIndexes, newData.PartyIndex)
	}

	// TODO: encrypt the modified small block (at per-pokemon level)
	for _, i := range updatedPokemonIndexes {
		pokemonOffset := base + i * consts.PARTY_POKEMON_SIZE
		encrypted := crypt.EncryptPokemon(savefile[pokemonOffset:])

		copy(savefile[AbsAddress(pokemonOffset):], encrypted)
	}

	start := chunk.SmallBlock.Address
	end := chunk.SmallBlock.Address + uint(chunk.SmallBlock.Footer.BlockSize) - 0x14
	newChecksum := crypt.CRC16_CCITT(savefile[start : end])
	// TODO: compute & update checksums at a party-pokemon level, and at the whole block-level
	fmt.Printf("Before (checksum): % x\n", savefile[end : end+20])

	binary.LittleEndian.PutUint16(savefile[end + 18:], newChecksum)

	fmt.Printf("After (checksum):  % x\n", savefile[end : end+20])
	fmt.Print("-------\n\n")

	return savefile, nil
}