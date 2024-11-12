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

				spammer, reasons = insights.Classify(email)

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
