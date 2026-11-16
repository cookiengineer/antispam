package afrinic

type Organization struct {
	Identifier  string   `whois:"organisation"`
	Name        string   `whois:"org-name"`
	Type        string   `whois:"org-type"`
	Country     string   `whois:"country"`
	Address     []string `whois:"address"`
	Phone       string   `whois:"phone"`
	Email       string   `whois:"e-mail"`
	Contacts    []string `whois:"abuse-c,admin-c,tech-c"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Source      string   `whois:"source"`
}
