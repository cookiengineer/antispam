package strings

import "strings"

func Cut(value string, prefix string, suffix string) string {

	var result string

	if strings.Contains(value, prefix) && strings.Contains(value, suffix) {

		index1 := strings.Index(value, prefix)
		index2 := strings.Index(value, suffix)

		if index2 > index1 {
			result = value[index1+len(prefix):index2]
		}

	}

	return result

}
