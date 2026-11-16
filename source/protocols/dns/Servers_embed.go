package dns

import "encoding/json"
import _ "embed"
import "log"

//go:embed Servers.json
var embedded_servers []byte

func init() {

	err := json.Unmarshal(embedded_servers, &Servers)

	if err != nil {
		log.Printf("Cannot decompress embedded Servers.json: %s", err.Error())
	}

}
