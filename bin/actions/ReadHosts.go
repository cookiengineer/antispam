package actions

import "postfix-spamdb/structs"
import "os"
import "strings"

func ReadHosts(filepath string) []structs.Entry {

	var entries []structs.Entry

	buffer, err1 := os.ReadFile(filepath)

	if err1 == nil && len(buffer) > 0 {

		lines := strings.Split(strings.TrimSpace(string(buffer)), "\n")

		for l := 0; l < len(lines); l++ {

			var domain = strings.TrimSpace(lines[l])

			if strings.Contains(domain, ".") {

				entry := structs.Entry{
					Domain: domain,
				}

				entries = append(entries, entry)

			}

		}

	}

	return entries

}
