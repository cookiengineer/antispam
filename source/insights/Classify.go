package insights

import "antispam/structs"
import "antispam/types"

func Classify(email *types.Email) (*structs.Spammer, []string) {

	var spammer *structs.Spammer = nil
	var reasons []string = make([]string, 0)

	if email != nil {

		if spammer == nil {

			for d := 0; d < len(email.Domains); d++ {

				check := Hosts.SearchDomain(email.Domains[d])

				if check == true {

					tmp := structs.Spammer{
						Domain: email.Domains[d],
					}

					reasons = append(reasons, email.Domains[d])
					spammer = &tmp
					break

				}

			}

		}

		if spammer == nil {

			for d := 0; d < len(email.Domains); d++ {

				check := Spammers.SearchDomain(email.Domains[d])

				if check != nil {
					reasons = append(reasons, email.Domains[d])
					spammer = check
					break
				}

			}

		}

		if spammer == nil {

			for i := 0; i < len(email.IPv4s); i++ {

				check := Spammers.SearchIPv4(email.IPv4s[i])

				if check != nil {
					reasons = append(reasons, email.IPv4s[i])
					spammer = check
					break
				}

			}

		}

		if spammer == nil {

			for i := 0; i < len(email.IPv6s); i++ {

				check := Spammers.SearchIPv6(email.IPv6s[i])

				if check != nil {
					reasons = append(reasons, email.IPv6s[i])
					spammer = check
					break
				}

			}

		}

	}

	return spammer, reasons

}
