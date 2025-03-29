package savefile

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/dingdongg/pkmn-rom-parser/v7/char"
	"github.com/dingdongg/pkmn-rom-parser/v7/crypt"
	"github.com/dingdongg/pkmn-rom-parser/v7/revamp/enums"

	"github.com/dingdongg/pkmn-rom-parser/v7/shuffler"
)

type Pokemon struct {
	Name      string
	Level     uint8
	Exp       uint32
	PokedexId uint16
	Nature    enums.Nature
}

func NewPlatSavefile(bytes []byte) *PlatSavefile {
	return &PlatSavefile{
		rawBytes:     bytes,
		partyPokemon: make([]PokemonPtr, 0),
	}
}

func printPokemonInfo(raw []byte) {
	type Test struct {
		field int
	}
	pid := binary.LittleEndian.Uint32(raw[0:4])
	fmt.Printf("PID: 0x%08X\n", pid)
	fmt.Printf("checksum: 0x%04X\n", raw[6:8])

	A, _ := shuffler.GetPokemonBlock(raw, shuffler.A, pid)
	fmt.Printf("Pokedex ID: %d\n", binary.LittleEndian.Uint16(A[0:2]))

	C, _ := shuffler.GetPokemonBlock(raw, shuffler.C, pid)
	name := ""

	for i := 0; i < 22; i += 2 {
		code := binary.LittleEndian.Uint16(C[i : i+2])
		if code == char.END_OF_STRING {
			break
		}
		chr, err := char.Char(code)

		if err != nil {
			log.Fatal("damn", err)
		}
		name += chr
	}

	fmt.Printf("name: '%s'\n", name)
}

func (pt *PlatSavefile) PartyPokemon() []PokemonPtr {
	partySize := binary.LittleEndian.Uint32(pt.rawBytes[0x9C:0xA0])
	for i := range partySize {
		offset := 0xA0 + i*236
		rawPokemon := crypt.DecryptPokemon(pt.rawBytes[offset : offset+236])
		printPokemonInfo(rawPokemon)
		fmt.Println("----------------")
	}

	return pt.partyPokemon
}

func (pt *PlatSavefile) validate() error {
	return nil
}

func (pt *PlatSavefile) Version() enums.GameVersion {
	return enums.PLAT
}

func (pt *PlatSavefile) Flush() error {
	return nil
}
