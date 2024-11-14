package main

import "antispam/insights"
import "antispam/types"
import "io"
import "os"

func main() {

	// Exit Code 0: Not Spam, keep Email
	// Exit Code 1: Spam, discard Email

	buffer, err0 := io.ReadAll(os.Stdin)

	if err0 == nil {

		if types.IsEmail(buffer) {

			email := types.ParseEmail(buffer)

			if email != nil {

				spammer, _ := insights.Classify(email)

				if spammer != nil {
					os.Exit(1)
				} else {
					os.Exit(0)
				}

			} else {
				os.Exit(0)
			}

		} else {
			os.Exit(0)
		}

	} else {
		os.Exit(0)
	}

}
