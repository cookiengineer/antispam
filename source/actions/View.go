package actions

import "antispam/console"
import "antispam/insights"
import "antispam/structs"
import "antispam/types"
import "os"
import "strings"
import "time"

func View(file string) (*types.Email, *structs.Spammer) {

	var email *types.Email = nil
	var spammer *structs.Spammer = nil

	stat, err0 := os.Stat(file)

	if err0 == nil && !stat.IsDir() {

		buffer, err1 := os.ReadFile(file)

		if err1 == nil && types.IsEmail(buffer) {

			email = types.ParseEmail(buffer)

			if email != nil {

				reasons := make([]string, 0)

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

				console.Group("Headers")

				console.Log("ID:       " + email.MessageID)
				console.Log("Boundary: " + email.Boundary)
				console.Log("From:     " + email.From)
				console.Log("To:       " + email.To)
				console.Log("Date:     " + email.Date.Format(time.RFC3339))

				if email.Subject != "" {
					console.Log("Subject:  " + email.Subject)
				} else {
					console.Log("Subject:  (no subject)")
				}

				console.Log("Domains:  " + strings.Join(email.Domains, ", "))
				console.Log("IPv4s:    " + strings.Join(email.IPv4s, ", "))
				console.Log("IPv6s:    " + strings.Join(email.IPv6s, ", "))

				if spammer != nil {
					console.Error("Rating:   E-Mail is classified as spam!")
					console.Error("Spammer:  \"" + spammer.Domain + "\"")
					console.Error("Reasons:  " + strings.Join(reasons, ", "))
				} else {
					console.Info("Rating:    E-Mail is not classified as spam.")
				}

				console.GroupEnd("-------")

				if email.Message != "" {

					console.Group("Message")

					lines := strings.Split(email.Message, "\n")

					for l := 0; l < len(lines); l++ {
						console.Log(lines[l])
					}

					console.GroupEnd("--------")

				}

			}

		}

	}

	return email, spammer

}
