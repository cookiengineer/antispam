package iana

type IPv6 struct {
	Refer        string `whois:"refer"`
	AddressBlock string `whois:"inet6num"`
	Status       string `whois:"status"`
	Whois        string `whois:"whois"`
	Changed      string `whois:"changed"`
	Source       string `whois:"source"`
}
