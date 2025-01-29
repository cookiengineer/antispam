package main

import "antispam/console"
import "antispam/insights"
import "antispam/types"
import "os"
import "slices"
import "strconv"
import "strings"

func main() {

	discover_domain := ""

	console.Group("Discover")

	if len(os.Args) == 2 {

		if strings.HasPrefix(os.Args[1], "--domain=") {

			tmp := os.Args[1][9:]

			if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
				tmp = tmp[1:len(tmp)-1]
			}

			discover_domain = strings.TrimSpace(tmp)

		}

	}

	if discover_domain != "" {

		// IPv4 map from "A.B" to []C subnets
		ipv4_candidates := make(map[string][]uint8, 0)

		for domain, spammer := range insights.Spammers.Domains {

			if domain == discover_domain {

				for n := 0; n < len(spammer.Networks); n++ {

					network := spammer.Networks[n]

					// Well, gotta change this once Russia's SVR changes this strategy
					if strings.HasSuffix(network, "/24") {

						ipv4 := types.ParseIPv4(network[0:strings.Index(network, "/")])

						if ipv4 != nil {

							identifier := strconv.FormatUint(uint64(uint8(ipv4[0])), 10) + "." + strconv.FormatUint(uint64(uint8(ipv4[1])), 10)

							_, ok := ipv4_candidates[identifier]

							if ok == false {
								ipv4_candidates[identifier] = []uint8{uint8(ipv4[2])}
							} else {
								ipv4_candidates[identifier] = append(ipv4_candidates[identifier], uint8(ipv4[2]))
							}

						}

					}

				}

			}

		}

		for a_b, subnets := range ipv4_candidates {

			slices.Sort(subnets)

			if len(subnets) >= 2 {

				min_subnet := subnets[0]
				max_subnet := subnets[1]

				if min_subnet < max_subnet && min_subnet + 1 < max_subnet {

					console.Group("Potential Candidates: " + a_b + "." + strconv.FormatUint(uint64(min_subnet), 10) + ".0/24 to " + a_b + "." + strconv.FormatUint(uint64(max_subnet), 10) + ".0/24")

					for s := min_subnet + 1; s < max_subnet; s++ {

						if !slices.Contains(subnets, s) {

							ipv4 := a_b + "." + strconv.FormatUint(uint64(s), 10) + ".0"
							check := insights.Spammers.SearchIPv4(ipv4)

							if check == nil {
								console.Log("> " + ipv4)
							}

						}

					}

					console.GroupEnd("")

				}

			}

		}

	}

	console.GroupEnd("Discover")

}
