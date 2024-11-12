package insights

import "antispam/types"
import "strings"

func IsUnblockable(reason string) bool {

	var result bool = false

	if types.IsIPv4(reason) {

		check := UnblockableSpammers.SearchIPv4(reason)

		if check != nil {
			result = true
		}

	} else if types.IsIPv6(reason) {

		check := UnblockableSpammers.SearchIPv6(reason)

		if check != nil {
			result = true
		}

	} else if types.IsDomain(reason) {

		// Strip out subdomain obfuscations
		tmp := strings.Split(reason, ".")

		if len(tmp) > 2 {
			reason = strings.Join(tmp[len(tmp)-2:], ".")
		}

		check := UnblockableSpammers.SearchDomain(reason)

		if check != nil {
			result = true
		}

	}

	return result

}
