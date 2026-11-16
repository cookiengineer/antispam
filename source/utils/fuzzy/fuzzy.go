package fuzzy

import "sort"
import "strings"

var stopwords map[string]bool = map[string]bool{
	"and":            true,
	"the":            true,
	"of":             true,
	"for":            true,
	"inc":            true,
	"llc":            true,
	"ltd":            true,
	"limited":        true,
	"gmbh":           true,
	"mbh":            true,
	"ag":             true,
	"sa":             true,
	"sas":            true,
	"sarl":           true,
	"bv":             true,
	"nv":             true,
	"corp":           true,
	"corporation":    true,
	"company":        true,
	"co":             true,
	"group":          true,
	"holding":        true,
	"holdings":       true,
	"technologies":   true,
	"technology":     true,
	"tech":           true,
	"services":       true,
	"service":        true,
	"solutions":      true,
	"solution":       true,
	"systems":        true,
	"system":         true,
	"network":        true,
	"networks":       true,
	"hosting":        true,
	"host":           true,
	"cloud":          true,
	"internet":       true,
	"online":         true,
	"digital":        true,
	"mail":           true,
	"email":          true,
	"mailer":         true,
	"smtp":           true,
	"mx":             true,
	"media":          true,
	"communications": true,
	"communication":  true,
	"com":            true,
	"net":            true,
	"org":            true,
	"info":           true,
	"biz":            true,
	"io":             true,
	"de":             true,
	"uk":             true,
	"us":             true,
	"eu":             false,
}

type Candidate struct {
	Key    string
	Tokens []string
	Text   string
}

type Match struct {
	Key    string
	Score  float64
	Reason string
}

func normalizeParts(value string) []string {

	value = strings.ToLower(value)

	tmp := strings.FieldsFunc(value, func(character rune) bool {
		return !((character >= 'a' && character <= 'z') || (character >= '0' && character <= '9'))
	})

	result := make([]string, 0)

	for _, part := range tmp {

		part = strings.TrimSpace(part)

		if len(part) < 2 {
			continue
		}

		if stopwords[part] == true {
			continue
		}

		result = append(result, part)

	}

	return result

}

// Normalize returns a lower-case token string for comparisons.
func Normalize(value string) string {
	return strings.Join(normalizeParts(value), " ")
}

// Tokens returns the significant normalized tokens of a value.
func Tokens(value string) []string {
	return normalizeParts(value)
}

// NewCandidate creates a match candidate from a key and a textual
// representation of the domain, aliases, and filename.
func NewCandidate(key string, text string) Candidate {

	return Candidate{
		Key:    key,
		Tokens: unique(Tokens(text)),
		Text:   text,
	}

}

func unique(values []string) []string {

	result := make([]string, 0)
	seen := make(map[string]bool)

	for _, value := range values {

		if seen[value] == false {
			seen[value] = true
			result = append(result, value)
		}

	}

	return result

}

func contains(values []string, value string) bool {

	for _, item := range values {

		if item == value {
			return true
		}

	}

	return false

}

func prefixOverlap(a string, b string) bool {

	if a == b {
		return true
	}

	if len(a) >= 4 && len(b) >= 4 {

		if strings.HasPrefix(a, b) == true || strings.HasPrefix(b, a) == true {
			return true
		}

	}

	return false

}

func ratio(a string, b string) float64 {

	if a == b {
		return 1.0
	}

	max_length := len(a)

	if len(b) > max_length {
		max_length = len(b)
	}

	if max_length == 0 {
		return 1.0
	}

	distance := levenshtein(a, b)

	return 1.0 - float64(distance)/float64(max_length)

}

func levenshtein(a string, b string) int {

	if a == "" {
		return len(b)
	}

	if b == "" {
		return len(a)
	}

	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)

	for j := 0; j <= len(b); j++ {
		previous[j] = j
	}

	for i := 1; i <= len(a); i++ {

		current[0] = i

		for j := 1; j <= len(b); j++ {

			cost := 1

			if a[i-1] == b[j-1] {
				cost = 0
			}

			deletion := previous[j] + 1
			insertion := current[j-1] + 1
			substitution := previous[j-1] + cost

			current[j] = deletion

			if insertion < current[j] {
				current[j] = insertion
			}

			if substitution < current[j] {
				current[j] = substitution
			}

		}

		previous, current = current, previous

	}

	return previous[len(b)]

}

// Score compares a query string against a candidate and returns a score
// between 0 and 1 plus a human-readable reason.
func Score(query string, candidate Candidate) (float64, string) {

	query_tokens := unique(Tokens(query))

	if len(query_tokens) == 0 || len(candidate.Tokens) == 0 {
		return 0.0, ""
	}

	query_key := strings.Join(query_tokens, " ")
	candidate_key := strings.Join(candidate.Tokens, " ")

	if query_key == candidate_key {
		return 1.0, "exact token match"
	}

	shared := make([]string, 0)
	prefix_shared := make([]string, 0)

	for _, query_token := range query_tokens {

		if contains(candidate.Tokens, query_token) == true {

			shared = append(shared, query_token)

		} else {

			for _, candidate_token := range candidate.Tokens {

				if prefixOverlap(query_token, candidate_token) == true && query_token != candidate_token {
					prefix_shared = append(prefix_shared, candidate_token)
					break
				}

			}

		}

	}

	score := 0.0
	reason := ""

	if len(shared) > 0 {

		// A shared significant token is a strong hint. Legal suffixes and
		// generic words were already removed above.
		score = 0.72 + 0.08*float64(len(shared))

		if len(shared) == len(candidate.Tokens) {
			score += 0.08
		}

		if len(shared) == len(query_tokens) {
			score += 0.05
		}

		if score > 0.95 {
			score = 0.95
		}

		reason = "token match: " + strings.Join(shared, ", ")

	}

	for _, candidate_token := range candidate.Tokens {

		value := ratio(query_key, candidate_token)

		if value >= 0.86 && value > score {
			score = value * 0.95
			reason = "similar token: " + candidate_token
		}

	}

	if len(prefix_shared) > 0 && score < 0.85 {

		score = 0.85
		reason = "prefix match: " + strings.Join(prefix_shared, ", ")

	}

	if score > 1.0 {
		score = 1.0
	}

	return score, reason

}

// MatchAll returns matches sorted by descending score.
func MatchAll(query string, candidates []Candidate, minimum float64) []Match {

	result := make([]Match, 0)

	for _, candidate := range candidates {

		score, reason := Score(query, candidate)

		if score >= minimum {
			result = append(result, Match{
				Key:    candidate.Key,
				Score:  score,
				Reason: reason,
			})
		}

	}

	sort.Slice(result, func(a int, b int) bool {

		if result[a].Score == result[b].Score {
			return result[a].Key < result[b].Key
		}

		return result[a].Score > result[b].Score

	})

	return result

}
