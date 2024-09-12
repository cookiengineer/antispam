package postfix

import "antispam/structs"
import "antispam/types"
import "bytes"
import "encoding/hex"
import "sort"
import "strconv"

func GenerateBlockedClients(hostmap *structs.HostMap, spammermap *structs.SpammerMap) []byte {

	var buffer bytes.Buffer

	hashes_ipv4 := make([]string, 0)
	hashes_ipv6 := make([]string, 0)
	networks_ipv4 := make(map[string]uint8, 0)
	networks_ipv6 := make(map[string]uint8, 0)

	for prefix, hashes := range spammermap.IPv4s {

		for hash, _ := range hashes {

			other_prefix, ok := networks_ipv4[hash]

			if ok == true {

				if prefix < other_prefix {
					networks_ipv4[hash] = prefix
				}

			} else {
				networks_ipv4[hash] = prefix
			}

		}

	}

	for prefix, hashes := range spammermap.IPv6s {

		for hash, _ := range hashes {

			other_prefix, ok := networks_ipv6[hash]

			if ok == true {

				if prefix < other_prefix {
					networks_ipv6[hash] = prefix
				}

			} else {
				networks_ipv6[hash] = prefix
			}

		}

	}

	for hash, _ := range networks_ipv4 {
		hashes_ipv4 = append(hashes_ipv4, hash)
	}

	if len(hashes_ipv4) > 0 {

		sort.Strings(hashes_ipv4)

		for h := 0; h < len(hashes_ipv4); h++ {

			hash := hashes_ipv4[h]
			prefix := networks_ipv4[hash]
			bytes, _ := hex.DecodeString(hash)
			ipv4 := types.IPv4(bytes)

			buffer.WriteString(ipv4.String() + "/" + strconv.FormatUint(uint64(prefix), 10) + " REJECT Your network is spam\n")

		}

		buffer.WriteString("\n")

	}

	for hash, _ := range networks_ipv6 {
		hashes_ipv6 = append(hashes_ipv6, hash)
	}

	if len(hashes_ipv6) > 0 {

		sort.Strings(hashes_ipv6)

		for h := 0; h < len(hashes_ipv6); h++ {

			hash := hashes_ipv6[h]
			prefix := networks_ipv6[hash]
			bytes, _ := hex.DecodeString(hash)
			ipv6 := types.IPv6(bytes)

			buffer.WriteString(ipv6.String() + "/" + strconv.FormatUint(uint64(prefix), 10) + " REJECT Your network is spam\n")

		}

		buffer.WriteString("\n")

	}

	return buffer.Bytes()

}
