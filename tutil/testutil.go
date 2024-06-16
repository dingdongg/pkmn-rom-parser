package tutil

// useful test helper functions for logging

import (
	"testing"
)

func ErrMsg(t *testing.T, message string) {
	t.Fatal(message)
}

func ErrUint(t *testing.T, expected uint, actual uint, inHex bool) {
	s := "expected %d but got %d\n"
	if inHex {
		s = "expected 0x%x but got 0x%x\n"
	}

	t.Fatalf(s, expected, actual)
}

func ErrInt(t *testing.T, expected int, actual int, inHex bool) {
	s := "expected %d but got %d\n"
	if inHex {
		s = "expected 0x%x but got 0x%x\n"
	}

	t.Fatalf(s, expected, actual)
}

func ErrString(t *testing.T, expected string, actual string) {
	t.Fatalf("expected '%s' but got '%s'\n", expected, actual)
}