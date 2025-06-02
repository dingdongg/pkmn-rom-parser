package sav

import "github.com/dingdongg/pkmn-rom-parser/v7/consts/gamever"

// type ISave interface {
// 	Chunk(offset uint) Chunk
// 	Validate() error
// 	LatestData() *Chunk
// 	PartySection() []byte
// 	PartySize() uint32
// 	PartyOffset() uint
// 	Get(start uint, numBytes uint) []byte
// 	Data() []byte
// }

func NewSavBW(savefile []byte) *savBW {
	return &savBW{
		version: gamever.BW,
		data: savefile,
		smallBlockSize: 0xBB00,
		bigBlockSize: 0x18400,
		partyOffset: 0x18E00,
	}
}

// TODO: figure out the folllowing:
// - are the correspoding big-small block pairs always adjacent to one another?
//   or can they be spatially separated like DPPT/HGSS savefiles?
// * for now, assume that they are adjacent to one another
func (sav *savBW) Chunk(offset uint) Chunk {
	sbData := sav.data[0x18400+offset : 0x18400+sav.smallBlockSize+offset]
	// sbFooter := ? not sure if BW savefiles still have footers 
	// purpose of footers in DPPT+HGSS was for sb-bb linking, 
	// so if they're right next to each other (the main assumption for now),
	// footer wouldn't even be necessary

	small := NewBlock(sbData, []byte{}, 0x18400+offset)

	bbData := sav.data[0x0+offset : sav.bigBlockSize+offset]
	big := NewBlock(bbData, []byte{}, 0x0+offset)

	return Chunk{
		SmallBlock: small,
		BigBlock: big,
	}
}

func (sav *savBW) Validate() error {
	return nil
}

func (sav *savBW) LatestData() *Chunk {
	return nil
}

func (sav *savBW) PartySection() []byte {
	return []byte{}
}

func (sav *savBW) PartySize() uint32 {
	return 0
}

func (sav *savBW) PartyOffset() uint {
	return 0
}

func (sav *savBW) Get(start uint, numBytes uint) []byte {
	return []byte{}
}

func (sav *savBW) Data() []byte {
	return []byte{}
}