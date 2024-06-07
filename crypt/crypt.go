package crypt

import (
	"encoding/binary"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v3/prng"
)

func DecryptPokemon(ciphertext []byte) []byte {
	personality := binary.LittleEndian.Uint32(ciphertext[0 : 4])
	checksum := binary.LittleEndian.Uint16(ciphertext[6 : 8])

	rand := prng.Init(checksum, personality)

	buffer := ciphertext[:8]
	plaintextSum := uint16(0)

	for i := 0x8; i < 0x87; i += 2 {
		word := binary.LittleEndian.Uint16(ciphertext[i : i+2])
		plaintext := word ^ rand.Next()
		plaintextSum += plaintext
		littleByte := byte(plaintext & 0xFFFF)
		bigByte := byte((plaintext >> 8) & 0xFFFF)
		buffer = append(buffer, littleByte, bigByte)
	}

	if plaintextSum != checksum {
		log.Fatalf("Checksum invalid. expected 0x%x, got 0x%x\n", checksum, plaintextSum)
	}

	return buffer
}

// first block of ciphertext points to offset 0x88 in a whole party pokemon block
func DecryptBattleStats(ciphertext []byte, personality uint32) []byte {
	bsprng := prng.InitBattleStatPRNG(personality)
	var plaintext []byte

	for i := 0; i < 0x14; i += 2 {
		decrypted := bsprng.Next() ^ binary.LittleEndian.Uint16(ciphertext[i:i+2])
		plaintext = append(plaintext, byte(decrypted&0xFF), byte((decrypted>>8)&0xFF))
	}

	return plaintext
}