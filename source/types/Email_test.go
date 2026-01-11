package types

import _ "embed"
import "strings"
import "testing"

//go:embed Email_test_google_business.eml
var sample_google_business []byte

//go:embed Email_test_zimbra.eml
var sample_zimbra []byte

//go:embed Email_test_rfc2047.eml
var sample_rfc2047 []byte

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
			t.Errorf("Expected Message to be parsed correctly")
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

	t.Run("ParseEmail(zimbra)", func(t *testing.T) {

		email := ParseEmail(sample_zimbra)

		if email.MessageID != "938764576.194494.1722920907974.JavaMail.zimbra@saludquillota.cl" {
			t.Errorf("Expected MessageID %s to be %s", email.MessageID, "938764576.194494.1722920907974.JavaMail.zimbra@saludquillota.cl")
		}

		if email.Boundary != "" {
			t.Errorf("Expected Boundary %s to be empty", email.Boundary)
		}

		if email.From != "gladys.vargas@saludquillota.cl" {
			t.Errorf("Expected From %s to be %s", email.From, "gladys.vargas@saludquillota.cl")
		}

		if email.To != "johndoe@example.com" {
			t.Errorf("Expected To %s to be %s", email.To, "johndoe@example.com")
		}

		if email.Subject != "Project Financing" {
			t.Errorf("Expected Subject \"%s\" to be \"%s\"", email.Subject, "Project Financing")
		}

		if !strings.HasPrefix(email.Message, "Greetings,") || !strings.HasSuffix(email.Message, "United Arab Emirates") {
			t.Errorf("Expected Message to be parsed correctly")
		}

		if email.Date.Format("2006-01-02 15:04:05") != "2024-08-06 01:08:27" {
			t.Errorf("Expected Date \"%s\" to be \"%s\"", email.Date, "2024-08-06 01:08:27")
		}

		if len(email.Domains) == 2 {

			if email.Domains[0] != "relaystgo2.tchile.com" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[0], "relaystgo2.tchile.com")
			}

			if email.Domains[1] != "saludquillota.cl" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[1], "saludquillota.cl")
			}

		} else {
			t.Errorf("Expected %d Domains to be %d", len(email.Domains), 2)
		}

		if len(email.IPv4s) == 2 {

			if email.IPv4s[0] != "131.108.210.19" {
				t.Errorf("Expected IPv4 \"%s\" to be \"%s\"", email.IPv4s[0], "131.108.210.19")
			}

			if email.IPv4s[1] != "127.0.0.1" {
				t.Errorf("Expected IPv4 \"%s\" to be \"%s\"", email.IPv4s[1], "127.0.0.1")
			}

		} else {
			t.Errorf("Expected %d IPv4s to be %d", len(email.IPv4s), 2)
		}

		if len(email.IPv6s) == 0 {
			// Do Nothing
		} else {
			t.Errorf("Expected %d IPv6s to be %d", len(email.IPv6s), 0)
		}

	})

	t.Run("ParseEmail(rfc2047)", func(t *testing.T) {

		email := ParseEmail(sample_rfc2047)

		if email.MessageID != "1579173458.453595.1738289482391.JavaMail.zimbra@consultoresdeempresas.com" {
			t.Errorf("Expected MessageID %s to be %s", email.MessageID, "1579173458.453595.1738289482391.JavaMail.zimbra@consultoresdeempresas.com")
		}

		if email.Boundary != "=_a32fe805-f8f3-491f-bb72-381c11415531" {
			t.Errorf("Expected Boundary %s to be %s", email.Boundary, "=_a32fe805-f8f3-491f-bb72-381c11415531")
		}

		if email.From != "garganaraz@consultoresdeempresas.com" {
			t.Errorf("Expected From %s to be %s", email.From, "garganaraz@consultoresdeempresas.com")
		}

		if email.To != "johndoe@example.com" {
			t.Errorf("Expected To %s to be %s", email.To, "johndoe@example.com")
		}

		if email.Subject != "€1,800.000 EUR Spende für Sie" {
			t.Errorf("Expected Subject \"%s\" to be \"%s\"", email.Subject, "€1,800.000 EUR Spende für Sie")
		}

		if !strings.HasPrefix(email.Message, "Hallo ,") || !strings.HasSuffix(email.Message, "für weitere Informationen:  doug64411@gmail.com") {
			t.Errorf("Expected Message to be parsed correctly")
		}

		if email.Date.Format("2006-01-02 15:04:05") != "2025-01-31 02:11:22" {
			t.Errorf("Expected Date \"%s\" to be \"%s\"", email.Date, "2025-01-31 02:11:22")
		}


		if len(email.Domains) == 2 {

			if email.Domains[0] != "mail.consultoresdeempresas.com" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[0], "mail.consultoresdeempresas.com")
			}

			if email.Domains[1] != "consultoresdeempresas.com" {
				t.Errorf("Expected Domain \"%s\" to be \"%s\"", email.Domains[1], "consultoresdeempresas.com")
			}

		} else {
			t.Errorf("Expected %d Domains to be %d", len(email.Domains), 2)
		}

		if len(email.IPv4s) == 3 {

			if email.IPv4s[0] != "190.221.161.226" {
				t.Errorf("Expected IPv4 \"%s\" to be \"%s\"", email.IPv4s[0], "190.221.161.226")
			}

			if email.IPv4s[1] != "127.0.0.1" {
				t.Errorf("Expected IPv4 \"%s\" to be \"%s\"", email.IPv4s[1], "127.0.0.1")
			}

			if email.IPv4s[2] != "172.16.10.102" {
				t.Errorf("Expected IPv4 \"%s\" to be \"%s\"", email.IPv4s[2], "172.16.10.102")
			}

		} else {
			t.Errorf("Expected %d IPv4s to be %d", len(email.IPv4s), 3)
		}

		if len(email.IPv6s) == 0 {
			// Do Nothing
		} else {
			t.Errorf("Expected %d IPv6s to be %d", len(email.IPv6s), 0)
		}

	})

}
