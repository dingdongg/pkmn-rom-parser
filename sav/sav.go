package sav

import (
	"encoding/binary"
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"
)

type ISave interface {
	GetChunk(offset uint) Chunk
	Validate() error
	GetLatestData() *Chunk
	GetPartySection() []byte
	GetPartySize() uint32
}

type gen4Savefile struct {
	version        gamever.GameVer
	data           []byte
	smallBlockSize uint
	bigBlockSize   uint
	partyOffset    uint
}

// tODO: include important offsets as fields
type savPLAT gen4Savefile
type savHGSS gen4Savefile

func NewSavPLAT(savefile []byte) *savPLAT {
	return &savPLAT{
		version:        gamever.PLAT,
		data:           savefile,
		smallBlockSize: 0xCF2C,
		bigBlockSize:   0x121E4,
		partyOffset:    0xA0,
	}
}

func (sav *savPLAT) GetChunk(offset uint) Chunk {
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
	firstChunk := sav.GetChunk(0x0)
	secondChunk := sav.GetChunk(0x40000)

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

func (sav *savPLAT) GetLatestData() *Chunk {
	chunk1 := sav.GetChunk(0x0)
	chunk2 := sav.GetChunk(0x40000)

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

func (sav *savPLAT) GetPartySection() []byte {
	return sav.data[sav.partyOffset:]
}

func (sav *savPLAT) GetPartySize() uint32 {
	return binary.LittleEndian.Uint32(sav.data[sav.partyOffset-4 : sav.partyOffset])
}

func NewSavHGSS(savefile []byte) *savHGSS {
	return &savHGSS{
		version:        gamever.HGSS,
		data:           savefile,
		smallBlockSize: 0xF628,
		bigBlockSize:   0x12310,
		partyOffset:    0x98,
	}
}

func (sav *savHGSS) GetChunk(offset uint) Chunk {
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
	firstChunk := sav.GetChunk(0x0)
	secondChunk := sav.GetChunk(0x40000)

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

func (sav *savHGSS) GetLatestData() *Chunk {
	chunk1 := sav.GetChunk(0x0)
	chunk2 := sav.GetChunk(0x40000)

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

func (sav *savHGSS) GetPartySection() []byte {
	return sav.data[sav.partyOffset:]
}

func (sav *savHGSS) GetPartySize() uint32 {
	return binary.LittleEndian.Uint32(sav.data[sav.partyOffset-4 : sav.partyOffset])
}

func Validate(savefile []byte) (ISave, error) {
	game, err := identifyGameVersion(savefile)
	if err != nil {
		return nil, err
	}

	if err = game.Validate(); err != nil {
		return nil, err
	}

	return game, nil
}
