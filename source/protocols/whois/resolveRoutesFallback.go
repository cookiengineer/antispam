package whois

import "antispam/types"
import "antispam/protocols/whois/level3"
import "antispam/protocols/whois/schema"

func resolveRoutesFallback(subject string) []schema.Route {

	result := make([]schema.Route, 0)

	if types.IsASN(subject) {

		asn := types.ParseASN(subject)

		if asn != nil {

			level3_routes := make([]level3.Route, 0)
			level3_server := NewServer("rr.level3.net")

			buffer_routes := Resolve("-i origin "+asn.String(), &level3_server)

			err3 := Unmarshal(buffer_routes, &level3_routes)

			if err3 == nil && len(level3_routes) > 0 {

				for r := 0; r < len(level3_routes); r++ {

					route := level3_routes[r]

					if route.Origin != "" {

						new_route := schema.NewRoute()
						new_route.SetAddressBlock(route.Route)
						new_route.SetOrigin(route.Origin)
						new_route.AddContacts(route.Contacts)
						new_route.AddMaintainers(route.Maintainers)
						new_route.SetSource("LEVEL3")

						result = append(result, new_route)

					}

				}

			}

		}

	}

	return result

}
