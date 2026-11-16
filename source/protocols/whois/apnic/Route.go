package apnic

type Route struct {
	Route       string   `whois:"route,route6"`
	Origin      string   `whois:"origin"`
	Description string   `whois:"descr"`
	Maintainers []string `whois:"mnt-by"`
	Updated     string   `whois:"last-modified"`
	Source      string   `whois:"source"`
}
