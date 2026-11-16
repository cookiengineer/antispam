package structs

import "antispam/types"
import "time"

type Verdict string

const (
	VerdictUnknown     Verdict = "unknown"
	VerdictSpam        Verdict = "spam"
	VerdictUnblockable Verdict = "unblockable"
	VerdictNotSpam     Verdict = "not-spam"
)

type Match struct {
	Domain      string `json:"domain"`
	File        string `json:"file"`
	Kind        string `json:"kind"`
	Unblockable bool   `json:"unblockable"`
}

type FuzzyCandidate struct {
	Domain      string  `json:"domain"`
	File        string  `json:"file"`
	Score       float64 `json:"score"`
	Reason      string  `json:"reason"`
	Unblockable bool    `json:"unblockable"`
}

type RouteCandidate struct {
	Prefix       string `json:"prefix"`
	Source       string `json:"source"`
	ContainsIP   bool   `json:"contains_ip"`
	Missing      bool   `json:"missing"`
	AddressCount uint64 `json:"address_count"`
}

type Indicator struct {
	Kind             string           `json:"kind"`
	Value            string           `json:"value"`
	Private          bool             `json:"private"`
	Authoritative    bool             `json:"authoritative"`
	Known            *Match           `json:"known"`
	Unblockable      *Match           `json:"unblockable"`
	PTR              []string         `json:"ptr"`
	ForwardConfirmed []string         `json:"forward_confirmed"`
	ASN              string           `json:"asn"`
	Organization     string           `json:"organization"`
	NetworkName      string           `json:"network_name"`
	Registry         string           `json:"registry"`
	Country          string           `json:"country"`
	AddressBlock     string           `json:"address_block"`
	Routes           []string         `json:"routes"`
	Candidates       []FuzzyCandidate `json:"candidates"`
	RouteCandidates  []RouteCandidate `json:"route_candidates"`
	LookedUp         bool             `json:"looked_up"`
	Error            string           `json:"error"`
}

type IPLookup struct {
	IP               string    `json:"ip"`
	PTR              []string  `json:"ptr"`
	ForwardConfirmed []string  `json:"forward_confirmed"`
	ASN              string    `json:"asn"`
	Organization     string    `json:"organization"`
	NetworkName      string    `json:"network_name"`
	Registry         string    `json:"registry"`
	Country          string    `json:"country"`
	AddressBlock     string    `json:"address_block"`
	Routes           []string  `json:"routes"`
	Cached           bool      `json:"cached"`
	Errors           []string  `json:"errors"`
	RetrievedAt      time.Time `json:"retrieved_at"`
}

type Analysis struct {
	File       string       `json:"file"`
	Email      *types.Email `json:"email"`
	Verdict    Verdict      `json:"verdict"`
	Match      *Match       `json:"match"`
	Reasons    []string     `json:"reasons"`
	Indicators []*Indicator `json:"indicators"`
}
