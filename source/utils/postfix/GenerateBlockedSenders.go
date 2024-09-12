package postfix

import "antispam/structs"
import "bytes"
import "sort"

func GenerateBlockedSenders(hostmap *structs.HostMap, spammermap *structs.SpammerMap) []byte {

	domains := make(map[string]bool, 0)

	for host, _ := range hostmap.Hosts {

		_, ok := domains[host]

		if ok == false {
			domains[host] = true
		}

	}

	for domain, _ := range spammermap.Domains {

		_, ok := domains[domain]

		if ok == false {
			domains[domain] = true
		}

	}

	collected := make([]string, 0)

	for domain, _ := range domains {
		collected = append(collected, domain)
	}

	sort.Strings(collected)

	var buffer bytes.Buffer

	for c := 0; c < len(collected); c++ {
		buffer.WriteString(collected[c] + " REJECT Your domain is spam\n")
	}

	buffer.WriteString("\n")

	return buffer.Bytes()

}
