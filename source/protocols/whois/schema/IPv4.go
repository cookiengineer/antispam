package schema

import "strings"

type IPv4 struct {
	AddressBlock           string   `json:"inetnum" whois:"inetnum"`
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

func NewIPv4() IPv4 {

	var ipv4 IPv4

	ipv4.Contacts = make([]string, 0)
	ipv4.Maintainers = make([]string, 0)

	return ipv4

}

func (schema *IPv4) AddContact(value string) {

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

func (schema *IPv4) AddContacts(values []string) {

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

func (schema *IPv4) AddMaintainer(value string) {

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

func (schema *IPv4) AddMaintainers(values []string) {

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

func (schema *IPv4) SetAddressBlock(value string) {

	if strings.Contains(value, " - ") {
		schema.AddressBlock = strings.TrimSpace(value)
	} else if strings.Contains(value, "/") {
		schema.AddressBlock = strings.TrimSpace(value)
	}

}

func (schema *IPv4) SetCountry(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.Country = value
	}

}

func (schema *IPv4) SetNetworkName(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.NetworkName = value
	}

}

func (schema *IPv4) SetOrganizationIdentifier(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.OrganizationIdentifier = value
	}

}

func (schema *IPv4) SetOrganizationName(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.OrganizationName = value
	}

}

func (schema *IPv4) SetOrigin(value string) {

	if strings.HasPrefix(value, "AS") {
		schema.Origin = value
	}

}

func (schema *IPv4) SetSource(value string) {

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
	}

}

func (schema *IPv4) SetStatus(value string) {

	if value == "LEGACY" {
		schema.Status = "LEGACY"
	} else if value == "ALLOCATED" {
		schema.Status = "ALLOCATED"
	} else if value == "ASSIGNED" || value == "ASSIGNED PA" {
		schema.Status = "ASSIGNED"
	}

}
