package iana

type Domain struct {
	Refer       string   `whois:"refer"`
	Name        string   `whois:"domain"`
	Nameservers []string `whois:"nserver"`
	Whois       string   `whois:"whois"`
	Status      string   `whois:"status"`
	Created     string   `whois:"created"`
	Changed     string   `whois:"changed"`
	Source      string   `whois:"source"`
}
