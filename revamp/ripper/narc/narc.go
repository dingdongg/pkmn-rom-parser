package narc

import (
	"fmt"

	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/utils"
)

/*

things to rip out of ROM files
- pokemon data (gen 4 + gen 5)
	- id
	- name
	- types
	- gender threshold
	- base stats
	- growth type
- pokemon evolution mechanisms (gen 4 + gen 5)
- pokemon moves (id, name, type, PP?)

- table 1: Pokemons
	- PK: pokedex id
	- name
	- type1
	- type2 (optional)
	- gender threshold
	- base stats
		- hp, atk, def, spa, spd, spe
	- growth type

- table 2: Evolutions
	- PK: id/uuid
	- FK: pokedex id
	- evolution mechanism (int? bitmask?)
	- relevant level/criterion
	- think more about required fields

- table 3: Moves
	- PK: move id
	- move type
	- move base PP
	- move max PP?
	- move name

*/

// header found in all Nitro Format files
type NitroFileHeader struct {
	magic [4]byte
	bom uint16
	unknown uint16
	filesize uint32
	unknown2 uint16
	numFrames uint16
}

type NarcFrame interface { FrameFATB | FrameFNTB | FrameFIMG }

type NitroFrame[TFrame NarcFrame] struct {
	magic []byte // will always be 4 bytes, but golang static array inits are inconvenient and annoying
	frameSize uint32
	data TFrame
}

type NarcFile struct {
	NitroHeader NitroFileHeader
	FrameFATB NitroFrame[FrameFATB]
	FrameFNTB NitroFrame[FrameFNTB]
	FrameFIMG NitroFrame[FrameFIMG]
}

func (narc NarcFile) String() string {
	ret := "--- NARC FILE ---\n"
	ret += fmt.Sprint(narc.NitroHeader)
	ret += fmt.Sprint(narc.FrameFATB)
	ret += fmt.Sprint(narc.FrameFNTB)
	ret += fmt.Sprint(narc.FrameFIMG)

	return ret
}

func (nfh NitroFileHeader) String() string {
	ret :=             "--- NFF header ---\n"
	ret += fmt.Sprintf("magic: '%s'\n", string(nfh.magic[:]))
	ret += fmt.Sprintf("filesize: 0x%08X bytes\n", nfh.filesize)
	ret += fmt.Sprintf("# frames: %d\n", nfh.numFrames)
	return ret
}

func (nf NitroFrame[TFrame]) String() string {
	ret := "--- NF Frame ---\n"
	ret += fmt.Sprintf("magic: '%s'\n", string(nf.magic[:]))
	ret += fmt.Sprintf("frame size: 0x%08X bytes\n", nf.frameSize)
	ret += fmt.Sprint(nf.data)

	return ret
}

func newNitroHeader(rom []byte, start uint32) NitroFileHeader {
	file := rom[start:]
	return NitroFileHeader{
		magic: [4]byte{ file[0], file[1], file[2], file[3] },
		bom: utils.U16(file, 4),
		unknown: utils.U16(file ,6),
		filesize: utils.U32(file, 8),
		unknown2: utils.U16(file, 12),
		numFrames: utils.U16(file, 14),
	}
}

func NewNarcFile(rom []byte, start uint32) NarcFile {
	file := rom[start:]
	header := newNitroHeader(file, 0)
	fatb := newFrameFATB(file, 0x10)
	fntb := newFrameFNTB(file, 0x10+fatb.frameSize, fatb.data.numEntries)
	fimg := newFrameFIMG(file, 0x10+fatb.frameSize+fntb.frameSize)
	return NarcFile{
		NitroHeader: header,
		FrameFATB: fatb,
		FrameFNTB: fntb,
		FrameFIMG: fimg,
	}
}
