package audit_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/vaultdiff/audit"
	"github.com/vaultdiff/diff"
)

func sampleChanges() []diff.Change {
	return []diff.Change{
		{Key: "DB_PASS", Type: diff.Added, NewValue: "secret123"},
		{Key: "API_KEY", Type: diff.Removed, OldValue: "old"},
		{Key: "HOST", Type: diff.Modified, OldValue: "localhost", NewValue: "prod.example.com"},
		{Key: "PORT", Type: diff.Unchanged, OldValue: "5432", NewValue: "5432"},
	}
}

func TestBuildEntry_SummaryCounts(t *testing.T) {
	changes := sampleChanges()
	entry := audit.BuildEntry("production", "secret/app", 1, 2, changes)

	if entry.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", entry.Summary.Added)
	}
	if entry.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", entry.Summary.Removed)
	}
	if entry.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", entry.Summary.Modified)
	}
	if entry.Summary.Total != 4 {
		t.Errorf("expected total 4, got %d", entry.Summary.Total)
	}
}

func TestBuildEntry_Fields(t *testing.T) {
	entry := audit.BuildEntry("staging", "secret/db", 3, 4, nil)

	if entry.Environment != "staging" {
		t.Errorf("unexpected environment: %s", entry.Environment)
	}
	if entry.Path != "secret/db" {
		t.Errorf("unexpected path: %s", entry.Path)
	}
	if entry.FromVersion != 3 || entry.ToVersion != 4 {
		t.Errorf("unexpected versions: %d -> %d", entry.FromVersion, entry.ToVersion)
	}
	if entry.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestWriteJSON_CreatesFileAndAppendsLines(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	e1 := audit.BuildEntry("prod", "secret/app", 1, 2, sampleChanges())
	e2 := audit.BuildEntry("prod", "secret/app", 2, 3, nil)

	if err := audit.WriteJSON(tmpFile.Name(), e1); err != nil {
		t.Fatalf("WriteJSON e1: %v", err)
	}
	if err := audit.WriteJSON(tmpFile.Name(), e2); err != nil {
		t.Fatalf("WriteJSON e2: %v", err)
	}

	data, _ := os.ReadFile(tmpFile.Name())
	lines := splitLines(data)
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSON lines, got %d", len(lines))
	}

	var decoded audit.Entry
	if err := json.Unmarshal([]byte(lines[0]), &decoded); err != nil {
		t.Fatalf("unmarshal line 1: %v", err)
	}
	if decoded.Environment != "prod" {
		t.Errorf("decoded environment mismatch: %s", decoded.Environment)
	}
}

func splitLines(data []byte) []string {
	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' && i > start {
			lines = append(lines, string(data[start:i]))
			start = i + 1
		}
	}
	return lines
}
