package arin

type Organization struct {
	Name       string   `whois:"OrgName"`
	Identifier string   `whois:"OrgId"`
	Address    []string `whois:"Address,City,StateProv,PostalCode"`
	Country    string   `whois:"Country"`
	Created    string   `whois:"RegDate"`
	Updated    string   `whois:"Updated"`
	Reference  string   `whois:"Ref"`
}
