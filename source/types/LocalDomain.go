package types

import "strings"

var local_domains []string = []string{
	"local",
	"localdomain",
	"fritz.box",
	"home",
	"internal",
	"lan",
}


func IsLocalDomain(value string) bool {

	if IsDomain(value) {

		result := false

		for l := 0; l < len(local_domains); l++ {

			domain := local_domains[l]

			if strings.HasSuffix(value, "." + domain) {
				result = true
				break
			}

		}

		return result

	}

	return false

}
