package insights

import "antispam/structs"
import "embed"
import "strings"

//go:embed hosts/*.hosts
var embedded_Hosts embed.FS

func init() {

	Hosts = structs.NewHostMap()
	files, err1 := embedded_Hosts.ReadDir("hosts")

	if err1 == nil {

		for _, file := range files {

			name := file.Name()
			buffer, err2 := embedded_Hosts.ReadFile("hosts/"+name)

			if err2 == nil {

				lines := strings.Split(strings.TrimSpace(string(buffer)), "\n")

				for l := 0; l < len(lines); l++ {

					line := strings.TrimSpace(lines[l])
					Hosts.AddDomain(line)

				}

			}

		}

	}

}
