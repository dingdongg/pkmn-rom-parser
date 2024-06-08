package crypt

import (
	"encoding/binary"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v4/prng"
)

// Computes a checksum via the CRC16-CCITT algorithm on the given data
func CRC16_CCITT(data []byte) uint16 {
	sum := uint(0xFFFF)

	for _, b := range data {
		sum = (sum << 8) ^ seeds[b^byte((sum>>8))]
	}

	return uint16(sum)
}

// Encrypts the given pokemon. Checksum will be updated as part of encryption
func EncryptPokemon(plaintext []byte) []byte {
	personality := binary.LittleEndian.Uint32(plaintext[0:4])
	// do i use the previous checksum? or the new plaintextSum calculated below?
	// UPDATE: I think im supposed to use the new checksum
	// checksum := binary.LittleEndian.Uint16(plaintext[6 : 8])

	buffer := make([]byte, 8)
	copy(buffer, plaintext[:8])

	plaintextSum := uint16(0)

	for i := 0x8; i < 0x87; i += 2 {
		word := binary.LittleEndian.Uint16(plaintext[i : i+2])
		plaintextSum += word
	}

	rand := prng.Init(plaintextSum, personality)

	for i := 0x8; i < 0x87; i += 2 {
		word := binary.LittleEndian.Uint16(plaintext[i : i+2])
		encrypted := word ^ rand.Next()
		little := byte(encrypted & 0xFF)
		big := byte((encrypted >> 8) & 0xFF)
		buffer = append(buffer, little, big)
	}

	binary.LittleEndian.PutUint16(buffer[6:8], plaintextSum)
	return append(buffer, EncryptBattleStats(plaintext[0x88:], personality)...)
}

func EncryptBattleStats(plaintext []byte, personality uint32) []byte {
	bsprng := prng.InitBattleStatPRNG(personality)
	var buffer []byte

	for i := 0; i < 0x64; i += 2 {
		decrypted := bsprng.Next() ^ binary.LittleEndian.Uint16(plaintext[i:i+2])
		buffer = append(buffer, byte(decrypted&0xFF), byte((decrypted>>8)&0xFF))
	}

	return buffer
}

func DecryptPokemon(ciphertext []byte) []byte {
	personality := binary.LittleEndian.Uint32(ciphertext[0:4])
	checksum := binary.LittleEndian.Uint16(ciphertext[6:8])

	rand := prng.Init(checksum, personality)

	buffer := make([]byte, 8)
	copy(buffer, ciphertext[:8])
	plaintextSum := uint16(0)

	for i := 0x8; i < 0x87; i += 2 {
		word := binary.LittleEndian.Uint16(ciphertext[i : i+2])
		plaintext := word ^ rand.Next()
		plaintextSum += plaintext
		littleByte := byte(plaintext & 0xFF)
		bigByte := byte((plaintext >> 8) & 0xFF)
		buffer = append(buffer, littleByte, bigByte)
	}

	if plaintextSum != checksum {
		log.Fatalf("Checksum invalid. expected 0x%x, got 0x%x\n", checksum, plaintextSum)
	}

	return append(buffer, DecryptBattleStats(ciphertext[0x88:], personality)...)
}

// first block of ciphertext points to offset 0x88 in a whole party pokemon block
func DecryptBattleStats(ciphertext []byte, personality uint32) []byte {
	bsprng := prng.InitBattleStatPRNG(personality)
	var plaintext []byte

	for i := 0; i < 0x64; i += 2 {
		decrypted := bsprng.Next() ^ binary.LittleEndian.Uint16(ciphertext[i:i+2])
		plaintext = append(plaintext, byte(decrypted&0xFF), byte((decrypted>>8)&0xFF))
	}

	return plaintext
}
