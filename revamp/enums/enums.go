package enums

type GameVersion int

const (
	DP GameVersion = iota
	PLAT
	HGSS
	BW
	B2W2
)

func (gv GameVersion) String() string {
	if gv == DP {
		return "DP"
	} else if gv == PLAT {
		return "PLAT"
	} else if gv == HGSS {
		return "HGSS"
	} else if gv == BW {
		return "BW"
	} else if gv == B2W2 {
		return "B2W2"
	} else {
		return "unrecognized"
	}
}

type Nature uint8

const (
	Hardy Nature = iota // inc. Attack
	Lonely
	Brave
	Adamant
	Naughty
	Bold // inc. Defense
	Docile
	Relaxed
	Impish
	Lax
	Timid // inc. Speed
	Hasty
	Serious
	Jolly
	Naive
	Modest // inc. SpA
	Mild
	Quiet
	Bashful
	Rash
	Calm // inc. SpD
	Gentle
	Sassy
	Careful
	Quirky
)