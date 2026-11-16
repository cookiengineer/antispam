package iana

type ASN struct {
	Refer        string `whois:"refer"`
	AddressBlock string `whois:"as-block"`
	Organization string `whois:"organisation"`
	Whois        string `whois:"whois"`
	Description  string `whois:"descr"`
	Source       string `whois:"source"`
}
