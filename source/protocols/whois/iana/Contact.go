package iana

type Contact struct {
	Type         string   `whois:"contact"`
	Name         string   `whois:"name"`
	Organization string   `whois:"organisation"`
	Address      []string `whois:"address"`
	Phone        string   `whois:"phone"`
	Fax          string   `whois:"fax-no"`
	EMail        string   `whois:"e-mail"`
}
