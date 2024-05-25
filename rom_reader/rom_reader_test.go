package rom_reader

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

/*

Party pokemon structure

+--------------+
| personality  | 4B   \
+--------------+      |
| unused       | 2B   |--> "metadata"
+--------------+	  |
| checksum     | 2B	  /
+--------------+
|              |
|              |
| 4 32B blocks | 128B
| (A, B, C, D) |
|              |
|              |
+--------------+
|              |
|              |
| battle stats | 100B
|              |
|              |
+--------------+
                = 236B
*/

func TestGetUnshuffledPos(t *testing.T) {
	blocks := []uint{A, B, C, D}

	for _, bo := range unshuffleTable {
		for _, b := range blocks {
			res := bo.getUnshuffledPos(b)
			expected := 0x8 + (bo.OriginalPos[b] * BLOCK_SIZE_BYTES)
			idx := (res - 0x8) / BLOCK_SIZE_BYTES

			if bo.ShuffledPos[idx] != b {
				t.Fatalf("expected 0x%x, got 0x%x\n", expected, res)
			}
		}
	}
}

func TestGetPokemonInvalidBlock(t *testing.T) {
	var invalidBlock uint = 4
	personality := uint32(5)

	_, err := getPokemonBlock(make([]byte, 1024), invalidBlock, personality)

	if err == nil {
		t.Fatalf("Error not thrown for invalid block: %d\n", invalidBlock)
	}
}

func TestGetPokemonBlock(t *testing.T) {
	blocks := []uint{A, B, C, D}
	personality := uint32(5)
	mockBuffer := make([]byte, 136)

	for _, b := range blocks {
		res, err := getPokemonBlock(mockBuffer, b, personality)
		if err != nil {
			t.Fatal("Unexpected error ", err)
		}

		if len(res) != int(BLOCK_SIZE_BYTES) {
			t.Fatalf("BUF LENGTH: expected %d but got %d\n", BLOCK_SIZE_BYTES, len(res))
		}
	}
}

func TestGetPokemon(t *testing.T) {
	// file takes care of the personality offset
	savefile, err := os.ReadFile("./mocks/mock_pokemon_data")
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	firstPokemon := GetPokemon(savefile[:], 0)

	expectedPokemon := Pokemon{
		461,
		"WEAVILE",
		BattleStat{
			58,
			Stats{163, 181, 93, 63, 106, 215},
		},
		"None",
		"Jolly",
		"Pressure",
		Stats{0, 255, 0, 0, 3, 252},
	}

	if !cmp.Equal(firstPokemon, expectedPokemon) {
		t.Fatalf("expected %+v, but got %+v\n", expectedPokemon, firstPokemon)
	}
}
