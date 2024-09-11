package structs

type Spammer struct {
	Domain   string   `json:"domain"`
	Aliases  []string `json:"aliases"`
	Networks []string `json:"networks"`
}
