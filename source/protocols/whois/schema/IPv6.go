package schema

import "strings"
import "fmt"

type IPv6 struct {
	AddressBlock           string   `json:"inet6num" whois:"inet6num"`
	NetworkName            string   `json:"netname" whois:"netname"`
	OrganizationIdentifier string   `json:"org-id" whois:"org-id"`
	OrganizationName       string   `json:"org-name" whois:"org-name"`
	Origin                 string   `json:"origin" whois:"origin"`
	Country                string   `json:"country" whois:"country"`
	Status                 string   `json:"status" whois:"status"`
	Contacts               []string `json:"contacts" whois:"abuse-c,admin-c,tech-c"`
	Maintainers            []string `json:"maintainers" whois:"mnt-by,mnt-irt"`
	Source                 string   `json:"source" whois:"source"`
}

func NewIPv6() IPv6 {

	var ipv6 IPv6

	ipv6.Contacts = make([]string, 0)
	ipv6.Maintainers = make([]string, 0)

	return ipv6

}

func (schema *IPv6) AddContact(value string) {

	found := false

	for c := 0; c < len(schema.Contacts); c++ {

		if schema.Contacts[c] == value {
			found = true
			break
		}

	}

	if found == false {
		schema.Contacts = append(schema.Contacts, value)
	}

}

func (schema *IPv6) AddContacts(values []string) {

	for _, value := range values {

		found := false

		for c := 0; c < len(schema.Contacts); c++ {

			if schema.Contacts[c] == value {
				found = true
				break
			}

		}

		if found == false {
			schema.Contacts = append(schema.Contacts, value)
		}

	}

}

func (schema *IPv6) AddMaintainer(value string) {

	found := false

	for m := 0; m < len(schema.Maintainers); m++ {

		if schema.Maintainers[m] == value {
			found = true
			break
		}

	}

	if found == false {
		schema.Maintainers = append(schema.Maintainers, value)
	}

}

func (schema *IPv6) AddMaintainers(values []string) {

	for _, value := range values {

		found := false

		for m := 0; m < len(schema.Maintainers); m++ {

			if schema.Maintainers[m] == value {
				found = true
				break
			}

		}

		if found == false {
			schema.Maintainers = append(schema.Maintainers, value)
		}

	}

}

func (schema *IPv6) SetAddressBlock(value string) {

	if strings.Contains(value, " - ") {
		schema.AddressBlock = strings.TrimSpace(value)
	} else if strings.Contains(value, "/") {
		schema.AddressBlock = strings.TrimSpace(value)
	}

}

func (schema *IPv6) SetCountry(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.Country = value
	}

}

func (schema *IPv6) SetNetworkName(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.NetworkName = value
	}

}

func (schema *IPv6) SetOrganizationIdentifier(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.OrganizationIdentifier = value
	}

}

func (schema *IPv6) SetOrganizationName(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.OrganizationName = value
	}

}

func (schema *IPv6) SetOrigin(value string) {

	if strings.HasPrefix(value, "AS") {
		schema.Origin = value
	}

}

func (schema *IPv6) SetSource(value string) {

	if strings.Contains(value, "#") {
		value = strings.TrimSpace(value[0:strings.Index(value, "#")])
	}

	if value == "IANA" {
		schema.Source = "IANA"
	} else if value == "AFRINIC" {
		schema.Source = "AFRINIC"
	} else if value == "APNIC" {
		schema.Source = "APNIC"
	} else if value == "ARIN" {
		schema.Source = "ARIN"
	} else if value == "LACNIC" {
		schema.Source = "LACNIC"
	} else if value == "RIPE" {
		schema.Source = "RIPE"
	} else {
		fmt.Println(value)
	}

}

func (schema *IPv6) SetStatus(value string) {

	if value == "LEGACY" {
		schema.Status = "LEGACY"
	} else if value == "ALLOCATED" || value == "ALLOCATED-BY-RIR" {
		schema.Status = "ALLOCATED"
	} else if value == "ASSIGNED" || value == "ASSIGNED PA" {
		schema.Status = "ASSIGNED"
	} else {
		fmt.Println(value)
	}

}
