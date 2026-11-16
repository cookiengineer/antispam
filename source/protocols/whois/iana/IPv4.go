package iana

type IPv4 struct {
	Refer        string `whois:"refer"`
	AddressBlock string `whois:"inetnum"`
	Status       string `whois:"status"`
	Whois        string `whois:"whois"`
	Changed      string `whois:"changed"`
	Source       string `whois:"source"`
}
