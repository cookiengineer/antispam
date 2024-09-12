package main

import "antispam/actions"
import "antispam/console"
import "encoding/json"
import "fmt"
import "os"
import "strings"

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

	console.Group("Usage: antispam [Action] [Flag] [File]")

	console.Log("")
	console.Group("Exit Code | Description |")
	console.Log("----------|-------------|")
	console.Log("    0     | Not Spam    |")
	console.Log("    1     | Spam        |")
	console.Log("    2     | Parse Error |")
	console.GroupEnd("----------|-------------|")

	console.Log("")
	console.Group("Action   | Description                            | Output          |")
	console.Log("---------|----------------------------------------|-----------------|")
	console.Log("classify | Classify an E-Mail as Spam or Not Spam | (Known) Spammer |")
	console.Log("mark     | Mark an E-Mail as Spam                 | (New)   Spammer |")
	console.Log("view     | View an E-Mail for manual inspection   |         Email   |")
	console.GroupEnd("---------|----------------------------------------|-----------------|")

	console.Log("")
	console.Group("Flag    | Description                          |")
	console.Log("--------|--------------------------------------|")
	console.Log(" --json | Output parsed data structure as JSON |")
	console.GroupEnd("--------|--------------------------------------|")

	console.GroupEnd("------")

	console.Group("Examples:")
	console.Log("    # Views Email in Terminal")
	console.Log("    antispam view path/to/email.eml;")
	console.Log("")
	console.Log("    # Exit code 1 if Email is spam");
	console.Log("    antispam classify path/to/email.eml;")
	console.Log("    if [[ $? == 1 ]]; then echo \"is spam!\"; fi;")
	console.Log("")
	console.Log("    # Mark email as spam, output Spammer JSON for pull-request")
	console.Log("    antispam mark --json path/to/email.eml;")
	console.GroupEnd("---------")

}

func main() {

	action := ""
	file := ""
	json_output := false

	if len(os.Args) == 4 {

		if os.Args[1] == "classify" {

			if os.Args[2] == "--json" {

				// antispam classify --json path/to/email.eml
				if isFile(os.Args[3]) {
					action = "classify"
					file = os.Args[3]
					json_output = true
				}

			}

		} else if os.Args[1] == "view" {

			if os.Args[2] == "--json" {

				// antispam view --json path/to/email.eml
				if isFile(os.Args[3]) {
					action = "view"
					file = os.Args[3]
					json_output = true
				}

			}

		}

	} else if len(os.Args) == 3 {

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

			if json_output == true {

				console.Disable(console.FeatureAll)

				spammer, _ := actions.Classify(file)

				if spammer != nil {

					json, err := json.MarshalIndent(spammer, "", "\t")

					if err == nil {
						fmt.Println(strings.TrimSpace(string(json)))
						os.Exit(1)
					} else {
						os.Exit(1)
					}

				} else {
					os.Exit(0)
				}

			} else {

				spammer, reasons := actions.Classify(file)

				if spammer != nil {

					console.Error("Rating:   E-Mail is classified as spam!")
					console.Error("Spammer:  \"" + spammer.Domain + "\"")
					console.Error("Reasons:  " + strings.Join(reasons, ", "))
					os.Exit(1)

				} else {

					console.Info("Rating:    E-Mail is not classified as spam.")
					os.Exit(0)

				}

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

			if json_output == true {

				console.Disable(console.FeatureAll)

				email, spammer := actions.View(file)

				if email != nil {

					json, err := json.MarshalIndent(email, "", "\t")

					if err == nil {
						fmt.Println(strings.TrimSpace(string(json)))
					}

					if spammer != nil {
						os.Exit(1)
					} else {
						os.Exit(0)
					}

				} else {
					os.Exit(2)
				}

			} else {

				email, spammer := actions.View(file)

				if email != nil {

					if spammer != nil {
						os.Exit(1)
					} else {
						os.Exit(0)
					}

				} else {
					console.Error("Cannot parse file \"" + file + "\"")
					os.Exit(2)
				}

			}

		} else {

			showHelp()
			os.Exit(2)

		}

	} else {

		showHelp()
		os.Exit(2)

	}

}
