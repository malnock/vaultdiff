package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/diff"
)

func TestRenderText_NoChanges(t *testing.T) {
	result := diff.Result{Path: "secret/noop"}
	var buf bytes.Buffer
	diff.RenderText(&buf, result, false)

	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestRenderText_ShowsAdded(t *testing.T) {
	result := diff.Result{
		Path: "secret/app",
		Changes: []diff.Change{
			{Key: "api_key", Type: diff.Added, NewValue: "xyz"},
		},
	}
	var buf bytes.Buffer
	diff.RenderText(&buf, result, false)

	if !strings.Contains(buf.String(), "+ api_key: xyz") {
		t.Errorf("expected added key in output, got: %s", buf.String())
	}
}

func TestRenderText_MasksSecrets(t *testing.T) {
	result := diff.Result{
		Path: "secret/db",
		Changes: []diff.Change{
			{Key: "password", Type: diff.Modified, OldValue: "old", NewValue: "new"},
		},
	}
	var buf bytes.Buffer
	diff.RenderText(&buf, result, true)

	output := buf.String()
	if strings.Contains(output, "old") || strings.Contains(output, "new") {
		t.Errorf("secret values should be masked, got: %s", output)
	}
	if !strings.Contains(output, "***") {
		t.Errorf("expected masked value '***', got: %s", output)
	}
}

func TestRenderSummary(t *testing.T) {
	result := diff.Result{
		Path: "secret/summary",
		Changes: []diff.Change{
			{Key: "a", Type: diff.Added},
			{Key: "b", Type: diff.Removed},
			{Key: "c", Type: diff.Modified},
			{Key: "d", Type: diff.Unchanged},
		},
	}
	var buf bytes.Buffer
	diff.RenderSummary(&buf, result)

	output := buf.String()
	if !strings.Contains(output, "+1") {
		t.Errorf("expected +1 added, got: %s", output)
	}
	if !strings.Contains(output, "-1") {
		t.Errorf("expected -1 removed, got: %s", output)
	}
	if !strings.Contains(output, "~1") {
		t.Errorf("expected ~1 modified, got: %s", output)
	}
}
