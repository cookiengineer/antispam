package insights

import "antispam/structs"
import "encoding/json"
import "embed"
import "strings"

//go:embed phishers/*.json spammers/*.json
var embedded_Spammers embed.FS

func init() {

	Spammers = structs.NewSpammerMap()

	files1, err1 := embedded_Spammers.ReadDir("phishers")

	if err1 == nil {

		for _, file := range files1 {

			name := file.Name()
			buffer, err12 := embedded_Spammers.ReadFile("phishers/"+name)

			if err12 == nil {

				var instances []structs.Spammer

				if strings.HasSuffix(name, ".json") {

					err13 := json.Unmarshal(buffer, &instances)

					if err13 == nil {

						for i := 0; i < len(instances); i++ {
							Spammers.AddSpammer(instances[i])
						}

					}

				}

				Spammers.FlushAliases()

			}

		}

	}

	files2, err2 := embedded_Spammers.ReadDir("spammers")

	if err2 == nil {

		for _, file := range files2 {

			name := file.Name()
			buffer, err22 := embedded_Spammers.ReadFile("spammers/"+name)

			if err22 == nil {

				var instances []structs.Spammer

				if strings.HasSuffix(name, ".json") {

					err23 := json.Unmarshal(buffer, &instances)

					if err23 == nil {

						for i := 0; i < len(instances); i++ {
							Spammers.AddSpammer(instances[i])
						}

					}

				}

				Spammers.FlushAliases()

			}

		}

	}

}
