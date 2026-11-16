package fuzzy

import "testing"

func TestScore(t *testing.T) {

	candidates := []Candidate{
		NewCandidate("amazon", "amazon.com amazon.cn amazon.de amazonaws.com"),
		NewCandidate("mailgun", "mailgun.net mailgun.com substack.com"),
		NewCandidate("microsoft", "microsoft.com outlook.com hotmail.com"),
	}

	matches := MatchAll("Amazon EU S.a.r.l.", candidates, 0.55)

	if len(matches) == 0 || matches[0].Key != "amazon" {
		t.Errorf("Expected amazon match, got %v", matches)
	}

	matches = MatchAll("Mailgun Technologies Inc.", candidates, 0.55)

	if len(matches) == 0 || matches[0].Key != "mailgun" {
		t.Errorf("Expected mailgun match, got %v", matches)
	}

	matches = MatchAll("mail-swedencentralazon11023086.outbound.protection.outlook.com", candidates, 0.55)

	if len(matches) == 0 || matches[0].Key != "microsoft" {
		t.Errorf("Expected microsoft/outlook match, got %v", matches)
	}

	if Normalize("Amazon EU S.a.r.l.") != "amazon eu" {
		t.Errorf("Expected normalized amazon eu, got %s", Normalize("Amazon EU S.a.r.l."))
	}

}
