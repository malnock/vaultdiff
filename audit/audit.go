package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/vaultdiff/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp   time.Time        `json:"timestamp"`
	Environment string           `json:"environment"`
	Path        string           `json:"path"`
	FromVersion int              `json:"from_version"`
	ToVersion   int              `json:"to_version"`
	Changes     []diff.Change    `json:"changes"`
	Summary     Summary          `json:"summary"`
}

// Summary holds aggregate counts for an audit entry.
type Summary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Total    int `json:"total"`
}

// BuildEntry constructs an Entry from a set of diff changes.
func BuildEntry(env, path string, fromVer, toVer int, changes []diff.Change) Entry {
	s := Summary{Total: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		}
	}
	return Entry{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Path:        path,
		FromVersion: fromVer,
		ToVersion:   toVer,
		Changes:     changes,
		Summary:     s,
	}
}

// WriteJSON appends an audit Entry to the given file path as a JSON line.
func WriteJSON(filePath string, entry Entry) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}
