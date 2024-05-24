package prng

type PRNG struct {
	Checksum    uint16
	Personality uint32
	PrevResult  uint
}

type BattleStatPRNG struct {
	PrevResult  uint
	Personality uint32
}

func InitBattleStatPRNG(personality uint32) BattleStatPRNG {
	return BattleStatPRNG{uint(personality), personality}
}

func (bsprng *BattleStatPRNG) Next() uint16 {
	result := 0x041C64E6D*bsprng.PrevResult + 0x06073
	bsprng.PrevResult = result
	result >>= 16
	// return the upper 16 bits only for external use; internally, all bits should be held for future calls
	return uint16(result & 0xFFFF)
}

func Init(checksum uint16, personality uint32) PRNG {
	return PRNG{checksum, personality, uint(checksum)}
}

func (prng *PRNG) Next() uint16 {
	result := 0x041C64E6D*prng.PrevResult + 0x06073
	prng.PrevResult = result
	result >>= 16
	// return the upper 16 bits only for external use; internally, all bits should be held for future calls
	return uint16(result & 0xFFFF)
}
