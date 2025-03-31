package narc

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
)

type entryFATB struct {
	start uint32
	end   uint32
}

type FrameFATB struct {
	numEntries uint32
	entries    []entryFATB
}

func (fatb FrameFATB) String() string {
	ret := "--- FATB frame ---\n"
	ret += fmt.Sprintf("# entries: %d\n", fatb.numEntries)

	// show first 4 entries?
	for i := range 4 {
		start, end := fatb.entries[i].start, fatb.entries[i].end
		ret += fmt.Sprintf("entry %d\n  start=0x%08X\n  end  =0x%08X\n", i, start, end)
	}

	return ret
}

func newEntryFATB(fatb []byte, start int) entryFATB {
	return entryFATB{
		start: utils.U32(fatb, start),
		end:   utils.U32(fatb, start+4),
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
		fatbFrame.entries[i] = newEntryFATB(fatb[4:], i*0x8)
	}

	return NitroFrame[FrameFATB]{
		magic:     frame[:4],
		frameSize: frameSize,
		Data:      fatbFrame,
	}
}
