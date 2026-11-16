package apnic

type ASN struct {
	Identifier  string   `whois:"aut-num"`
	Name        string   `whois:"as-name"`
	Description string   `whois:"descr"`
	Country     string   `whois:"country"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Updated     string   `whois:"last-modified"`
	Source      string   `whois:"source"`
}
