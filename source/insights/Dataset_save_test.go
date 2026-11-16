package insights

import "os"
import "path/filepath"
import "strings"
import "testing"

func TestDatasetSaveTemp(t *testing.T) {

	root := t.TempDir()

	os.MkdirAll(filepath.Join(root, "phishers"), 0755)
	os.MkdirAll(filepath.Join(root, "spammers", "unblockable"), 0755)
	os.MkdirAll(filepath.Join(root, "hosts"), 0755)

	json_path := filepath.Join(root, "spammers", "example.json")
	json_data := "[{\n\t\"domain\": \"example.com\",\n\t\"networks\": [\n\t\t\"203.0.113.0/24\"\n\t]\n}]\n"

	err := os.WriteFile(json_path, []byte(json_data), 0644)

	if err != nil {
		t.Fatal(err)
	}

	hosts_path := filepath.Join(root, "hosts", "spammers.hosts")
	hosts_data := "spam.example\n"

	err = os.WriteFile(hosts_path, []byte(hosts_data), 0644)

	if err != nil {
		t.Fatal(err)
	}

	dataset, err := Load(root)

	if err != nil {
		t.Fatal(err)
	}

	entry := dataset.FindEntry("example.com")

	if entry == nil {
		t.Fatal("example.com missing")
	}

	_, err = dataset.AddNetwork(entry, "203.0.114.0/24")

	if err != nil {
		t.Fatal(err)
	}

	if dataset.AddHost(dataset.HostsFiles[0], "other.example") == false {
		t.Fatal("expected host to be staged")
	}

	if len(dataset.Diff()) != 2 {
		t.Fatalf("expected two diffs, got %d", len(dataset.Diff()))
	}

	changed, err := dataset.Save()

	if err != nil {
		t.Fatal(err)
	}

	if len(changed) != 2 {
		t.Errorf("expected two changed files, got %v", changed)
	}

	buffer, err := os.ReadFile(json_path)

	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(buffer), "203.0.114.0/24") == false {
		t.Errorf("expected new network in saved JSON")
	}

	buffer, err = os.ReadFile(hosts_path)

	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(buffer), "other.example") == false {
		t.Errorf("expected new host in saved hosts file")
	}

	if len(dataset.PendingFiles()) != 0 {
		t.Errorf("expected no pending files after save")
	}

}
