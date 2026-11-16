package level3

type ASN struct {
	Identifier  string   `whois:"aut-num"`
	Name        string   `whois:"as-name"`
	Description string   `whois:"descr"`
	Imports     []string `whois:"import"`
	Exports     []string `whois:"export"`
	MPImports   []string `whois:"mp-import"`
	MPExports   []string `whois:"mp-export"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Created     string   `whois:"created"`
	Updated     string   `whois:"last-modified"`
	Source      string   `whois:"source"`
}
