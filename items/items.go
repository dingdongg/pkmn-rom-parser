package items

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dingdongg/pkmn-rom-parser/v3/path_resolver"
)

// fetch item names and cache in RAM to prevent multiple file IO operations
var cache []string

func GetItemName(index uint16) (string, error) {
	if len(cache) == 0 {
		fmt.Println("populating cache")
		path := filepath.Join(path_resolver.GetRoot(), "data", "items-gen4.txt")

		b, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("Unexpected error while reading items file: ", err)
		}

		scanner := bufio.NewScanner(bytes.NewReader(b))

		for scanner.Scan() {
			res := scanner.Text()
			tokens := strings.Split(res, "|")
			id, err := strconv.ParseUint(tokens[0], 0, 16)
			if err != nil {
				log.Fatalf("Unexpected error while parsing line %d: %s\n", id, err)
			}

			cache = append(cache, tokens[1])
		}
	}

	if index >= uint16(len(cache)) {
		return "", errors.New("invalid item ID")
	}

	return cache[index], nil
}
