package main

import "antispam/console"
import "antispam/insights"
import "antispam/types"
import "os"
import "path/filepath"
import "slices"
import "sort"
import "strings"
import "time"

var NotSpam map[string]*types.Email
var Spam map[string]*types.Email

func init() {
	NotSpam = make(map[string]*types.Email)
	Spam = make(map[string]*types.Email)
}

func main() {

	console.Group("Learn")

	toolchain_folder, err0 := os.Getwd()

	if err0 == nil && strings.HasSuffix(toolchain_folder, "/toolchain") {

		mails_folder := filepath.Dir(toolchain_folder) + "/mails"

		entries, err1 := os.ReadDir(mails_folder)

		if err1 == nil {

			for _, file := range entries {

				name := file.Name()

				// Ignore .gitkeep and other hidden files
				if strings.HasPrefix(name, ".") != true {

					buffer, err2 := os.ReadFile(mails_folder + "/" + name)

					if err2 == nil && types.IsEmail(buffer) {

						email := types.ParseEmail(buffer)

						if email != nil {

							spammer, _ := insights.Classify(email)

							if spammer != nil {
								Spam[name] = email
							} else {
								NotSpam[name] = email
							}

						} else {
							console.Error("Cannot parse \"" + name + "\"")
						}

					} else {
						console.Error("Cannot parse \"" + name + "\"")
					}

				}

			}

			names_spam := make([]string, 0)

			for name, _ := range Spam {
				names_spam = append(names_spam, name)
			}

			sort.Strings(names_spam)

			for n := 0; n < len(names_spam); n++ {

				name := names_spam[n]
				email := Spam[name]

				spammer, reasons := insights.Classify(email)

				if len(reasons) == 1 && spammer.Domain == reasons[0] {

					console.Group(name)
					console.Error("Blocked via \"" + spammer.Domain + "\" because of " + strings.Join(reasons, ", "))
					console.GroupEnd(name)

				} else {

					console.Group(name)
					console.Error("Blocked via \"" + spammer.Domain + "\" because of " + strings.Join(reasons, ", "))

					for d := 0; d < len(email.Domains); d++ {

						domain := email.Domains[d]

						if strings.Contains(domain, ".") && !slices.Contains(reasons, domain) && !types.IsLocalDomain(domain) {

							if !insights.IsBlocked(domain) && !insights.IsUnblockable(domain) {
								console.Warn("> " + domain)
							}

						}

					}

					console.GroupEnd(name)

				}

			}

			names_notspam := make([]string, 0)

			for name, _ := range NotSpam {
				names_notspam = append(names_notspam, name)
			}

			sort.Strings(names_notspam)

			for n := 0; n < len(names_notspam); n++ {

				name := names_notspam[n]
				email := NotSpam[name]

				console.Log("")
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
				console.Info("Rating:   E-Mail is not classified as spam.")

				console.GroupEnd("-------")

				console.GroupEnd(name)

			}

		}

	}

	console.GroupEnd("Learn")

}
