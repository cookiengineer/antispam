package strings

import "encoding/json"
import _ "embed"
import "strings"

//go:embed ToASCII.UTF8.json
var embedded_utf8_table []byte

var utf8_table map[string][]uint16

func init() {
	json.Unmarshal(embedded_utf8_table, &utf8_table)
}

func lookup(character uint16) string {

	var found string

	for str, codes := range utf8_table {

		for c := 0; c < len(codes); c++ {

			if codes[c] == character {
				found = str
				break
			}

		}

	}

	return found

}

func ToASCII(value string) string {

	tmp := strings.TrimSpace(value)

	var filtered string

	for _, character := range tmp {

		var mapped = lookup(uint16(character))

		if mapped != "" {
			filtered += mapped
		}

	}

	return filtered

}
