package utils

import (
	"fmt"
	"strings"
)

func PrintBuffer(buf []byte, w int) {
	output := ""
	buffer := make([]string, 0)

	update := func() {
		buffer = append(buffer, "\n")
		output += strings.Join(buffer, " ")
		buffer = make([]string, 0)
	}

	for i, b := range buf {
		buffer = append(buffer, fmt.Sprintf("%02X", b))
		if i % w == w - 1 {
			update()
		}
	}
	if len(buffer) != 0 {
		update()
	}
	fmt.Println(output)
}