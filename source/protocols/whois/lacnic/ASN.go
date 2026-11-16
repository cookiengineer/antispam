package lacnic

type ASN struct {
	Identifier    string   `whois:"aut-num"`
	Name          string   `whois:"as-name"`
	Status        string   `whois:"status"`
	Maintainers   []string `whois:"mnt-by,mnt-irt"`
	Created       string   `whois:"created"`
	Updated       string   `whois:"changed"`
	AddressBlocks []string `whois:"inetnum"`
}
