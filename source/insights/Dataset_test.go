package insights

import "antispam/structs"
import "antispam/types"
import "bytes"
import "testing"

func TestDatasetRoundTrip(t *testing.T) {

	dataset, err := Load(".")

	if err != nil {
		t.Fatal(err)
	}

	for _, file := range dataset.Files {

		rendered := file.Render()

		if bytes.Equal(rendered, file.Original) == false {
			t.Errorf("JSON round-trip changed formatting of %s", file.RelPath)
		}

	}

	for _, file := range dataset.UnblockableFiles {

		rendered := file.Render()

		if bytes.Equal(rendered, file.Original) == false {
			t.Errorf("JSON round-trip changed formatting of %s", file.RelPath)
		}

	}

	for _, file := range dataset.HostsFiles {

		rendered := file.Render()

		if bytes.Equal(rendered, file.Original) == false {
			t.Errorf("Hosts round-trip changed formatting of %s", file.RelPath)
		}

	}

}

func testEmail(domains []string, ipv4s []string) *types.Email {

	email := types.NewEmail()
	email.MessageID = "test@example.com"
	email.Domains = domains
	email.IPv4s = ipv4s

	return &email

}

func TestDatasetClassify(t *testing.T) {

	dataset, err := Load(".")

	if err != nil {
		t.Fatal(err)
	}

	email := testEmail([]string{"amazon.com"}, []string{"52.94.0.1"})
	verdict, match, reasons := dataset.Classify(email)

	if verdict != structs.VerdictSpam {
		t.Errorf("Expected amazon IP to be spam, got %s", verdict)
	}

	if match == nil || match.Domain != "amazon.com" {
		t.Errorf("Expected amazon.com match")
	}

	if len(reasons) != 1 || reasons[0] != "amazon.com" {
		t.Errorf("Expected domain reason, got %v", reasons)
	}

	email = testEmail([]string{"example.com"}, []string{"209.85.128.1"})
	verdict, match, _ = dataset.Classify(email)

	if verdict != structs.VerdictUnblockable {
		t.Errorf("Expected google IP to be unblockable, got %s", verdict)
	}

	if match == nil || match.Domain != "google.com" {
		t.Errorf("Expected google.com match")
	}

}
