package insights

import "antispam/structs"
import "encoding/json"
import "embed"
import "strings"

//go:embed phishers/*.json spammers/*.json
var embedded_Spammers embed.FS

func init() {

	Spammers = structs.NewSpammerMap()
	files, err1 := embedded_Spammers.ReadDir("spammers")

	if err1 == nil {

		for _, file := range files {

			name := file.Name()
			buffer, err2 := embedded_Spammers.ReadFile("spammers/"+name)

			if err2 == nil {

				var instances []structs.Spammer

				if strings.HasSuffix(name, ".json") {

					err3 := json.Unmarshal(buffer, &instances)

					if err3 == nil {

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
