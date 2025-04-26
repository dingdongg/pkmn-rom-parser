package block

import (
	"encoding/binary"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/walker"
)


type Footer struct {
	Bridge uint32
	SaveCount uint32
	BlockSize uint32
	MagicTimestamp uint32
	unknown uint16
	Checksum uint16
}

func NewFooter(block []byte) Footer {
	// last 20 bytes of block
	footerData := block[len(block)-0x14:]
	w := walker.NewWalker(footerData)
	return Footer{
		Bridge: w.U32(),
		SaveCount: w.U32(),
		BlockSize: w.U32(),
		MagicTimestamp: w.U32(),
		unknown: w.U16(),
		Checksum: w.U16(),
	}
}

func (f Footer) Bytes() []byte {
	buf := make([]byte, 0)
	buf = binary.LittleEndian.AppendUint32(buf, f.Bridge)
	buf = binary.LittleEndian.AppendUint32(buf, f.SaveCount)
	buf = binary.LittleEndian.AppendUint32(buf, f.BlockSize)
	buf = binary.LittleEndian.AppendUint32(buf, f.MagicTimestamp)
	buf = binary.LittleEndian.AppendUint16(buf, f.unknown)
	buf = binary.LittleEndian.AppendUint16(buf, f.Checksum)

	return buf
}

type Block struct {
	data []byte
	offset int
	Footer Footer
}

func NewBlock(block []byte, offset int) *Block {
	return &Block{
		data: block[:len(block)-0x14],
		offset: offset,
		Footer: NewFooter(block),
	}
}

func (b *Block) Size() uint32 {
	return b.Footer.BlockSize
}

func (b *Block) Bytes() []byte {
	output := make([]byte, len(b.data))
	copy(output, b.data)
	output = append(output, b.Footer.Bytes()...)
	return output
}

func (b *Block) Data() []byte {
	return b.data
}

func (b *Block) Offset() int {
	return b.offset
}