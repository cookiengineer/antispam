package ripe

type ASN struct {
	Identifier  string   `whois:"aut-num"`
	Name        string   `whois:"as-name"`
	Status      string   `whois:"status"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Created     string   `whois:"created"`
	Updated     string   `whois:"last-modified"`
	Source      string   `whois:"source"`
}
