package main

import "antispam/console"
import "antispam/insights"
import "antispam/utils/postfix"
import "os"
import "path/filepath"

func main() {

	folder, err := os.Getwd()

	if err == nil {

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
		os.Exit(1)
	}

}
