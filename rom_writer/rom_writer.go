package rom_writer

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v3/consts"
	"github.com/dingdongg/pkmn-rom-parser/v3/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v3/rom_writer/req"
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

type AbsAddress uint

type StagingBuffer []byte

// maps a party pokemon index to the updated pokemon data structure
type StagingMap map[uint]StagingBuffer

func UpdatePartyPokemon(savefile []byte, chunk validator.Chunk, newData []req.WriteRequest) ([]byte, error) {
	updatedPokemonIndexes := make(map[uint]bool, 0)
	base := chunk.SmallBlock.Address + consts.PERSONALITY_OFFSET
	changes := make(StagingMap)

	for _, wr := range newData {
		for request, data := range wr.Contents {
			bytes, err := data.Bytes()
			if err != nil {
				return []byte{}, err
			}

			offset := base + wr.PartyIndex*consts.PARTY_POKEMON_SIZE
			personality := binary.LittleEndian.Uint32(savefile[offset : offset+4])

			dataOffset, blockIndex, err := req.GetWriteLocation(request)
			if err != nil {
				return []byte{}, err
			}

			var blockAddress uint = 0x88

			if blockIndex != -1 {
				blockAddress, err = shuffler.GetPokemonBlockLocation(uint(blockIndex), personality)
				if err != nil {
					return []byte{}, err
				}
			}
			
			if _, ok := changes[wr.PartyIndex]; !ok {
				changes[wr.PartyIndex] = crypt.DecryptPokemon(savefile[offset:]) 
			} 
			copy(changes[wr.PartyIndex][blockAddress+uint(dataOffset):], bytes)

			if _, seen := updatedPokemonIndexes[wr.PartyIndex]; !seen {
				updatedPokemonIndexes[wr.PartyIndex] = true
			}
		}
	}

	for i := range updatedPokemonIndexes {
		pokemonOffset := base + i*consts.PARTY_POKEMON_SIZE
		fmt.Printf("changes for partyPokemon[%d]: % x\n", i, changes[i])
		encrypted := crypt.EncryptPokemon(changes[i])

		copy(savefile[AbsAddress(pokemonOffset):], encrypted)
	}

	updateBlockChecksum(savefile, chunk)
	return savefile, nil
}

func updateBlockChecksum(savefile []byte, chunk validator.Chunk) {
	start := chunk.SmallBlock.Address
	end := chunk.SmallBlock.Address + uint(chunk.SmallBlock.Footer.BlockSize) - 0x14
	newChecksum := crypt.CRC16_CCITT(savefile[start:end])

	fmt.Printf("Before (checksum): % x\n", savefile[end:end+20])
	binary.LittleEndian.PutUint16(savefile[end+18:], newChecksum)
	fmt.Printf("After (checksum):  % x\n", savefile[end:end+20])
	fmt.Print("-------\n\n")
}
