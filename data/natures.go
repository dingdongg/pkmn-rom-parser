package data

import "errors"

var natureTable [25]string = [25]string{
	"Hardy",		// attack
	"Lonely",
	"Brave",
	"Adamant",
	"Naughty",
	"Bold",			// defense
	"Docile",
	"Relaxed",
	"Impish",
	"Lax",
	"Timid",		// speed
	"Hasty",
	"Serious",
	"Jolly",
	"Naive",
	"Modest",		// sp. attack
	"Mild",
	"Quiet",
	"Bashful",
	"Rash",
	"Calm",			// sp. defense
	"Gentle",
	"Sassy",
	"Careful",
	"Quirky",
}

func GetNature(index uint) (string, error) {
	if index >= uint(len(natureTable)) {
		return "", errors.New("invalid index")
	}

	return natureTable[index], nil
}

func GenerateNatureMap() map[string][5]uint {
	natureInfo := make(map[string][5]uint)

	for i, n := range natureTable {
		boostStat := i / 5
		hinderStat := i % 5

		value := [5]uint{100, 100, 100, 100, 100}
		value[boostStat] += 10
		value[hinderStat] -= 10

		natureInfo[n] = value
	}

	return natureInfo
}