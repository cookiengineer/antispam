package whois

import "antispam/types"
import "antispam/protocols/whois/level3"
import "strings"

func ResolveASN(subject string) string {

	result := ""

	if strings.Contains(subject, "/") {

		ip := subject[0:strings.Index(subject, "/")]

		if types.IsIPv4(ip) {
			subject = ip
		} else if types.IsIPv6(ip) {
			subject = ip
		}

	}

	if types.IsIPv4(subject) {

		ipv4 := types.ParseIPv4(subject)

		if ipv4 != nil {

			level3_routes := make([]level3.Route, 0)
			level3_server := NewServer("rr.level3.net")
			buffer_routes := Resolve(ipv4.String(), &level3_server)

			err3 := Unmarshal(buffer_routes, &level3_routes)

			if err3 == nil && len(level3_routes) > 0 {

				for r := 0; r < len(level3_routes); r++ {

					route := level3_routes[r]

					if types.IsASN(route.Origin) {
						// Use first route's ASN
						result = types.ParseASN(route.Origin).String()
						break
					}

				}

			}

		}

	} else if types.IsIPv6(subject) {

		ipv6 := types.ParseIPv6(subject)

		if ipv6 != nil {

			level3_routes := make([]level3.Route, 0)
			level3_server := NewServer("rr.level3.net")

			str := ipv6.String()
			buffer_routes := Resolve(str[1:len(str)-1], &level3_server)

			err3 := Unmarshal(buffer_routes, &level3_routes)

			if err3 == nil && len(level3_routes) > 0 {

				for r := 0; r < len(level3_routes); r++ {

					route := level3_routes[r]

					if types.IsASN(route.Origin) {
						// Use first route's ASN
						result = types.ParseASN(route.Origin).String()
						break
					}

				}

			}

		}

	} else if types.IsASN(subject) {

		asn := types.ParseASN(subject)

		if asn != nil {
			result = asn.String()
		}

	}

	return result

}
