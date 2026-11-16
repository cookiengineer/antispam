package whois

import "antispam/types"
import "antispam/protocols/whois/afrinic"
import "antispam/protocols/whois/apnic"
import "antispam/protocols/whois/arin"
import "antispam/protocols/whois/lacnic"
import "antispam/protocols/whois/ripe"
import "antispam/protocols/whois/schema"
import recon_types "antispam/types"

func ResolveIPv4(subject string, server *recon_types.Server) schema.IPv4 {

	result := schema.NewIPv4()

	if types.IsIPv4(subject) {

		ipv4 := types.ParseIPv4(subject)

		if ipv4 != nil {

			if server == nil {
				server = NewRegistryServer(subject)
			}

			if server != nil && server.Domain != "" {

				if server.Schema == "AFRINIC" {

					var tmp_ipv4 afrinic.IPv4
					var tmp_organization afrinic.Organization
					var tmp_routes []afrinic.Route

					buffer := Resolve(ipv4.String(), server)

					err1 := Unmarshal(buffer, &tmp_ipv4)
					err2 := Unmarshal(buffer, &tmp_organization)
					err3 := Unmarshal(buffer, &tmp_routes)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv4.AddressBlock)
						result.SetNetworkName(tmp_ipv4.NetName)
						result.SetOrganizationName(tmp_ipv4.Organization)
						result.SetCountry(tmp_ipv4.Country)
						result.AddContacts(tmp_ipv4.Contacts)
						result.AddMaintainers(tmp_ipv4.Maintainers)
						result.SetStatus(tmp_ipv4.Status)
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

						origin := ResolveASN(ipv4.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "APNIC" {

					var tmp_ipv4 apnic.IPv4
					var tmp_organization apnic.Organization
					var tmp_routes []apnic.Route

					buffer := Resolve(ipv4.String(), server)

					err1 := Unmarshal(buffer, &tmp_ipv4)
					err2 := Unmarshal(buffer, &tmp_organization)
					err3 := Unmarshal(buffer, &tmp_routes)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv4.AddressBlock)
						result.SetNetworkName(tmp_ipv4.NetName)
						result.SetCountry(tmp_ipv4.Country)
						result.AddContacts(tmp_ipv4.Contacts)
						result.AddMaintainers(tmp_ipv4.Maintainers)
						result.SetStatus(tmp_ipv4.Status)
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

						origin := ResolveASN(ipv4.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "ARIN" {

					var tmp_ipv4 arin.IPv4
					var tmp_organization arin.Organization

					buffer := Resolve(ipv4.String(), server)

					err1 := Unmarshal(buffer, &tmp_ipv4)
					err2 := Unmarshal(buffer, &tmp_organization)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv4.AddressBlock)
						result.SetNetworkName(tmp_ipv4.NetName)
						result.SetOrganizationName(tmp_ipv4.Organization)
						result.SetOrigin(tmp_ipv4.Origin)
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

						origin := ResolveASN(ipv4.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "LACNIC" {

					var tmp_ipv4 lacnic.IPv4
					var tmp_contact lacnic.Contact

					buffer := Resolve(ipv4.String(), server)

					err1 := Unmarshal(buffer, &tmp_ipv4)
					err2 := Unmarshal(buffer, &tmp_contact)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv4.AddressBlock)
						result.SetOrigin(tmp_ipv4.Origin)
						result.SetOrganizationIdentifier(tmp_ipv4.OwnerIdentifier)
						result.SetOrganizationName(tmp_ipv4.Owner)
						result.SetCountry(tmp_ipv4.Country)
						result.AddContacts(tmp_ipv4.Contacts)
						result.AddMaintainers(tmp_ipv4.Maintainers)
						result.SetStatus(tmp_ipv4.Status)
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

						origin := ResolveASN(ipv4.String())

						if origin != "" {
							result.SetStatus("ASSIGNED")
							result.SetOrigin(origin)
						}

					}

				} else if server.Schema == "RIPE" {

					var tmp_ipv4 ripe.IPv4
					var tmp_organization ripe.Organization
					var tmp_routes []ripe.Route

					buffer := Resolve(ipv4.String(), server)

					err1 := Unmarshal(buffer, &tmp_ipv4)
					err2 := Unmarshal(buffer, &tmp_organization)
					err3 := Unmarshal(buffer, &tmp_routes)

					if err1 == nil {
						result.SetAddressBlock(tmp_ipv4.AddressBlock)
						result.SetNetworkName(tmp_ipv4.NetName)
						result.SetOrganizationName(tmp_ipv4.Organization)
						result.SetCountry(tmp_ipv4.Country)
						result.AddContacts(tmp_ipv4.Contacts)
						result.AddMaintainers(tmp_ipv4.Maintainers)
						result.SetStatus(tmp_ipv4.Status)
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

						origin := ResolveASN(ipv4.String())

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
