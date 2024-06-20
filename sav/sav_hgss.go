package sav

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
)

func NewSavHGSS(savefile []byte) *savHGSS {
	return &savHGSS{
		version:        gamever.HGSS,
		data:           savefile,
		smallBlockSize: 0xF628,
		bigBlockSize:   0x12310,
		partyOffset:    0x98,
	}
}

func (sav *savHGSS) Chunk(offset uint) Chunk {
	sbData := sav.data[0x0+offset : sav.smallBlockSize+offset-0x10]
	sbFooter := sav.data[sav.smallBlockSize+offset-0x14 : sav.smallBlockSize+offset]
	small := NewBlock(sbData, sbFooter, 0x0+offset)

	padding := uint(0xD8)
	bbStart := sav.smallBlockSize + padding
	bbData := sav.data[bbStart+offset : bbStart+offset+sav.bigBlockSize-0x10]
	bbFooter := sav.data[bbStart+offset+sav.bigBlockSize-0x14 : bbStart+offset+sav.bigBlockSize]
	big := NewBlock(bbData, bbFooter, bbStart+offset)

	return Chunk{
		SmallBlock: small,
		BigBlock:   big,
	}
}

func (sav *savHGSS) Validate() error {
	firstChunk := sav.Chunk(0x0)
	secondChunk := sav.Chunk(0x40000)

	if !firstChunk.IsValid() {
		fmt.Println("First chunk invalid")
		return fmt.Errorf("invalid savefile")
	}

	if !secondChunk.IsValid() {
		fmt.Println("Second chunk invalid")
		return fmt.Errorf("invalid savefile")
	}

	return nil
}

func (sav *savHGSS) LatestData() *Chunk {
	chunk1 := sav.Chunk(0x0)
	chunk2 := sav.Chunk(0x40000)

	var latestSmallBlock Block
	if chunk1.SmallBlock.Footer.SaveNumber >= chunk2.SmallBlock.Footer.SaveNumber {
		latestSmallBlock = chunk1.SmallBlock
	} else {
		latestSmallBlock = chunk2.SmallBlock
	}

	var latestBigBlock Block
	if latestSmallBlock.Footer.Identifier == chunk1.BigBlock.Footer.Identifier {
		latestBigBlock = chunk1.BigBlock
	} else {
		latestBigBlock = chunk2.BigBlock
	}

	return &Chunk{
		SmallBlock: latestSmallBlock,
		BigBlock:   latestBigBlock,
	}
}

func (sav *savHGSS) PartySection() []byte {
	latest := sav.LatestData()
	return sav.data[latest.SmallBlock.Address+sav.partyOffset:]
}

func (sav *savHGSS) PartySize() uint32 {
	latest := sav.LatestData()
	offset := latest.SmallBlock.Address + sav.partyOffset
	return binary.LittleEndian.Uint32(sav.data[offset-4 : offset])
}

func (sav *savHGSS) PartyOffset() uint {
	return sav.partyOffset
}

func (sav *savHGSS) Get(start uint, numBytes uint) []byte {
	return sav.data[start : start+numBytes]
}

func (sav *savHGSS) Data() []byte {
	return sav.data
}