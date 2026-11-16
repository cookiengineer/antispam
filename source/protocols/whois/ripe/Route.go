package ripe

type Route struct {
	Route        string   `whois:"route,route6"`
	AddressBlock string   `whois:"inetnum,inet6num"`
	Origin       string   `whois:"origin"`
	Description  string   `whois:"descr"`
	Maintainers  []string `whois:"mnt-by"`
	Created      string   `whois:"created"`
	Updated      string   `whois:"last-modified"`
	Source       string   `whois:"source"`
}
