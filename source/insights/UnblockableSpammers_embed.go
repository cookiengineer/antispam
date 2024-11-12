package insights

import "antispam/structs"
import "encoding/json"
import "embed"
import "strings"

//go:embed spammers/unblockable/*.json
var embedded_UnblockableSpammers embed.FS

func init() {

	UnblockableSpammers = structs.NewSpammerMap()
	files, err1 := embedded_UnblockableSpammers.ReadDir("spammers/unblockable")

	if err1 == nil {

		for _, file := range files {

			name := file.Name()
			buffer, err2 := embedded_UnblockableSpammers.ReadFile("spammers/unblockable/"+name)

			if err2 == nil {

				var instances []structs.Spammer

				if strings.HasSuffix(name, ".json") {

					err3 := json.Unmarshal(buffer, &instances)

					if err3 == nil {

						for i := 0; i < len(instances); i++ {
							UnblockableSpammers.AddSpammer(instances[i])
						}

					}

				}

				UnblockableSpammers.FlushAliases()

			}

		}

	}

}
