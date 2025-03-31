package narc

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
)

type FrameFIMG struct {
	data []byte
}

func (fimg FrameFIMG) String() string {
	ret := "--- FIMG Frame ---\n"
	ret += fmt.Sprintf("filesize: 0x%08X bytes\n", len(fimg.data))
	return ret
}

func newFrameFIMG(rom []byte, offset uint32) NitroFrame[FrameFIMG] {
	frame := rom[offset:]
	frameSize := utils.U32(frame, 4)
	fimg := frame[8 : frameSize] // 8 + frameSize - 8

	return NitroFrame[FrameFIMG]{
		magic: frame[:4],
		frameSize: frameSize,
		data: FrameFIMG{ fimg },
	}
}