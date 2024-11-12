package main

import "antispam/console"
import "antispam/insights"
import "antispam/types"
import "os"
import "path/filepath"
import "strings"

func main() {

	console.Group("Cleanup")

	cleanup_domain := ""
	cleanup_from := ""
	cleanup_spam := false

	if len(os.Args) == 2 {

		if strings.HasPrefix(os.Args[1], "--domain=") {

			tmp := os.Args[1][9:]

			if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
				tmp = tmp[1:len(tmp)-1]
			}

			cleanup_domain = tmp

		} else if strings.HasPrefix(os.Args[1], "--from=") {

			tmp := os.Args[1][7:]

			if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
				tmp = tmp[1:len(tmp)-1]
			}

			cleanup_from = tmp

		} else if os.Args[1] == "--spam" {

			cleanup_spam = true

		}

	}

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

							if cleanup_from != "" {

								if strings.HasSuffix(email.From, cleanup_from) {

									err := os.Remove(mails_folder + "/" + name)

									if err == nil {
										console.Info("Removed \"" + name + "\"")
									}

								}

							} else if cleanup_domain != "" {

								matches_domain := false

								for d := 0; d < len(email.Domains); d++ {

									domain := email.Domains[d]

									if domain == cleanup_domain || strings.HasSuffix(domain, "." + cleanup_domain) {
										matches_domain = true
										break
									}

								}

								if matches_domain == true {

									err := os.Remove(mails_folder + "/" + name)

									if err == nil {
										console.Info("Removed \"" + name + "\"")
									}

								}

							} else if cleanup_spam == true {

								spammer, _ := insights.Classify(email)

								if spammer != nil {

									err := os.Remove(mails_folder + "/" + name)

									if err == nil {
										console.Info("Removed \"" + name + "\"")
									}

								}

							}

						}

					}

				}

			}

		}

	}

	console.GroupEnd("Cleanup")

}
