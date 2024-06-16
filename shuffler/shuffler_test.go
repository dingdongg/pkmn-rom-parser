package shuffler

import (
	"testing"

	"github.com/dingdongg/pkmn-rom-parser/v6/consts"
)

func TestGetUnshuffledPos(t *testing.T) {
	blocks := []uint{A, B, C, D}

	for _, bo := range unshuffleTable {
		for _, b := range blocks {
			res := bo.GetUnshuffledPos(b)
			expected := 0x8 + (bo.OriginalPos[b] * consts.BLOCK_SIZE_BYTES)
			idx := (res - 0x8) / consts.BLOCK_SIZE_BYTES

			if bo.ShuffledPos[idx] != b {
				t.Fatalf("expected 0x%x, got 0x%x\n", expected, res)
			}
		}
	}
}

func TestGetPokemonInvalidBlock(t *testing.T) {
	var invalidBlock uint = 4
	personality := uint32(5)

	_, err := GetPokemonBlock(make([]byte, 1024), invalidBlock, personality)

	if err == nil {
		t.Fatalf("Error not thrown for invalid block: %d\n", invalidBlock)
	}
}

func TestGetPokemonBlock(t *testing.T) {
	blocks := []uint{A, B, C, D}
	personality := uint32(5)
	mockBuffer := make([]byte, 136)

	for _, b := range blocks {
		res, err := GetPokemonBlock(mockBuffer, b, personality)
		if err != nil {
			t.Fatal("Unexpected error ", err)
		}

		if len(res) != int(consts.BLOCK_SIZE_BYTES) {
			t.Fatalf("BUF LENGTH: expected %d but got %d\n", consts.BLOCK_SIZE_BYTES, len(res))
		}
	}
}
