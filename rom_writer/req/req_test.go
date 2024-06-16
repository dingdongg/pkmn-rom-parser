package req

import (
	"encoding/binary"
	"testing"

	"github.com/dingdongg/pkmn-rom-parser/v5/tutil"
)

var templates = tutil.GetTemplates()

func TestNewWriteRequest(t *testing.T) {
	wr := NewWriteRequest(0)

	if wr.PartyIndex != 0 {
		t.Fatalf(templates.Uint, 0, wr.PartyIndex)
	}

	if len(wr.Contents) != 0 {
		t.Fatalf(templates.Int, 0, len(wr.Contents))
	}
}

func TestWriteItem(t *testing.T) {
	wr := NewWriteRequest(0)

	wr.WriteItem("Master Ball") // ID is 1

	res, ok := wr.Contents[ITEM]
	if !ok {
		t.Fatal("expected map entry to exist, but got DNE")
	}

	t.Log(res)

	byteForm, err := res.Bytes()
	if err != nil {
		t.Fatal("byte conversion failed")
	}

	t.Log(byteForm)

	// item ID should take up 2 bytes only
	if len(byteForm) != 2 {
		t.Fatalf(templates.Int, 2, len(byteForm))
	}

	id := uint16(byteForm[0]) | (uint16(byteForm[1]) << 8)
	if id != 0x0001 {
		t.Fatalf(templates.UintHex, 0x0001, id)
	}
}

func TestWriteAbility(t *testing.T) {
	wr := NewWriteRequest(0)

	wr.WriteAbility("Levitate") // ID is 26

	res, ok := wr.Contents[ABILITY]
	if !ok {
		t.Fatal("expected map entry to exist, but got DNE")
	}

	byteForm, err := res.Bytes()
	if err != nil {
		t.Fatal("byte conversion failed")
	}

	if len(byteForm) != 1 {
		t.Fatalf(templates.Int, 1, len(byteForm))
	}

	if byteForm[0] != 26 {
		t.Fatal(templates.UintHex, 26, byteForm[0])
	}
}

func TestWriteBattleStats(t *testing.T) {
	wr := NewWriteRequest(0)
	stats := [6]uint{65535, 0, 124, 7000, 333, 255}
	wr.WriteBattleStats(stats[0], stats[1], stats[2], stats[4], stats[5], stats[3])
	
	res, ok := wr.Contents[BATTLE_STATS]

	if !ok {
		t.Fatal("expected map entry to exeist, but got DNE")
	}

	byteForm, err := res.Bytes()
	if err != nil {
		t.Fatal("byte conversion failed")
	}

	if len(byteForm) != 12 {
		t.Fatalf(templates.Int, 12, len(byteForm))
	}

	for i := 0; i < len(byteForm); i += 2 {
		t.Logf("index %d\n", i)
		expected := stats[i / 2]
		actual := binary.LittleEndian.Uint16(byteForm[i : i+2])
		if actual != uint16(expected) {
			t.Fatalf(templates.UintHex, expected, actual)
		}
	}
}	