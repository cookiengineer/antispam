package schema

import "strings"

type Route struct {
	AddressBlock string   `json:"route" whois:"route,route6"`
	Origin       string   `json:"origin" whois:"origin"`
	Country      string   `json:"country" whois:"country"`
	Contacts     []string `json:"contacts" whois:"abuse-c,admin-c,tech-c"`
	Maintainers  []string `json:"maintainers" whois:"mnt-by,mnt-irt"`
	Source       string   `json:"source" whois:"source"`
}

func NewRoute() Route {

	var route Route

	route.Contacts = make([]string, 0)
	route.Maintainers = make([]string, 0)
	route.Country = "XX" // ISO Code for unknown country

	return route

}

func (schema *Route) AddContact(value string) {

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

func (schema *Route) AddContacts(values []string) {

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

func (schema *Route) AddMaintainer(value string) {

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

func (schema *Route) AddMaintainers(values []string) {

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

func (schema *Route) SetAddressBlock(value string) {

	if strings.Contains(value, " - ") {
		schema.AddressBlock = strings.TrimSpace(value)
	} else if strings.Contains(value, "/") {
		schema.AddressBlock = strings.TrimSpace(value)
	}

}

func (schema *Route) SetCountry(value string) {

	value = strings.TrimSpace(value)

	if value != "" {
		schema.Country = value
	}

}

func (schema *Route) SetOrigin(value string) {

	if strings.HasPrefix(value, "AS") {
		schema.Origin = value
	}

}

func (schema *Route) SetSource(value string) {

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
	} else if value == "LEVEL3" {
		schema.Source = "LEVEL3"
	} else if value == "RIPE" {
		schema.Source = "RIPE"
	}

}
