package locator

import "github.com/dingdongg/pkmn-rom-parser/v5/validator"

// return address to the start of the latest save chunk
// REQUIREMENT: savefile must be the entire .SAV file
func GetLatestSaveChunk(savefile []byte) *validator.Chunk {
	firstChunk := validator.GetChunk(savefile, 0)
	secondChunk := validator.GetChunk(savefile, 0x40000)

	var latestSmallBlock validator.Block
	if firstChunk.SmallBlock.Footer.SaveNumber >= secondChunk.SmallBlock.Footer.SaveNumber {
		latestSmallBlock = firstChunk.SmallBlock
	} else {
		latestSmallBlock = secondChunk.SmallBlock
	}

	var latestBigBlock validator.Block
	if latestSmallBlock.Footer.Identifier == firstChunk.BigBlock.Footer.Identifier {
		latestBigBlock = firstChunk.BigBlock
	} else {
		latestBigBlock = secondChunk.BigBlock
	}

	return &validator.Chunk{
		SmallBlock: latestSmallBlock,
		BigBlock:   latestBigBlock,
	}
}
