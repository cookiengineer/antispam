package main

import "antispam/console"
import "antispam/insights"
import "antispam/utils/postfix"
import "os"
import "path/filepath"

func showPostfixUsage() {

	console.Info("")
	console.Info("toolchain/postfix")
	console.Info("")

	console.Group("Usage: postfix [Action]")
	console.GroupEnd("------")

	console.Group("Examples")
	console.Log("# Generate postmap files to /build folder")
	console.Log("go run postfix.go generate;")
	console.GroupEnd("--------")

}

func main() {

	action := ""

	if len(os.Args) == 2 {

		if os.Args[1] == "generate" {
			action = "generate"
		}

	}

	folder, err := os.Getwd()

	if err == nil {

		if action == "generate" {

			build_folder := filepath.Dir(folder) + "/build"

			blocked_clients := postfix.GenerateBlockedClients(&insights.Hosts, &insights.Spammers)
			blocked_senders := postfix.GenerateBlockedSenders(&insights.Hosts, &insights.Spammers)


			err1 := os.WriteFile(build_folder + "/blocked_clients", blocked_clients, 0666)
			err2 := os.WriteFile(build_folder + "/blocked_senders", blocked_senders, 0666)

			if err1 == nil {
				console.Log("Generated \"" + build_folder + "/blocked_clients\"")
			}

			if err2 == nil {
				console.Log("Generated \"" + build_folder + "/blocked_senders\"")
			}

			if err1 == nil && err2 == nil {
				os.Exit(0)
			} else {
				os.Exit(1)
			}

		} else {
			showPostfixUsage()
			os.Exit(1)
		}

	}

}
