package arin

type Route struct {
	Route       string   `whois:"route,route6"`
	Origin      string   `whois:"origin"`
	Description string   `whois:"descr"`
	Contacts    []string `whois:"abuse-c,admin-c,tech-c"`
	Maintainers []string `whois:"mnt-by"`
	Created     string   `whois:"created"`
	Updated     string   `whois:"last-modified"`
	Source      string   `whois:"source"`
}
