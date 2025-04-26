package enums

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

// String returns the string representation of a Nature
func (n Nature) String() string {
	switch n {
	case Hardy: return "Hardy"
	case Lonely: return "Lonely"
	case Brave: return "Brave"
	case Adamant: return "Adamant"
	case Naughty: return "Naughty"
	case Bold: return "Bold"
	case Docile: return "Docile"
	case Relaxed: return "Relaxed"
	case Impish: return "Impish"
	case Lax: return "Lax"
	case Timid: return "Timid"
	case Hasty: return "Hasty"
	case Serious: return "Serious"
	case Jolly: return "Jolly"
	case Naive: return "Naive"
	case Modest: return "Modest"
	case Mild: return "Mild"
	case Quiet: return "Quiet"
	case Bashful: return "Bashful"
	case Rash: return "Rash"
	case Calm: return "Calm"
	case Gentle: return "Gentle"
	case Sassy: return "Sassy"
	case Careful: return "Careful"
	case Quirky: return "Quirky"
	default: return "Unknown Nature"
	}
}

