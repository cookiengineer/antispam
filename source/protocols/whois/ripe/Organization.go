package ripe

type Organization struct {
	Identifier  string   `whois:"organisation"`
	Name        string   `whois:"org-name"`
	Country     string   `whois:"country"`
	Type        string   `whois:"org-type"`
	Address     []string `whois:"address"`
	Phone       string   `whois:"phone"`
	Email       string   `whois:"e-mail"`
	Contacts    []string `whois:"abuse-c,admin-c,tech-c"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Created     string   `whois:"created"`
	Updated     string   `whois:"last-modified"`
	Source      string   `whois:"source"`
}
