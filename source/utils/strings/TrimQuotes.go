package strings

import "strings"

func TrimQuotes(value string) string {

	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1:len(value)-1]
	}

	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		value = value[1:len(value)-1]
	}

	return value

}
