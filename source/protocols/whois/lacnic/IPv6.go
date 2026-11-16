package lacnic

type IPv6 struct {
	AddressBlock    string   `whois:"inetnum"`
	Status          string   `whois:"status"`
	Origin          string   `whois:"aut-num"`
	Owner           string   `whois:"owner"`
	OwnerIdentifier string   `whois:"ownerid"`
	Country         string   `whois:"country"`
	Contacts        []string `whois:"abuse-c,admin-c,owner-c,tech-c"`
	Maintainers     []string `whois:"mnt-by,mnt-irt"`
	Created         string   `whois:"created"`
	Updated         string   `whois:"changed"`
}
