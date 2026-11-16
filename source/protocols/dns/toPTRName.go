package dns

import "antispam/types"
import "strings"

func toPTRName(name string) string {

	result := ""

	if types.IsIPv4(name) == true {

		ipv4 := types.ParseIPv4(name)

		if ipv4 != nil {

			tmp := strings.Split(ipv4.String(), ".")

			for t := len(tmp) - 1; t >= 0; t-- {
				result += tmp[t] + "."
			}

			result += "in-addr.arpa"

		}

	} else if types.IsIPv6(name) == true {

		ipv6 := types.ParseIPv6(name)

		if ipv6 != nil {

			str := ipv6.String()
			tmp := strings.Split(str[1:len(str)-1], ":")

			for t := len(tmp) - 1; t >= 0; t-- {

				result += string(tmp[t][3]) + "."
				result += string(tmp[t][2]) + "."
				result += string(tmp[t][1]) + "."
				result += string(tmp[t][0]) + "."

			}

			result += "ip6.arpa"

		}

	} else {

		result = name

	}

	return result

}
