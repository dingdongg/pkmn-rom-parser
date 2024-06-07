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
