package sav

// this file is ignored from version control since it makes use of gitignored savefiles

import (
	"fmt"
	"os"
	"testing"
)

func TestPlatSAV(t *testing.T) {
	f, err := os.ReadFile("./../savefiles/new.sav")

	if err != nil {
		t.Fatal("error opening platinum savefile")
	}

	var savefile Savefile = NewSavPLAT(f)

	first := savefile.Chunk(0x0)
	second := savefile.Chunk(0x40000)

	fmt.Println(first)
	fmt.Println(second)
}

func TestPlatHGSS(t *testing.T) {
	f, err := os.ReadFile("./../savefiles/soulsilver.sav")

	if err != nil {
		t.Fatal("error opening HGSS savefile")
	}

	var savefile Savefile = NewSavHGSS(f)

	first := savefile.Chunk(0x0)
	second := savefile.Chunk(0x40000)

	fmt.Println(first)
	fmt.Println(second)
}

func TestArbitrarySavefile(t *testing.T) {
	f, err := os.ReadFile("./../savefiles/soulsilver.sav")

	if err != nil {
		t.Fatal("error opening HGSS savefile")
	}

	// incorrectly parse HGSS as a PLAT savefile
	var savefile Savefile = NewSavPLAT(f)

	chunk1 := savefile.Chunk(0x0)
	chunk2 := savefile.Chunk(0x40000)

	if chunk1.SmallBlock.Footer.K == MAGIC_TIMESTAMP_JP_INTL {
		t.Fatalf(
			"shouldn't equal 0x%x but got 0x%x\n",
			MAGIC_TIMESTAMP_JP_INTL,
			chunk1.SmallBlock.Footer.K,
		)
	}

	if chunk2.SmallBlock.Footer.K == MAGIC_TIMESTAMP_JP_INTL {
		t.Fatalf(
			"shouldn't equal 0x%x but got 0x%x\n",
			MAGIC_TIMESTAMP_JP_INTL,
			chunk2.SmallBlock.Footer.K,
		)
	}
}
