package whois

import "antispam/types"
import "antispam/protocols/whois/afrinic"
import "antispam/protocols/whois/apnic"
import "antispam/protocols/whois/arin"
import "antispam/protocols/whois/lacnic"
import "antispam/protocols/whois/ripe"
import "antispam/protocols/whois/schema"
import recon_types "antispam/types"

func ResolveRoutes(subject string, server *recon_types.Server) []schema.Route {

	result := make([]schema.Route, 0)

	if types.IsIPv4(subject) {

		ipv4 := types.ParseIPv4(subject)

		if ipv4 != nil {

			origin := ResolveASN(ipv4.String())

			if origin != "" {
				return ResolveRoutes(origin, server)
			}

		}

	} else if types.IsIPv6(subject) {

		ipv6 := types.ParseIPv6(subject)

		if ipv6 != nil {

			origin := ResolveASN(ipv6.String())

			if origin != "" {
				return ResolveRoutes(origin, server)
			}

		}

	} else if types.IsASN(subject) {

		asn := types.ParseASN(subject)

		if asn != nil {

			if server == nil {
				server = NewRegistryServer(asn.String())
			}

			if server != nil && server.Domain != "" {

				if server.Schema == "AFRINIC" {

					var routes []afrinic.Route

					buffer := Resolve("-i origin "+asn.String(), server)
					err := Unmarshal(buffer, &routes)

					if err == nil && len(routes) > 0 {

						for r := 0; r < len(routes); r++ {

							route := routes[r]

							if route.Origin != "" {

								tmp := schema.NewRoute()
								tmp.SetAddressBlock(route.Route)
								tmp.SetOrigin(route.Origin)
								tmp.AddMaintainers(route.Maintainers)
								tmp.SetSource("AFRINIC")

								result = append(result, tmp)

							}

						}

					} else {

						result = resolveRoutesFallback(asn.String())

					}

				} else if server.Schema == "APNIC" {

					var routes []apnic.Route

					buffer := Resolve("-i origin "+asn.String(), server)
					err := Unmarshal(buffer, &routes)

					if err == nil && len(routes) > 0 {

						for r := 0; r < len(routes); r++ {

							route := routes[r]

							if route.Origin != "" {

								tmp := schema.NewRoute()
								tmp.SetAddressBlock(route.Route)
								tmp.SetOrigin(route.Origin)
								tmp.AddMaintainers(route.Maintainers)
								tmp.SetSource("APNIC")

								result = append(result, tmp)

							}

						}

					} else {

						result = resolveRoutesFallback(asn.String())

					}

				} else if server.Schema == "ARIN" {

					var routes []arin.Route
					var server = NewServer("rr.arin.net")

					buffer := Resolve("-i origin "+asn.String(), &server)
					err := Unmarshal(buffer, &routes)

					if err == nil && len(routes) > 0 {

						for r := 0; r < len(routes); r++ {

							route := routes[r]

							if route.Origin != "" {

								tmp := schema.NewRoute()
								tmp.SetAddressBlock(route.Route)
								tmp.SetOrigin(route.Origin)
								tmp.AddContacts(route.Maintainers)
								tmp.AddMaintainers(route.Maintainers)
								tmp.SetSource("ARIN")

								result = append(result, tmp)

							}

						}

					} else {

						result = resolveRoutesFallback(asn.String())

					}

				} else if server.Schema == "LACNIC" {

					var routes []lacnic.Route
					var server = NewServer("irr.lacnic.net")

					buffer := Resolve("-i origin "+asn.String(), &server)
					err := Unmarshal(buffer, &routes)

					if err == nil && len(routes) > 0 {

						for r := 0; r < len(routes); r++ {

							route := routes[r]

							if route.Origin != "" {

								tmp := schema.NewRoute()
								tmp.SetAddressBlock(route.Route)
								tmp.SetOrigin(route.Origin)
								tmp.AddMaintainers(route.Maintainers)
								tmp.SetSource("LACNIC")

								result = append(result, tmp)

							}

						}

					} else {

						result = resolveRoutesFallback(asn.String())

					}

				} else if server.Schema == "RIPE" {

					var routes []ripe.Route

					buffer := Resolve("-i origin "+asn.String(), server)
					err := Unmarshal(buffer, &routes)

					if err == nil && len(routes) > 0 {

						for r := 0; r < len(routes); r++ {

							route := routes[r]

							if route.Origin != "" {

								tmp := schema.NewRoute()
								tmp.SetAddressBlock(route.Route)
								tmp.SetOrigin(route.Origin)
								tmp.AddMaintainers(route.Maintainers)
								tmp.SetSource("RIPE")

								result = append(result, tmp)

							}

						}

					} else {

						result = resolveRoutesFallback(asn.String())

					}

				} else {

					result = resolveRoutesFallback(asn.String())

				}

			}

		}

	}

	return result

}
