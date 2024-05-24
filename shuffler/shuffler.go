package shuffler

import (
	"fmt"
	"strings"
)

// from https://projectpokemon.org/home/docs/gen-4/pkm-structure-r65/
const tableValues = `00	ABCD	ABCD
01	ABDC	ABDC
02	ACBD	ACBD
03	ACDB	ADBC
04	ADBC	ACDB
05	ADCB	ADCB
06	BACD	BACD
07	BADC	BADC
08	BCAD	CABD
09	BCDA	DABC
10	BDAC	CADB
11	BDCA	DACB
12	CABD	BCAD
13	CADB	BDAC
14	CBAD	CBAD
15	CBDA	DBAC
16	CDAB	CDAB
17	CDBA	DCAB
18	DABC	BCDA
19	DACB	BDCA
20	DBAC	CBDA
21	DBCA	DBCA
22	DCAB	CDBA
23	DCBA	DCBA`

func Extract() {
	rows := strings.FieldsFunc(tableValues, func(r rune) bool { return r == '\n' })
	var sequences [24]([2]string)

	for i, r := range rows {
		tokens := strings.FieldsFunc(r, func(r rune) bool { return r == '\t' })
		sequences[i] = [2]string{tokens[1], tokens[2]}

		fmt.Println(getShuffleInfo(sequences[i]))
	}
}

func getShuffleInfo(sequence [2]string) string {
	return fmt.Sprintf(
		"{ [4]uint{%c, %c, %c, %c}, [4]uint{%c, %c, %c, %c} }, // %s %s",
		sequence[0][0], sequence[0][1], sequence[0][2], sequence[0][3],
		sequence[1][0], sequence[1][1], sequence[1][2], sequence[1][3],
		sequence[0], sequence[1],
	)
}
