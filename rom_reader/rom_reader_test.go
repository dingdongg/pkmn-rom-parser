package rom_reader

import (
	"os"
	"testing"

	"github.com/dingdongg/pkmn-rom-parser/v3/consts"
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

		if len(res) != int(consts.BLOCK_SIZE_BYTES) {
			t.Fatalf("BUF LENGTH: expected %d but got %d\n", consts.BLOCK_SIZE_BYTES, len(res))
		}
	}
}

func TestParsePokemon(t *testing.T) {
	// file takes care of the personality offset
	savefile, err := os.ReadFile("./mocks/mock_pokemon_data")
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	firstPokemon := parsePokemon(savefile[:], 0)

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
		Stats{25, 1, 23, 25, 5, 17},
	}

	if !cmp.Equal(firstPokemon, expectedPokemon) {
		t.Fatalf("expected %+v, but got %+v\n", expectedPokemon, firstPokemon)
	}
}
