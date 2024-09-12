package main

import "antispam/actions"
import "antispam/console"
import "os"

func isFile(file string) bool {

	var result bool = false

	stat, err := os.Stat(file)

	if err == nil && !stat.IsDir() {
		result = true
	}

	return result

}

func showHelp() {

	console.Info("Antispam")
	console.Info("E-Mail Spam Classification Tool")

	console.Group("Usage:")
	console.Log("    antispam view path/to/email.eml")
	console.Log("    antispam classify path/to/email.eml")
	console.Log("    antispam mark-spam path/to/email.eml")
	console.Log("")
	console.GroupEnd("------")

	console.Group("Examples:")
	console.Log("    # Views email in Terminal")
	console.Log("    antispam view path/to/email.eml;")
	console.Log("")
	console.Log("    # Exit code 1 if email is spam");
	console.Log("    antispam classify path/to/email.eml;")
	console.Log("    if [[ $? == 1 ]]; then echo \"is spam!\"; fi;")
	console.Log("")
	console.Log("    # Mark email as spam, output new json file for pull-request")
	console.Log("    antispam mark-spam path/to/email.eml;")
	console.GroupEnd("---------")

}

func main() {

	action := ""
	file := ""

	if len(os.Args) == 3 {

		if os.Args[1] == "classify" {

			// antispam classify path/to/email.eml
			if isFile(os.Args[2]) {
				action = "classify"
				file = os.Args[2]
			}

		} else if os.Args[1] == "mark" || os.Args[1] == "mark-spam" {

			// antispam mark path/to/email.eml
			if isFile(os.Args[2]) {
				action = "mark"
				file = os.Args[2]
			}

		} else if os.Args[1] == "view" {

			// antispam view path/to/email.eml
			if isFile(os.Args[2]) {
				action = "view"
				file = os.Args[2]
			}

		}

	}

	if action != "" && file != "" {

		if action == "classify" {

			spammer := actions.Classify(file)

			if spammer != nil {

				console.Warn("Email is classified as spam!")

				console.Warn("Reason: \"" + spammer.Domain + "\"")
				console.Inspect(spammer)

				os.Exit(1)

			} else {

				console.Info("Email is not classified as spam.")
				os.Exit(0)

			}

		} else if action == "mark" {

			spammer := actions.Mark(file)

			if spammer != nil {

				console.Log("Please add a Pull-Request to https://github.com/cookiengineer/antispam")
				console.Group(spammer.Domain + ".json")
				console.Inspect(spammer)
				console.GroupEnd("")

				os.Exit(0)

			} else {

				console.Error("Some error happened. Maybe this was not a known email format?")

				os.Exit(1)

			}

		} else if action == "view" {

			result := actions.View(file)

			if result == true {
				os.Exit(0)
			} else {
				os.Exit(1)
			}

		}

	}

}
