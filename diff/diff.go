package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType represents the type of change detected between two secret versions.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference between two secret snapshots.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff result between two secret snapshots.
type Result struct {
	Path    string
	Changes []Change
}

// HasChanges returns true if any non-unchanged entries exist.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Compare diffs two secret data maps (string->string) and returns a Result.
func Compare(path string, before, after map[string]interface{}) Result {
	result := Result{Path: path}

	allKeys := make(map[string]struct{})
	for k := range before {
		allKeys[k] = struct{}{}
	}
	for k := range after {
		allKeys[k] = struct{}{}
	}

	sortedKeys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		oldVal := fmt.Sprintf("%v", before[key])
		newVal := fmt.Sprintf("%v", after[key])

		_, inBefore := before[key]
		_, inAfter := after[key]

		var ct ChangeType
		switch {
		case inBefore && !inAfter:
			ct = Removed
			newVal = ""
		case !inBefore && inAfter:
			ct = Added
			oldVal = ""
		case strings.TrimSpace(oldVal) != strings.TrimSpace(newVal):
			ct = Modified
		default:
			ct = Unchanged
		}

		result.Changes = append(result.Changes, Change{
			Key:      key,
			Type:     ct,
			OldValue: oldVal,
			NewValue: newVal,
		})
	}

	return result
}
