package iana

type Organization struct {
	Name    string   `whois:"organisation"`
	Address []string `whois:"address"`
}
