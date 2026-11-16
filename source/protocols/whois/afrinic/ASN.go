package afrinic

type ASN struct {
	Identifier  string   `whois:"aut-num"`
	Name        string   `whois:"as-name"`
	Description string   `whois:"descr"`
	Status      string   `whois:"status"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Source      string   `whois:"source"`
}
