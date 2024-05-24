package prng

import (
	"testing"
)

func TestInit(t *testing.T) {
	checksum := uint16(0)
	personality := uint32(5)
	actual := Init(checksum, personality)

	if actual.Checksum != checksum {
		t.Fatalf("Invalid checksum; expected 0x%x, got 0x%x", checksum, actual.Checksum)
	}

	if actual.Personality != personality {
		t.Fatalf("Invalid personality; expected 0x%x, got 0x%x", personality, actual.Personality)
	}

	if actual.PrevResult != uint(checksum) {
		t.Fatalf("Invalid prevResult; expected 0x%x, got 0x%x", uint(checksum), actual.PrevResult)
	}
}

func TestNextRetVal(t *testing.T) {
	checksum := uint16(0)
	personality := uint32(5)
	prng := Init(checksum, personality)

	ret := prng.Next()
	expected := uint16((0x6073 >> 16) & 0xFFFF)
	if ret != expected {
		t.Fatalf("expected 0x%x, got 0x%x", expected, ret)
	}
}

func TestNextInternals(t *testing.T) {
	checksum := uint16(0)
	personality := uint32(5)
	prng := Init(checksum, personality)

	numCalls := 3
	expectedValues := []uint{0x6073, 0x18C7E97E7B6A, 0xF47F2B6C52713895}

	for i := 0; i < numCalls; i++ {
		prng.Next()
		if prng.PrevResult != expectedValues[i] {
			t.Fatalf("expected 0x%x, got 0x%x", expectedValues[i], prng.PrevResult)
		}
	}
}

func TestBattleStatInit(t *testing.T) {
	personality := uint32(5)
	actual := InitBattleStatPRNG(personality)

	if actual.Personality != personality {
		t.Fatalf("Invalid personality; expected 0x%x, got 0x%x", personality, actual.Personality)
	}

	if actual.PrevResult != uint(personality) {
		t.Fatalf("Invalid prevResult; expected 0x%x, got 0x%x", uint(personality), actual.PrevResult)
	}
}

func TestBSNextVal(t *testing.T) {
	personality := uint32(4)
	prng := InitBattleStatPRNG(personality)

	ret := prng.Next()
	expected := uint16((((0x041C64E6D << 2) + 0x6073) >> 16) & 0xFFFF)
	if ret != expected {
		t.Fatalf("expected 0x%x, got 0x%x", expected, ret)
	}
}

func TestBSNextInternals(t *testing.T) {
	personality := uint32(0)
	prng := InitBattleStatPRNG(personality)

	numCalls := 3
	expectedValues := []uint{0x6073, 0x18C7E97E7B6A, 0xF47F2B6C52713895}

	for i := 0; i < numCalls; i++ {
		prng.Next()
		if prng.PrevResult != expectedValues[i] {
			t.Fatalf("expected 0x%x, got 0x%x", expectedValues[i], prng.PrevResult)
		}
	}
}
