package types

import _ "embed"
import "strings"
import "testing"

//go:embed Email_test_google_business.eml
var sample_google_business []byte

func TestParseEmail(t *testing.T) {

	t.Run("ParseEmail(google business)", func(t *testing.T) {

		email := ParseEmail(sample_google_business)

		if email.MessageID != "CAOmPQAPn1pW+MdO7Vh89+dADDohCQzNkxVxY8wZfe8-EcTmH8A@mail.gmail.com" {
			t.Errorf("Expected MessageID %s to be %s", email.MessageID, "CAOmPQAPn1pW+MdO7Vh89+dADDohCQzNkxVxY8wZfe8-EcTmH8A@mail.gmail.com")
		}

		if email.Boundary != "0000000000004699430629eb9f1a" {
			t.Errorf("Expected Boundary %s to be %s", email.Boundary, "0000000000004699430629eb9f1a")
		}

		if email.From != "brandon.p@hello.tempolabs.ai" {
			t.Errorf("Expected From %s to be %s", email.From, "brandon.p@hello.tempolabs.ai")
		}

		if email.To != "johndoe@example.com" {
			t.Errorf("Expected To %s to be %s", email.To, "johndoe@example.com")
		}

		if email.Subject != "Reducing Build Time" {
			t.Errorf("Expected Subject \"%s\" to be \"%s\"", email.Subject, "Reducing Build Time")
		}

		if !strings.HasPrefix(email.Message, "Hi Chris,") || !strings.HasSuffix(email.Message, "We respect your decision.") {
			t.Errorf("Expected Message to be text/plain message")
		}

		if email.Date.Format("2006-01-02 15:04:05") != "2024-12-23 00:17:06" {
			t.Errorf("Expected Date \"%s\" to be \"%s\"", email.Date, "2024-12-23 00:17:06")
		}

		if len(email.Domains) == 4 {

			if email.Domains[0] != "mail.tholian.network" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[0], "mail.tholian.network")
			}

			if email.Domains[1] != "mail-vk1-f181.google.com" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[1], "mail-vk1-f181.google.com")
			}

			if email.Domains[2] != "gmailapi.google.com" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[2], "gmailapi.google.com")
			}

			if email.Domains[3] != "hello.tempolabs.ai" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[3], "hello.tempolabs.ai")
			}

		} else {
			t.Errorf("Expected %d Domains to be %d", len(email.Domains), 4)
		}

		if len(email.IPv4s) == 1 {

			if email.IPv4s[0] != "209.85.221.181" {
				t.Errorf("Expected IPv4 \"%s\" to be \"%s\"", email.IPv4s[0], "209.85.221.181")
			}

		} else {
			t.Errorf("Expected %d IPv4s to be %d", len(email.IPv4s), 1)
		}

		if len(email.IPv6s) == 0 {
			// Do Nothing
		} else {
			t.Errorf("Expected %d IPv6s to be %d", len(email.IPv6s), 0)
		}

	})

}
