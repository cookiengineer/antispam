package whois

import "strings"

func Parse(buffer []byte) []map[string]string {

	result := make([]map[string]string, 0)
	blocks := strings.Split(strings.TrimSpace(string(buffer)), "\n\n")

	for b := 0; b < len(blocks); b++ {

		block := blocks[b]
		lines := strings.Split(strings.TrimSpace(block), "\n")

		if len(lines) > 0 {

			key := ""
			value := ""
			entry := make(map[string]string, 0)

			for l := 0; l < len(lines); l++ {

				line := lines[l]

				if strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") {

					continue

				} else {

					if key == "descr" && !strings.Contains(line, ": ") && strings.HasPrefix(line, " ") {

						value = value + "\n" + strings.TrimSpace(line)

					} else if key == "descr" && !strings.Contains(line, ": ") && !strings.HasPrefix(line, " ") {

						value = value + strings.TrimSpace(line)

					} else if strings.Contains(line, ": ") && !strings.HasPrefix(line, " ") {

						key = strings.TrimSpace(line[0:strings.Index(line, ":")])
						value = strings.TrimSpace(line[strings.Index(line, ":")+1:])

					} else if strings.Contains(line, ":") && !strings.HasPrefix(line, " ") {

						key = strings.TrimSpace(line[0:strings.Index(line, ":")])
						value = strings.TrimSpace(line[strings.Index(line, ":")+1:])

					} else {

						value = value + "\n" + strings.TrimSpace(line)

					}

					if key != "" && value != "" {
						entry[key] = value
					}

				}

			}

			if len(entry) > 0 {
				result = append(result, entry)
			}

		}

	}

	return result

}
