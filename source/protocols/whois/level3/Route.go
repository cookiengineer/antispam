package level3

type Route struct {
	Route       string   `whois:"route,route6"`
	Description string   `whois:"descr"`
	Origin      string   `whois:"origin"`
	Contacts    []string `whois:"abuse-c,admin-c,tech-c"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Created     string   `whois:"created"`
	Changed     string   `whois:"changed,last-modified"`
	Source      string   `whois:"source"`
}
