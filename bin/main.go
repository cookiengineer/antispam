package main

import "postfix-spamdb/actions"
import "postfix-spamdb/structs"
import "path/filepath"
import "fmt"
import "os"
import "strings"

func main() {

	var target_folder string
	var root string

	cwd, err1 := os.Getwd()

	if err1 == nil {
		root = cwd
	}

	if len(os.Args) == 2 {

		tmp, err1 := filepath.Abs(os.Args[1])

		if err1 == nil {
			target_folder = tmp
		}

	} else {

		if root != "" {

			tmp, err2 := filepath.Abs(root + "/../build")

			if err2 == nil {
				target_folder = tmp
			}

		}

	}

	if target_folder != "" {

		var all_entries []structs.Entry

		stat, err1 := os.Stat(target_folder)

		if err1 == nil && stat.IsDir() {

			source_folder, err2 := filepath.Abs(root+"/../source")

			if err2 == nil {

				files, err3 := os.ReadDir(source_folder)

				if err3 == nil && len(files) > 0 {

					for f := 0; f < len(files); f++ {

						var file = files[f].Name()

						if strings.HasSuffix(file, ".json") {

							entries := actions.ReadJSON(source_folder+"/"+file)

							for e := 0; e < len(entries); e++ {
								all_entries = append(all_entries, entries[e])
							}

						} else if strings.HasSuffix(file, ".hosts") {

							entries := actions.ReadHosts(source_folder+"/"+file)

							for e := 0; e < len(entries); e++ {
								all_entries = append(all_entries, entries[e])
							}

						}

					}

				}

			}

		}

		if len(all_entries) > 0 {

			result1 := actions.WritePostmap(all_entries, target_folder+"/blocked_clients")

			if result1 == true {
				fmt.Println("- Successfully generated \""+target_folder+"/blocked_clients\"")
			} else {
				fmt.Println("! Could not generate \""+target_folder+"/blocked_clients\"!")
			}

			result2 := actions.WritePostmap(all_entries, target_folder+"/blocked_senders")

			if result2 == true {
				fmt.Println("- Successfully generated \""+target_folder+"/blocked_senders\"")
			} else {
				fmt.Println("! Could not generate \""+target_folder+"/blocked_senders\"!")
			}

		}

	}

}
