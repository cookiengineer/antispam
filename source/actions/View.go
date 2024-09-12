package actions

import "antispam/console"
import "antispam/insights"
import "antispam/structs"
import "antispam/types"
import "os"
import "strings"

func View(file string) bool {

	var result bool = false
	var spammer *structs.Spammer = nil

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
							spammer = check
							break
						}

					}

				}

				if spammer == nil {

					for i := 0; i < len(email.IPv4s); i++ {

						check := insights.Spammers.SearchIPv4(email.IPv4s[i])

						if check != nil {
							spammer = check
							break
						}

					}

				}

				if spammer == nil {

					for i := 0; i < len(email.IPv6s); i++ {

						check := insights.Spammers.SearchIPv6(email.IPv6s[i])

						if check != nil {
							spammer = check
							break
						}

					}

				}

				console.Group("Headers")

				console.Log("ID:       " + email.MessageID)
				console.Log("Boundary: " + email.Boundary)

				if email.From != "" {

					if spammer != nil {
						console.Error("From:     " + email.From)
					} else {
						console.Log("From:     " + email.From)
					}

				}

				if email.To != "" {
					console.Log("To:       " + email.To)
				}

				if email.Subject != "" {
					console.Log("Subject:  " + email.Subject)
				} else {
					console.Log("Subject:  (no subject)")
				}

				if len(email.Domains) > 0 {
					console.Log("Domains:  " + strings.Join(email.Domains, ", "))
				}

				if len(email.IPv4s) > 0 {
					console.Log("IPv4s:    " + strings.Join(email.IPv4s, ", "))
				}

				if len(email.IPv6s) > 0 {
					console.Log("IPv6s:    " + strings.Join(email.IPv6s, ", "))
				}

				console.GroupEnd("-------")

				// TODO: email.Date

				if email.Message != "" {

					console.Group("Message")

					lines := strings.Split(email.Message, "\n")

					for l := 0; l < len(lines); l++ {
						console.Log(lines[l])
					}

					console.GroupEnd("--------")

				}

				result = true

			} else {

				console.Error("Cannot parse file \"" + file + "\"")

			}

		} else {

			console.Error("Cannot parse file \"" + file + "\"")

		}

	} else {

		console.Error("Cannot read file \"" + file + "\"")

	}

	return result

	// TODO: Parse E-Mail
	// TODO: Lookup Domains in HostMap
	// TODO: Lookup IPs in SpammerMap

	// TODO: Render E-Mail as text
	// TODO: If spammer found, then render email headers via console.Warn()
	// TODO: else render them via console.Log()

}
