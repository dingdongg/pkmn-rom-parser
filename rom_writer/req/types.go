package req

type Writable interface {
	Bytes() ([]byte, error)
}

// for battle stats. implements Writable
type WriteStats struct {
	Hp        uint
	Attack    uint
	Defense   uint
	SpAttack  uint
	SpDefense uint
	Speed     uint
}

// implements Writable
type WriteEffortValue WriteStats

// implements Writable

type WriteIndivValue WriteStats

// for IDs/level. implements Writable
type WriteUint struct {
	Val uint
	NumBytes uint	// number of bytes. Used in Bytes() implementation
}

// for nicknames. implements Writable
type WriteString struct {
	Val string
}

type NewData map[string]Writable

type WriteRequest struct {
	PartyIndex uint
	Contents   NewData
}