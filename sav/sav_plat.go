package sav

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
)

func NewSavPLAT(savefile []byte) *savPLAT {
	return &savPLAT{
		version:        gamever.PLAT,
		data:           savefile,
		smallBlockSize: 0xCF2C,
		bigBlockSize:   0x121E4,
		partyOffset:    0xA0,
	}
}

func (sav *savPLAT) Chunk(offset uint) Chunk {
	sbData := sav.data[0x0+offset : sav.smallBlockSize+offset-0x14]
	sbFooter := sav.data[sav.smallBlockSize+offset-0x14 : sav.smallBlockSize+offset]
	small := NewBlock(sbData, sbFooter, 0x0+offset)

	bbData := sav.data[sav.smallBlockSize+offset : sav.smallBlockSize+offset+sav.bigBlockSize-0x14]
	bbFooter := sav.data[sav.smallBlockSize+offset+sav.bigBlockSize-0x14 : sav.smallBlockSize+offset+sav.bigBlockSize]
	big := NewBlock(bbData, bbFooter, sav.smallBlockSize+offset)

	return Chunk{
		SmallBlock: small,
		BigBlock:   big,
	}
}

func (sav *savPLAT) Validate() error {
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

func (sav *savPLAT) LatestData() *Chunk {
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

func (sav *savPLAT) PartySection() []byte {
	return sav.data[sav.partyOffset:]
}

func (sav *savPLAT) PartySize() uint32 {
	return binary.LittleEndian.Uint32(sav.data[sav.partyOffset-4 : sav.partyOffset])
}

func (sav *savPLAT) PartyOffset() uint {
	return sav.partyOffset
}

func (sav *savPLAT) Get(start uint, numBytes uint) []byte {
	return sav.data[start : start+numBytes]
}

func (sav *savPLAT) Data() []byte {
	return sav.data
}