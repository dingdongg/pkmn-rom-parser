package narc

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/utils"
)

const FILENAMES_INCLUDED uint32 = 0x00_00_00_08

type entryFNTB struct {
	length uint8
	name   []byte // of size `length`
}

type FrameFNTB struct {
	filenames uint32
	unknown   uint32
	entries   []entryFNTB
}

func (fntb FrameFNTB) String() string {
	ret := "--- FNTB Frame ---\n"

	if len(fntb.entries) == 0 {
		ret += "  empty.\n"
		return ret
	}

	for i := 0; i < 4; i++ {
		length, name := fntb.entries[i].length, fntb.entries[i].name
		ret += fmt.Sprintf("entry %d\n  filename: '%s'\n", i, name[:length])
	}

	return ret
}

func newEntryFNTB(fntb []byte, start int) entryFNTB {
	length := utils.U8(fntb, start)
	buffer := fntb[start+1 : start+1+int(length)]

	return entryFNTB{
		length: length,
		name:   buffer,
	}
}

func newFrameFNTB(rom []byte, offset uint32, numFiles uint32) NitroFrame[FrameFNTB] {
	frame := rom[offset:]
	frameSize := utils.U32(frame, 4)
	fntb := frame[8:]

	fntbFrame := FrameFNTB{
		filenames: utils.U32(fntb, 0),
		unknown:   utils.U32(fntb, 4),
		entries:   make([]entryFNTB, 0),
	}

	if fntbFrame.filenames == FILENAMES_INCLUDED {
		// read entries
		prevEntrySize := 0
		for range int(numFiles) {
			entry := newEntryFNTB(fntb[8:], prevEntrySize)
			fntbFrame.entries = append(fntbFrame.entries, entry)
			prevEntrySize = int(entry.length) + 1
		}
	}

	return NitroFrame[FrameFNTB]{
		magic:     frame[:4],
		frameSize: frameSize,
		Data:      fntbFrame,
	}
}
