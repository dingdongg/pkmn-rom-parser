package req

type Writable interface {
	Bytes() ([]byte, error)
}

// for battle stats
type WriteStats struct {
	Hp        uint
	Attack    uint
	Defense   uint
	SpAttack  uint
	SpDefense uint
	Speed     uint
}

// for IDs/level/IVs/EVs
type WriteUint struct {
	Val uint
	NumItems uint	// for compressed values, the number of items compressed
	ItemBits uint	// for compressed values, number of bits to represent each item
}

// for nicknames
type WriteString struct {
	Val string
}

type NewData map[string]Writable

// types that implement CompressibleStat can have their stat values
// compressed into an unsigned integer. For instance,
// each IV stat uses 5 bits (=30) - so IVs can be packed in a uint32.
// each EV stat uses 8 bits (=48) - so EVs can be packed in a uint64.
type CompressibleStat interface {
	Compress(elemBits uint) uint
}

type WriteRequest struct {
	PartyIndex uint
	Contents   NewData
}