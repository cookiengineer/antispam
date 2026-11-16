package whois

import "antispam/types"
import "antispam/protocols/whois/afrinic"
import "antispam/protocols/whois/apnic"
import "antispam/protocols/whois/arin"
import "antispam/protocols/whois/lacnic"
import "antispam/protocols/whois/ripe"
import "antispam/protocols/whois/schema"
import recon_types "antispam/types"

func ResolveIPv6(subject string, server *recon_types.Server) schema.IPv6 {

	result := schema.NewIPv6()

	if types.IsIPv6(subject) {

		ipv6 := types.ParseIPv6(subject)

		if ipv6 != nil {

			if server == nil {
				server = NewRegistryServer(ipv6.String())
			}

			if server != nil && server.Domain != "" {

				if server.Schema == "AFRINIC" {

					var tmp_ipv6 afrinic.IPv6
					var tmp_organization afrinic.Organization
					var tmp_routes []afrinic.Route

					str := ipv6.String()
					buffer := Resolve(str[1:len(str)-1], server)

					err1 := Unmarshal(buffer, &tmp_ipv6)
					err2 := Unmarshal(buffer, &tmp_organization)
					err3 := Unmarshal(buffer, &tmp_routes)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv6.AddressBlock)
						result.SetNetworkName(tmp_ipv6.NetName)
						result.SetOrganizationName(tmp_ipv6.Organization)
						result.SetCountry(tmp_ipv6.Country)
						result.AddContacts(tmp_ipv6.Contacts)
						result.AddMaintainers(tmp_ipv6.Maintainers)
						result.SetStatus(tmp_ipv6.Status)
						result.SetSource("AFRINIC")
					}

					if err2 == nil {
						result.SetOrganizationIdentifier(tmp_organization.Identifier)
						result.SetOrganizationName(tmp_organization.Name)
						result.SetCountry(tmp_organization.Country)
						result.AddContacts(tmp_organization.Contacts)
						result.AddMaintainers(tmp_organization.Maintainers)
					}

					if err3 == nil && len(tmp_routes) > 0 {

						for r := 0; r < len(tmp_routes); r++ {

							route := tmp_routes[r]

							if route.Origin != "" {
								// Use first route's ASN
								result.SetOrigin(route.Origin)
								break
							}

						}

					} else {

						result.SetStatus("ALLOCATED")

						origin := ResolveASN(ipv6.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "APNIC" {

					var tmp_ipv6 apnic.IPv6
					var tmp_organization apnic.Organization
					var tmp_routes []apnic.Route

					str := ipv6.String()
					buffer := Resolve(str[1:len(str)-1], server)

					err1 := Unmarshal(buffer, &tmp_ipv6)
					err2 := Unmarshal(buffer, &tmp_organization)
					err3 := Unmarshal(buffer, &tmp_routes)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv6.AddressBlock)
						result.SetNetworkName(tmp_ipv6.NetName)
						result.SetCountry(tmp_ipv6.Country)
						result.AddContacts(tmp_ipv6.Contacts)
						result.AddMaintainers(tmp_ipv6.Maintainers)
						result.SetStatus(tmp_ipv6.Status)
						result.SetSource("APNIC")
					}

					if err2 == nil {
						result.SetOrganizationIdentifier(tmp_organization.Identifier)
						result.SetOrganizationName(tmp_organization.Name)
						result.SetCountry(tmp_organization.Country)
						result.AddContacts(tmp_organization.Contacts)
						result.AddMaintainers(tmp_organization.Maintainers)
					}

					if err3 == nil && len(tmp_routes) > 0 {

						result.SetStatus("ASSIGNED")

						for r := 0; r < len(tmp_routes); r++ {

							route := tmp_routes[r]

							if route.Origin != "" {
								// Use first route's ASN
								result.SetOrigin(route.Origin)
								break
							}

						}

					} else {

						result.SetStatus("ALLOCATED")

						origin := ResolveASN(ipv6.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "ARIN" {

					var tmp_ipv6 arin.IPv6
					var tmp_organization arin.Organization

					str := ipv6.String()
					buffer := Resolve(str[1:len(str)-1], server)

					err1 := Unmarshal(buffer, &tmp_ipv6)
					err2 := Unmarshal(buffer, &tmp_organization)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv6.AddressBlock)
						result.SetNetworkName(tmp_ipv6.NetName)
						result.SetOrganizationName(tmp_ipv6.Organization)
						result.SetOrigin(tmp_ipv6.Origin)
						result.SetCountry("US")
						result.SetSource("ARIN")
					}

					if err2 == nil {
						result.SetOrganizationIdentifier(tmp_organization.Identifier)
						result.SetOrganizationName(tmp_organization.Name)
						result.SetCountry(tmp_organization.Country)
					}

					// ARIN doesn't respond with routes
					if result.Origin != "" {

						result.SetStatus("ASSIGNED")

					} else {

						result.SetStatus("ALLOCATED")

						origin := ResolveASN(ipv6.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "LACNIC" {

					var tmp_ipv6 lacnic.IPv6
					var tmp_contact lacnic.Contact

					str := ipv6.String()
					buffer := Resolve(str[1:len(str)-1], server)

					err1 := Unmarshal(buffer, &tmp_ipv6)
					err2 := Unmarshal(buffer, &tmp_contact)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv6.AddressBlock)
						result.SetOrigin(tmp_ipv6.Origin)
						result.SetOrganizationIdentifier(tmp_ipv6.OwnerIdentifier)
						result.SetOrganizationName(tmp_ipv6.Owner)
						result.SetCountry(tmp_ipv6.Country)
						result.AddContacts(tmp_ipv6.Contacts)
						result.AddMaintainers(tmp_ipv6.Maintainers)
						result.SetStatus(tmp_ipv6.Status)
						result.SetSource("LACNIC")
					}

					if err2 == nil {
						result.SetCountry(tmp_contact.Country)
						result.AddMaintainer(tmp_contact.Identifier)
					}

					if result.Origin != "" {

						result.SetStatus("ASSIGNED")

					} else {

						result.SetStatus("ALLOCATED")

						origin := ResolveASN(ipv6.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "RIPE" {

					var tmp_ipv6 ripe.IPv6
					var tmp_organization ripe.Organization
					var tmp_routes []ripe.Route

					str := ipv6.String()
					buffer := Resolve(str[1:len(str)-1], server)

					err1 := Unmarshal(buffer, &tmp_ipv6)
					err2 := Unmarshal(buffer, &tmp_organization)
					err3 := Unmarshal(buffer, &tmp_routes)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv6.AddressBlock)
						result.SetNetworkName(tmp_ipv6.NetName)
						result.SetOrganizationName(tmp_ipv6.Organization)
						result.SetCountry(tmp_ipv6.Country)
						result.AddContacts(tmp_ipv6.Contacts)
						result.AddMaintainers(tmp_ipv6.Maintainers)
						result.SetStatus(tmp_ipv6.Status)
						result.SetSource("RIPE")
					}

					if err2 == nil {
						result.SetOrganizationIdentifier(tmp_organization.Identifier)
						result.SetOrganizationName(tmp_organization.Name)
						result.SetCountry(tmp_organization.Country)
						result.AddContacts(tmp_organization.Contacts)
						result.AddMaintainers(tmp_organization.Maintainers)
					}

					if err3 == nil && len(tmp_routes) > 0 {

						result.SetStatus("ASSIGNED")

						for r := 0; r < len(tmp_routes); r++ {

							route := tmp_routes[r]

							if route.Origin != "" {
								// Use first route's ASN
								result.SetOrigin(route.Origin)
								break
							}

						}

					} else {

						result.SetStatus("ALLOCATED")

						origin := ResolveASN(ipv6.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				}

			}

		}

	}

	return result

}
