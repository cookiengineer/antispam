package actions

import "antispam/console"
import "antispam/insights"
import "antispam/structs"
import "antispam/types"
import "os"

func Classify(file string) (*structs.Spammer, []string) {

	var spammer *structs.Spammer = nil
	var reasons []string = make([]string, 0)

	stat, err0 := os.Stat(file)

	if err0 == nil && !stat.IsDir() {

		buffer, err1 := os.ReadFile(file)

		if err1 == nil && types.IsEmail(buffer) {

			email := types.ParseEmail(buffer)

			if email != nil {

				if spammer == nil {

					for d := 0; d < len(email.Domains); d++ {

						check := insights.Spammers.SearchDomain(email.Domains[d])

						if check != nil {
							reasons = append(reasons, email.Domains[d])
							spammer = check
							break
						}

					}

				}

				if spammer == nil {

					for i := 0; i < len(email.IPv4s); i++ {

						check := insights.Spammers.SearchIPv4(email.IPv4s[i])

						if check != nil {
							reasons = append(reasons, email.IPv4s[i])
							spammer = check
							break
						}

					}

				}

				if spammer == nil {

					for i := 0; i < len(email.IPv6s); i++ {

						check := insights.Spammers.SearchIPv6(email.IPv6s[i])

						if check != nil {
							reasons = append(reasons, email.IPv6s[i])
							spammer = check
							break
						}

					}

				}

			} else {

				console.Error("Cannot parse file \"" + file + "\"")

			}

		} else {

			console.Error("Cannot parse file \"" + file + "\"")

		}

	} else {

		console.Error("Cannot read file \"" + file + "\"")

	}

	return spammer, reasons

}
