package diff_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/diff"
)

func TestCompare_Added(t *testing.T) {
	before := map[string]interface{}{}
	after := map[string]interface{}{"api_key": "abc123"}

	result := diff.Compare("secret/app", before, after)

	if !result.HasChanges() {
		t.Fatal("expected changes, got none")
	}
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added {
		t.Errorf("expected Added, got %s", result.Changes[0].Type)
	}
}

func TestCompare_Removed(t *testing.T) {
	before := map[string]interface{}{"db_pass": "secret"}
	after := map[string]interface{}{}

	result := diff.Compare("secret/db", before, after)

	if result.Changes[0].Type != diff.Removed {
		t.Errorf("expected Removed, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].NewValue != "" {
		t.Errorf("expected empty NewValue for removed key")
	}
}

func TestCompare_Modified(t *testing.T) {
	before := map[string]interface{}{"token": "old"}
	after := map[string]interface{}{"token": "new"}

	result := diff.Compare("secret/token", before, after)

	if result.Changes[0].Type != diff.Modified {
		t.Errorf("expected Modified, got %s", result.Changes[0].Type)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	data := map[string]interface{}{"host": "localhost"}
	result := diff.Compare("secret/host", data, data)

	if result.HasChanges() {
		t.Error("expected no changes for identical maps")
	}
}

func TestCompare_Mixed(t *testing.T) {
	before := map[string]interface{}{"a": "1", "b": "2"}
	after := map[string]interface{}{"a": "1", "c": "3"}

	result := diff.Compare("secret/mixed", before, after)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}

	changeMap := make(map[string]diff.ChangeType)
	for _, c := range result.Changes {
		changeMap[c.Key] = c.Type
	}

	if changeMap["a"] != diff.Unchanged {
		t.Errorf("key 'a' should be unchanged")
	}
	if changeMap["b"] != diff.Removed {
		t.Errorf("key 'b' should be removed")
	}
	if changeMap["c"] != diff.Added {
		t.Errorf("key 'c' should be added")
	}
}
