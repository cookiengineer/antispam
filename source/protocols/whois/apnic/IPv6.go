package apnic

type IPv6 struct {
	AddressBlock string   `whois:"inet6num"`
	NetName      string   `whois:"netname"`
	Organization string   `whois:"org"`
	Description  string   `whois:"descr"`
	Country      string   `whois:"country"`
	Status       string   `whois:"status"`
	Contacts     []string `whois:"abuse-c,admin-c,tech-c"`
	Maintainers  []string `whois:"mnt-by,mnt-irt"`
	Updated      string   `whois:"last-modified"`
	Source       string   `whois:"source"`
}
