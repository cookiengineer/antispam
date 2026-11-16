package actions

import "antispam/insights"
import "os"
import "path/filepath"
import "testing"

func TestAnalyzeMails(t *testing.T) {

	dataset, err := insights.Load("../insights")

	if err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir("../../mails")

	if err != nil {
		t.Skip("mails folder not available")
	}

	count := 0

	for _, entry := range entries {

		if entry.IsDir() == true {
			continue
		}

		if entry.Name()[0] == '.' {
			continue
		}

		path := filepath.Join("../../mails", entry.Name())
		analysis, err := Analyze(path, dataset)

		if err != nil {
			continue
		}

		if analysis.Email == nil {
			t.Errorf("Expected parsed email for %s", path)
		}

		if len(analysis.Indicators) == 0 {
			t.Errorf("Expected indicators for %s", path)
		}

		count++

	}

	if count == 0 {
		t.Skip("no mails available")
	}

}
