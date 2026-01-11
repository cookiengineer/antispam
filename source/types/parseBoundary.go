package types

import utils_strings "antispam/utils/strings"
import "strings"

func parseBoundary(content string) string {

	result := ""
	lines := strings.Split(strings.TrimSpace(content), "\n")

	for l := 0; l < len(lines); l++ {

		line := strings.TrimSpace(lines[l])

		if strings.HasPrefix(line, "boundary=") && strings.HasSuffix(line, ";") {

			result = utils_strings.TrimQuotes(strings.TrimSpace(line[9:len(line)-1]))
			break

		} else if strings.HasPrefix(line, "boundary=") {

			result = utils_strings.TrimQuotes(strings.TrimSpace(line[9:]))
			break

		} else if strings.HasSuffix(line, ";") {
			// Do Nothing
		} else if strings.Contains(line, ";") {

			chunks := strings.Split(line, ";")

			for c := 0; c < len(chunks); c++ {

				chunk := strings.TrimSpace(chunks[c])

				if strings.HasPrefix(chunk, "boundary=") {
					result = utils_strings.TrimQuotes(strings.TrimSpace(chunk[9:]))
					break
				}

			}

		}

		if strings.HasPrefix(line, "--") {
			break
		}

	}

	return result

}

