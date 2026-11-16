package whois

import "antispam/types"
import recon_types "antispam/types"

func ResolveRegistry(subject string, server *recon_types.Server) string {

	result := ""

	if types.IsIPv4(subject) {

		ipv4 := types.ParseIPv4(subject)

		if ipv4 != nil {

			origin := ResolveASN(ipv4.String())

			if origin != "" {
				return ResolveRegistry(origin, server)
			}

		}

	} else if types.IsIPv6(subject) {

		ipv6 := types.ParseIPv6(subject)

		if ipv6 != nil {

			origin := ResolveASN(ipv6.String())

			if origin != "" {
				return ResolveRegistry(origin, server)
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
					result = server.Schema
				} else if server.Schema == "APNIC" {
					result = server.Schema
				} else if server.Schema == "ARIN" {
					result = server.Schema
				} else if server.Schema == "LACNIC" {
					result = server.Schema
				} else if server.Schema == "RIPE" {
					result = server.Schema
				}

			}

		}

	}

	return result

}
