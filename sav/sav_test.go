package sav

import (
	"fmt"
	"os"
	"testing"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts"
	"github.com/dingdongg/pkmn-rom-parser/v7/tutil"
)

var templs = tutil.GetTemplates()

func TestPlatSAV(t *testing.T) {
	f, err := os.ReadFile("./../savefiles/new.sav")

	if err != nil {
		t.Fatal("error opening platinum savefile")
	}

	var savefile ISave = NewSavPLAT(f)

	first := savefile.GetChunk(0x0)
	second := savefile.GetChunk(0x40000)

	fmt.Println(first)
	fmt.Println(second)
}

func TestPlatHGSS(t *testing.T) {
	f, err := os.ReadFile("./../savefiles/soulsilver.sav")

	if err != nil {
		t.Fatal("error opening HGSS savefile")
	}

	var savefile ISave = NewSavHGSS(f)

	first := savefile.GetChunk(0x0)
	second := savefile.GetChunk(0x40000)

	fmt.Println(first)
	fmt.Println(second)
}

func TestArbitrarySavefile(t *testing.T) {
	f, err := os.ReadFile("./../savefiles/soulsilver.sav")

	if err != nil {
		t.Fatal("error opening HGSS savefile")
	}

	// incorrectly parse HGSS as a PLAT savefile
	var savefile ISave = NewSavPLAT(f)

	chunk1 := savefile.GetChunk(0x0)
	chunk2 := savefile.GetChunk(0x40000)

	if chunk1.SmallBlock.Footer.K == consts.MAGIC_TIMESTAMP_JP_INTL {
		t.Fatalf(
			"shouldn't equal 0x%x but got 0x%x\n",
			consts.MAGIC_TIMESTAMP_JP_INTL,
			chunk1.SmallBlock.Footer.K,
		)
	}

	if chunk2.SmallBlock.Footer.K == consts.MAGIC_TIMESTAMP_JP_INTL {
		t.Fatalf(
			"shouldn't equal 0x%x but got 0x%x\n",
			consts.MAGIC_TIMESTAMP_JP_INTL,
			chunk2.SmallBlock.Footer.K,
		)
	}
}
