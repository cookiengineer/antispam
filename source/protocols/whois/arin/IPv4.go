package arin

type IPv4 struct {
	AddressBlock string `whois:"NetRange"`
	CIDR         string `whois:"CIDR"`
	NetName      string `whois:"NetName"`
	Handle       string `whois:"NetHandle"`
	Parent       string `whois:"Parent"`
	Type         string `whois:"NetType"`
	Origin       string `whois:"OriginAS"`
	Organization string `whois:"Organization"`
	Created      string `whois:"RegDate"`
	Updated      string `whois:"Updated"`
	Reference    string `whois:"Ref"`
}
