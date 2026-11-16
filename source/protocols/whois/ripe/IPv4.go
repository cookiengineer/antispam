package ripe

type IPv4 struct {
	AddressBlock string   `whois:"inetnum"`
	NetName      string   `whois:"netname"`
	Organization string   `whois:"org"`
	Description  string   `whois:"descr"`
	Country      string   `whois:"country"`
	Status       string   `whois:"status"`
	Contacts     []string `whois:"abuse-c,admin-c,tech-c"`
	Maintainers  []string `whois:"mnt-by,mnt-irt"`
	Created      string   `whois:"created"`
	Updated      string   `whois:"last-modified"`
	Source       string   `whois:"source"`
}
