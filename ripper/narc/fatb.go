package narc

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/utils"
)

type entryFATB struct {
	Start uint32
	End   uint32
}

type FrameFATB struct {
	numEntries uint32
	Entries    []entryFATB
}

func (fatb FrameFATB) Entry(index int) entryFATB {
	return fatb.Entries[index]
}

func (fatb FrameFATB) String() string {
	ret := "--- FATB frame ---\n"
	ret += fmt.Sprintf("# entries: %d\n", fatb.numEntries)

	// show first 4 entries?
	for i := range 4 {
		start, end := fatb.Entries[i].Start, fatb.Entries[i].End
		ret += fmt.Sprintf("entry %d\n  start=0x%08X\n  end  =0x%08X\n", i, start, end)
	}

	return ret
}

func newEntryFATB(fatb []byte, start int) entryFATB {
	return entryFATB{
		Start: utils.U32(fatb, start),
		End:   utils.U32(fatb, start+4),
	}
}

/*
rom should start AFTER NFF header (offset = 0x10)
*/
func newFrameFATB(rom []byte, offset uint32) NitroFrame[FrameFATB] {
	frame := rom[offset:]
	frameSize := utils.U32(frame, 4)
	fatb := frame[8:] // includes FATB header
	numEntries := utils.U32(fatb, 0)
	fatbFrame := FrameFATB{numEntries, make([]entryFATB, numEntries)}

	for i := 0; i < int(numEntries); i += 1 {
		fatbFrame.Entries[i] = newEntryFATB(fatb[4:], i*0x8)
	}

	return NitroFrame[FrameFATB]{
		magic:     frame[:4],
		frameSize: frameSize,
		Data:      fatbFrame,
	}
}
