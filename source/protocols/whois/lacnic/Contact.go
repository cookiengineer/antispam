package lacnic

type Contact struct {
	Identifier string   `whois:"nic-hdl,nic-hdl-br"`
	Type       string   `whois:"role,person,irt"`
	Email      string   `whois:"e-mail,abuse-mailbox"`
	Address    []string `whois:"address"`
	Country    string   `whois:"country"`
	Phone      string   `whois:"phone"`
	Created    string   `whois:"created"`
	Updated    string   `whois:"changed"`
}
