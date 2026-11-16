package insights

import "strings"
import "testing"

func TestDatasetEditDiff(t *testing.T) {

	dataset, err := Load(".")

	if err != nil {
		t.Fatal(err)
	}

	entry := dataset.FindEntry("amazon.com")

	if entry == nil {
		t.Fatal("amazon.com entry missing")
	}

	added, err := dataset.AddNetwork(entry, "203.0.113.0/24")

	if err != nil {
		t.Fatal(err)
	}

	if len(added) != 1 || added[0] != "203.0.113.0/24" {
		t.Errorf("Expected one added network, got %v", added)
	}

	if entry.HasNetwork("203.0.113.42") == false {
		t.Errorf("Expected new network to cover test IP")
	}

	diffs := dataset.Diff()

	if len(diffs) != 1 {
		t.Fatalf("Expected one diff, got %d", len(diffs))
	}

	if strings.Contains(diffs[0].Diff, "+") == false {
		t.Errorf("Expected added lines in diff")
	}

	added, err = dataset.AddNetwork(entry, "203.0.113.0/25")

	if err != nil {
		t.Fatal(err)
	}

	if len(added) != 0 {
		t.Errorf("Expected covered network to be skipped, got %v", added)
	}

}
