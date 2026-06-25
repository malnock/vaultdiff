package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// PrintEntries reads a JSON-lines audit file and pretty-prints each entry.
func PrintEntries(filePath string, out io.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	index := 0
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return fmt.Errorf("audit: decode entry %d: %w", index, err)
		}
		printEntry(out, e, index+1)
		index++
	}
	if index == 0 {
		fmt.Fprintln(out, "No audit entries found.")
	}
	return nil
}

func printEntry(out io.Writer, e Entry, n int) {
	fmt.Fprintf(out, "── Entry #%d ─────────────────────────────\n", n)
	fmt.Fprintf(out, "  Time:        %s\n", e.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(out, "  Environment: %s\n", e.Environment)
	fmt.Fprintf(out, "  Path:        %s\n", e.Path)
	fmt.Fprintf(out, "  Versions:    v%d → v%d\n", e.FromVersion, e.ToVersion)
	fmt.Fprintf(out, "  Changes:     +%d added  -%d removed  ~%d modified  (total: %d)\n",
		e.Summary.Added, e.Summary.Removed, e.Summary.Modified, e.Summary.Total)
	fmt.Fprintln(out)
}
