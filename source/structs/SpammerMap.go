package structs

import "antispam/types"
import "encoding/hex"
import "sort"
import "strconv"
import "strings"

type SpammerMap struct {
	Domains map[string]*Spammer           `json:"domains"`
	IPv4s   map[uint8]map[string]*Spammer `json:"ipv4s"`
	IPv6s   map[uint8]map[string]*Spammer `json:"ipv6s"`
}

func NewSpammerMap() SpammerMap {

	var spammermap SpammerMap

	spammermap.Domains = make(map[string]*Spammer)
	spammermap.IPv4s = make(map[uint8]map[string]*Spammer)
	spammermap.IPv6s = make(map[uint8]map[string]*Spammer)

	return spammermap

}

func (spammermap *SpammerMap) AddSpammer(value Spammer) bool {

	var result bool = false

	if value.Domain != "" {

		spammermap.Domains[value.Domain] = &value

		for n := 0; n < len(value.Networks); n++ {

			network := value.Networks[n]

			if strings.Contains(network, "/") {

				raw_ip := network[0:strings.Index(network, "/")]
				raw_prefix := network[strings.Index(network, "/")+1:]

				if types.IsIPv4(raw_ip) {

					ipv4 := types.ParseIPv4(raw_ip)
					num, err := strconv.ParseUint(raw_prefix, 10, 8)

					if ipv4 != nil && err == nil && num >= 8 && num <= 32 {

						prefix := uint8(num)
						bytes := ipv4.Bytes(prefix)
						hash := hex.EncodeToString(bytes)

						_, ok1 := spammermap.IPv4s[prefix]

						if ok1 == false {
							spammermap.IPv4s[prefix] = make(map[string]*Spammer)
						}

						_, ok2 := spammermap.IPv4s[prefix][hash]

						if ok2 == false {
							spammermap.IPv4s[prefix][hash] = &value
						}

					}

				} else if types.IsIPv6(raw_ip) {

					ipv6 := types.ParseIPv6(raw_ip)
					num, err := strconv.ParseUint(raw_prefix, 10, 8)

					if ipv6 != nil && err == nil && num >= 16 && num <= 128 {

						prefix := uint8(num)
						bytes := ipv6.Bytes(prefix)
						hash := hex.EncodeToString(bytes)

						_, ok1 := spammermap.IPv6s[prefix]

						if ok1 == false {
							spammermap.IPv6s[prefix] = make(map[string]*Spammer)
						}

						_, ok2 := spammermap.IPv6s[prefix][hash]

						if ok2 == false {
							spammermap.IPv6s[prefix][hash] = &value
						}

					}

				}

			}

		}

	}

	return result

}

func (spammermap *SpammerMap) FlushAliases() {

	for domain, spammer := range spammermap.Domains {

		for a := 0; a < len(spammer.Aliases); a++ {

			alias := spammer.Aliases[a]
			_, ok := spammermap.Domains[alias]

			if ok == false {
				spammermap.Domains[alias] = spammermap.Domains[domain]
			}

		}

	}

}

func (spammermap *SpammerMap) SearchDomain(value string) *Spammer {

	var result *Spammer = nil

	tmp, ok := spammermap.Domains[value]

	if ok == true {
		result = tmp
	}

	return result

}

func (spammermap *SpammerMap) SearchIPv4(value string) *Spammer {

	var result *Spammer = nil

	if types.IsIPv4(value) {

		ipv4 := types.ParseIPv4(value)

		if ipv4 != nil {

			prefixes := make([]uint8, 0)

			for prefix := range spammermap.IPv4s {
				prefixes = append(prefixes, prefix)
			}

			sort.Slice(prefixes, func(a int, b int) bool {
				return prefixes[a] > prefixes[b]
			})

			for p := 0; p < len(prefixes); p++ {

				prefix := prefixes[p]
				hash := hex.EncodeToString(ipv4.Bytes(prefix))

				tmp, ok := spammermap.IPv4s[prefix][hash]

				if ok == true {
					result = tmp
					break
				}

			}

		}

	}

	return result

}

func (spammermap *SpammerMap) SearchIPv6(value string) *Spammer {

	var result *Spammer = nil

	if types.IsIPv6(value) {

		ipv6 := types.ParseIPv6(value)

		if ipv6 != nil {

			prefixes := make([]uint8, 0)

			for prefix := range spammermap.IPv6s {
				prefixes = append(prefixes, prefix)
			}

			sort.Slice(prefixes, func(a int, b int) bool {
				return prefixes[a] > prefixes[b]
			})

			for p := 0; p < len(prefixes); p++ {

				prefix := prefixes[p]
				hash := hex.EncodeToString(ipv6.Bytes(prefix))

				tmp, ok := spammermap.IPv6s[prefix][hash]

				if ok == true {
					result = tmp
					break
				}

			}

		}

	}

	return result

}

