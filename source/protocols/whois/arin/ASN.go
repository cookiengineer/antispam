package arin

type ASN struct {
	Identifier string `whois:"ASNumber"`
	Name       string `whois:"ASName"`
	Handle     string `whois:"ASHandle"`
	Created    string `whois:"RegDate"`
	Updated    string `whois:"Updated"`
	Reference  string `whois:"Ref"`
}
