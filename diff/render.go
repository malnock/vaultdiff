package diff

import (
	"fmt"
	"io"
	"strings"
)

// OutputFormat controls how the diff result is rendered.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

// RenderText writes a human-readable diff of the Result to the given writer.
func RenderText(w io.Writer, result Result, maskSecrets bool) {
	fmt.Fprintf(w, "Path: %s\n", result.Path)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	if !result.HasChanges() {
		fmt.Fprintln(w, "  No changes detected.")
		return
	}

	for _, c := range result.Changes {
		if c.Type == Unchanged {
			continue
		}

		oldDisplay := displayValue(c.OldValue, maskSecrets)
		newDisplay := displayValue(c.NewValue, maskSecrets)

		switch c.Type {
		case Added:
			fmt.Fprintf(w, "  + %s: %s\n", c.Key, newDisplay)
		case Removed:
			fmt.Fprintf(w, "  - %s: %s\n", c.Key, oldDisplay)
		case Modified:
			fmt.Fprintf(w, "  ~ %s: %s -> %s\n", c.Key, oldDisplay, newDisplay)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
}

// RenderSummary writes a one-line summary (counts of changes) to the writer.
func RenderSummary(w io.Writer, result Result) {
	var added, removed, modified int
	for _, c := range result.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	fmt.Fprintf(w, "[%s] +%d -%d ~%d\n", result.Path, added, removed, modified)
}

func displayValue(v string, mask bool) string {
	if v == "" {
		return "(empty)"
	}
	if mask {
		return "***"
	}
	return v
}
