package req

import (
	"testing"

	"github.com/dingdongg/pkmn-rom-parser/v5/tutil"
)

func TestNewWriteRequest(t *testing.T) {
	wr := NewWriteRequest(0)

	if wr.PartyIndex != 0 {
		tutil.ErrUint(t, 0, wr.PartyIndex, false)
	}

	if len(wr.Contents) != 0 {
		tutil.ErrInt(t, 0, len(wr.Contents), false)
	}
}

func TestWriteItem(t *testing.T) {
	wr := NewWriteRequest(0)

	wr.WriteItem("Master Ball") // ID is 1

	res, ok := wr.Contents[ITEM]
	if !ok {
		tutil.ErrMsg(t, "expected map entry to exist, but got DNE")
	}

	byteForm, err := res.Bytes()
	if err != nil {
		tutil.ErrMsg(t, "byte conversion failed")
	}

	// item ID should take up 2 bytes only
	if len(byteForm) != 2 {
		tutil.ErrInt(t, 2, len(byteForm), false)
	}

	id := uint16(byteForm[0]) | (uint16(byteForm[1]) << 8)
	if id != 0x0001 {
		tutil.ErrUint(t, 0x0001, uint(id), true)
	}
}

func TestWriteAbility(t *testing.T) {
	wr := NewWriteRequest(0)

	wr.WriteAbility("Levitate") // ID is 26

	res, ok := wr.Contents[ABILITY]
	if !ok {
		tutil.ErrMsg(t, "expected map entry to exist, but got DNE")
	}

	byteForm, err := res.Bytes()
	if err != nil {
		tutil.ErrMsg(t, "byte conversion failed")
	}

	if len(byteForm) != 1 {
		tutil.ErrInt(t, 1, len(byteForm), false)
	}

	if byteForm[0] != 26 {
		tutil.ErrUint(t, 26, uint(byteForm[0]), true)
	}
}