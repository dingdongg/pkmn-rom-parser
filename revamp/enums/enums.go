package enums

const (
	MAGIC_TS_JP_INTL = 0x20060623
	MAGIC_TS_KR = 0x20070903
)

type Gender int
const (
	Male Gender = iota
	Female 
	Unknown
)

func (g Gender) String() string {
	if g == Male {
		return "Male"
	} else if g == Female {
		return "Female"
	} else {
		return "Unknown"
	}
}

type GameVersion int

const (
	DP GameVersion = iota
	PLAT
	HGSS
	BW
	B2W2
	INVALID
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
