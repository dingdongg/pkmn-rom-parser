package walker

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
)

type Walker struct {
	buffer []byte
	index int
	done bool
}

func NewWalker(buf []byte) *Walker {
	return &Walker{
		buffer: buf,
		index: 0,
		done: false,
	}
}

func (w *Walker) Seek(newIndex int) error {
	if newIndex < 0 || newIndex >= len(w.buffer) {
		return fmt.Errorf("index out of bounds")
	}

	w.index = newIndex
	return nil
}

func (w *Walker) walk(steps int) {
	w.index += steps
	if w.index >= len(w.buffer) {
		w.done = true
	}
}

func (w *Walker) Index() int {
	return w.index
}

func (w *Walker) U8() uint8 {
	ret := utils.U8(w.buffer, w.index)
	w.walk(1)
	return ret
}

func (w *Walker) U16() uint16 {
	ret := utils.U16(w.buffer, w.index)
	w.walk(2)
	return ret
}

func (w *Walker) U32() uint32 {
	ret := utils.U32(w.buffer, w.index)
	w.walk(4)
	return ret
}
