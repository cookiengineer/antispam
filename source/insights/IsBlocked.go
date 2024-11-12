package insights

import "antispam/types"
import "strings"

func IsBlocked(reason string) bool {

	var result bool = false

	if types.IsIPv4(reason) {

		check := Spammers.SearchIPv4(reason)

		if check != nil {
			result = true
		}

	} else if types.IsIPv6(reason) {

		check := Spammers.SearchIPv6(reason)

		if check != nil {
			result = true
		}

	} else if types.IsDomain(reason) {

		// Strip out subdomain obfuscations
		tmp := strings.Split(reason, ".")

		if len(tmp) > 2 {
			reason = strings.Join(tmp[len(tmp)-2:], ".")
		}

		check1 := Hosts.SearchDomain(reason)

		if check1 == true {
			result = true
		}

		if result == false {

			check2 := Spammers.SearchDomain(reason)

			if check2 != nil {
				result = true
			}

		}

	}

	return result

}
