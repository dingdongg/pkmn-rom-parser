package validator

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/validator/block"
)

func validateBlock(b *block.Block) error {
	timestamp := b.Footer.MagicTimestamp
	if timestamp != enums.MAGIC_TS_JP_INTL && timestamp != enums.MAGIC_TS_KR {
		return fmt.Errorf("magic number invalid")
	}

	// checksum validations
	expected, actual := b.Footer.Checksum, crypt.CRC16_CCITT(b.Data())
	if expected != actual {
		msg := "checksum mismatch: 0x%04X (expected), 0x%04X (actual)"
		return fmt.Errorf(msg, expected, actual)
	}

	return nil
}

/*
	this functino may seem like it is returning a region of memory on the stack
	(ie. a static array), but golang actually returns a copy of this static array
	which is allocated within the stack frame of the CALLING function
*/
func getBlocks(savefile []byte, start, end uint) [2]*block.Block {
	sb1 := block.NewBlock(savefile[start : end+1], 0x0)
	start, end = start+0x40000, end+0x40000
	sb2 := block.NewBlock(savefile[start : end+1], 0x40000)

	return [2]*block.Block{ sb1, sb2 }
}

func LatestSmallBlock(savefile []byte, offsets utils.Range[uint]) (*block.Block, error) {
	blocks := getBlocks(savefile, offsets.Start, offsets.End)
	count1, count2 := blocks[0].Footer.SaveCount, blocks[1].Footer.SaveCount
	if count1 < count2 {
		// swap blocks
		blocks[0], blocks[1] = blocks[1], blocks[0]
	}

	var err error
	for _, b := range blocks {
		err = validateBlock(b)
		if err == nil {
			return b, nil
		}
	}

	// this should be very unlikely (savefile is just giga fucked gg)
	return nil, err
}

func LatestBigBlock(savefile []byte, offsets utils.Range[uint], bridgeNum uint32) (*block.Block, error) {
	blocks := getBlocks(savefile, offsets.Start, offsets.End)
	count1, count2 := blocks[0].Footer.SaveCount, blocks[1].Footer.SaveCount
	if count2 == bridgeNum {
		blocks[0], blocks[1] = blocks[1], blocks[0]
	} else if count1 != bridgeNum {
		return nil, fmt.Errorf("no big block match found")
	}

	var err error
	for _, b := range blocks {
		err = validateBlock(b)
		if err == nil {
			return b, nil
		}
	}

	return nil, err
}

func validateGen4(buf []byte, sbRange utils.Range[uint], bbRange utils.Range[uint]) error {
	// the LatestXXBlock() helpers validate the fetched blocks
	latestSb, err := LatestSmallBlock(buf, sbRange)
	if err != nil {
		return err
	}

	_, err = LatestBigBlock(buf, bbRange, latestSb.Footer.Bridge)
	if err != nil {
		return err
	}

	fmt.Println("Validation successful")
	return nil
}