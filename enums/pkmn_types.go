package enums

import "fmt"

type PokemonType uint8

const (
	NORMAL   PokemonType = iota // 0
	FIGHTING                    // 1
	FLYING                      // 2
	POISON                      // 3
	GROUND                      // 4
	ROCK                        // 5
	BUG                         // 6
	GHOST                       // 7
	STEEL                       // 8
	UNKNOWN						// 9
	FIRE                        // 10
	WATER                       // 11
	GRASS                       // 12
	ELECTRIC                    // 13
	PSYCHIC                     // 14
	ICE                         // 15
	DRAGON                      // 16
	DARK                        // 17
)

// TypeName returns the string representation of a type ID
func (pt PokemonType) String() string {
	typeNames := map[PokemonType]string{
		NORMAL:   "Normal",
		FIGHTING: "Fighting",
		FLYING:   "Flying",
		POISON:   "Poison",
		GROUND:   "Ground",
		ROCK:     "Rock",
		BUG:      "Bug",
		GHOST:    "Ghost",
		STEEL:    "Steel",
		FIRE:     "Fire",
		WATER:    "Water",
		GRASS:    "Grass",
		ELECTRIC: "Electric",
		PSYCHIC:  "Psychic",
		ICE:      "Ice",
		DRAGON:   "Dragon",
		DARK:     "Dark",
		UNKNOWN:  "???",
	}

	if name, ok := typeNames[pt]; ok {
		return name
	}
	
	return fmt.Sprintf("Invalid(%d)", pt)
}