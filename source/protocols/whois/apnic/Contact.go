package apnic

type Contact struct {
	Type        string   `whois:"role,person,irt"`
	Address     []string `whois:"address"`
	Identifier  string   `whois:"nic-hdl"`
	Phone       string   `whois:"phone"`
	Email       string   `whois:"e-mail,abuse-mailbox"`
	Contacts    []string `whois:"abuse-c,admin-c,tech-c"`
	Maintainers []string `whois:"mnt-by,mnt-irt"`
	Updated     string   `whois:"last-modified"`
	Source      string   `whois:"source"`
}
