package shuffler

import (
	"testing"

	"github.com/dingdongg/pkmn-rom-parser/v3/consts"
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