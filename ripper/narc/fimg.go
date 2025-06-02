package narc

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/utils"
)

type FrameFIMG struct {
	Data []byte
}

func (fimg FrameFIMG) String() string {
	ret := "--- FIMG Frame ---\n"
	ret += fmt.Sprintf("filesize: 0x%08X bytes\n", len(fimg.Data))
	return ret
}

func newFrameFIMG(rom []byte, offset uint32) NitroFrame[FrameFIMG] {
	frame := rom[offset:]
	frameSize := utils.U32(frame, 4)
	fimg := frame[8:frameSize] // 8 + frameSize - 8

	return NitroFrame[FrameFIMG]{
		magic:     frame[:4],
		frameSize: frameSize,
		Data:      FrameFIMG{fimg},
	}
}
