package req

import (
	"encoding/binary"
	"testing"

	"github.com/dingdongg/pkmn-rom-parser/v5/char"
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

func TestWriteEV(t *testing.T) {
	wr := NewWriteRequest(0)
	stats := [6]uint{0, 0, 0, 252, 252, 6}
	wr.WriteEV(stats[0], stats[1], stats[2], stats[4], stats[5], stats[3])

	res, ok := wr.Contents[EV]

	if !ok {
		t.Fatal("expected map entry to exeist, but got DNE")
	}

	byteForm, err := res.Bytes()
	if err != nil {
		t.Fatal("byte conversion failed")
	}

	if len(byteForm) != 6 {
		t.Fatalf(templates.Int, 6, len(byteForm))
	}

	for i := 0; i < len(byteForm); i++ {
		t.Logf("index %d\n", i)
		expected := stats[i]
		actual := byteForm[i]
		if actual != byte(expected) {
			t.Fatalf(templates.UintHex, expected, actual)
		}
	}
}

func TestWriteIV(t *testing.T) {
	wr := NewWriteRequest(0)
	stats := [6]uint{0, 31, 14, 5, 21, 30}
	wr.WriteIV(stats[0], stats[1], stats[2], stats[4], stats[5], stats[3])

	res, ok := wr.Contents[IV]

	if !ok {
		t.Fatal("expected map entry to exeist, but got DNE")
	}

	byteForm, err := res.Bytes()
	if err != nil {
		t.Fatal("byte conversion failed")
	}

	if len(byteForm) != 4 {
		t.Fatalf(templates.Int, 4, len(byteForm))
	}

	mask := uint32(0b11111)
	ivs := binary.LittleEndian.Uint32(byteForm)
	for i := 0; i < 6; i++ {
		val := (ivs >> (i * 5)) & mask
		t.Log(val, stats[i])
		if val != uint32(stats[i]) {
			t.Fatalf(templates.Uint, stats[i], val)
		}
	}
}

func TestWriteNicknameTooLong(t *testing.T) {
	wr := NewWriteRequest(0)
	name := "way too long of a anme";
	wr.WriteNickname(name)

	res, ok := wr.Contents[NICKNAME]

	if !ok {
		t.Fatal("expected map entry to exist, but got DNE")
	}

	if _, err := res.Bytes(); err == nil {
		t.Fatal("expected res.Bytes() to fail")
	}
}

func TestWriteNickname(t *testing.T) {
	wr := NewWriteRequest(0)
	name := "trainer"
	wr.WriteNickname(name)

	res, ok := wr.Contents[NICKNAME]

	if !ok {
		t.Fatal("expected map entry to exeist, but got DNE")
	}

	byteForm, err := res.Bytes()
	if err != nil {
		t.Fatal("byte conversion failed")
	}

	if len(byteForm) != 22 {
		t.Fatalf(templates.Int, 22, len(byteForm))
	}

	for i, r := range name {
		expected, _ := char.Index(string(r))
		actual := binary.LittleEndian.Uint16(byteForm[i*2 : (i*2)+2])
		t.Logf(templates.UintHex, expected, actual)
		if expected != actual {
			t.Fatalf(templates.UintHex, expected, actual)
		}
	}

	nullTerminator := binary.LittleEndian.Uint16(byteForm[len(name)*2  : len(name)*2+2])
	if nullTerminator != 65535 {
		t.Fatalf(templates.UintHex, 65535, nullTerminator)
	}

	zerosIndex := len(name) * 2 + 2
	for i := zerosIndex; i < 22; i++ {
		if byteForm[i] != 0 {
			t.Fatalf(templates.Uint, 0, byteForm[i])
		}
	}
}