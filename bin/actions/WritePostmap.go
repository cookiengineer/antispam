package actions

import "postfix-spamdb/structs"
import "bytes"
import "os"
import "sort"
import "strings"

func WritePostmap(entries []structs.Entry, filepath string) bool {

	var result bool = false
	var blocked_domains map[string]bool = make(map[string]bool)
	var blocked_networks map[string]bool = make(map[string]bool)

	for e := 0; e < len(entries); e++ {

		var entry = entries[e]

		if entry.Domain != "" {
			blocked_domains[entry.Domain] = true
		}

		if len(entry.Aliases) > 0 {

			for a := 0; a < len(entry.Aliases); a++ {
				blocked_domains[entry.Aliases[a]] = true
			}

		}

		if len(entry.Networks) > 0 {

			for n := 0; n < len(entry.Networks); n++ {
				blocked_networks[entry.Networks[n]] = true
			}

		}

	}

	var domains []string

	for domain, _ := range blocked_domains {
		domains = append(domains, domain)
	}

	sort.Strings(domains)

	var networks []string

	for network, _ := range blocked_networks {
		networks = append(networks, network)
	}

	sort.Strings(networks)


	if strings.Contains(filepath, "clients") {

		var buffer bytes.Buffer

		for n := 0; n < len(networks); n++ {
			buffer.WriteString(networks[n]+" REJECT Your network is spam\n")
		}

		buffer.WriteString("\n")

		err := os.WriteFile(filepath, buffer.Bytes(), 0666)

		if err == nil {
			result = true
		}

	} else if strings.Contains(filepath, "senders") {

		var buffer bytes.Buffer

		for d := 0; d < len(domains); d++ {
			buffer.WriteString(domains[d]+" REJECT Your domain is spam\n")
		}

		buffer.WriteString("\n")

		err := os.WriteFile(filepath, buffer.Bytes(), 0666)

		if err == nil {
			result = true
		}

	}

	return result

}
