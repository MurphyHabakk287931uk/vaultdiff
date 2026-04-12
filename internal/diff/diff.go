package diff

import (
	"fmt"
	"sort"

	"github.com/user/vaultdiff/internal/redact"
)

// ChangeType represents the kind of change between two secret values.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Entry represents a single key diff between two secret maps.
type Entry struct {
	Key      string
	Change   ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two secret maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if any entry is not Unchanged.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Change != Unchanged {
			return true
		}
	}
	return false
}

// Compare diffs two flat secret maps and returns a Result.
func Compare(src, dst map[string]string) *Result {
	keys := unionKeys(src, dst)
	sort.Strings(keys)

	var entries []Entry
	for _, k := range keys {
		srcVal, inSrc := src[k]
		dstVal, inDst := dst[k]

		switch {
		case inSrc && !inDst:
			entries = append(entries, Entry{Key: k, Change: Removed, OldValue: srcVal})
		case !inSrc && inDst:
			entries = append(entries, Entry{Key: k, Change: Added, NewValue: dstVal})
		case srcVal != dstVal:
			entries = append(entries, Entry{Key: k, Change: Modified, OldValue: srcVal, NewValue: dstVal})
		default:
			entries = append(entries, Entry{Key: k, Change: Unchanged, OldValue: srcVal, NewValue: dstVal})
		}
	}
	return &Result{Entries: entries}
}

// Format renders a Result as a human-readable diff string with redaction applied.
func Format(r *Result, mode redact.Mode) string {
	var out string
	for _, e := range r.Entries {
		switch e.Change {
		case Added:
			out += fmt.Sprintf("+ %s = %s\n", e.Key, redact.Apply(e.NewValue, mode))
		case Removed:
			out += fmt.Sprintf("- %s = %s\n", e.Key, redact.Apply(e.OldValue, mode))
		case Modified:
			out += fmt.Sprintf("~ %s: %s -> %s\n", e.Key, redact.Apply(e.OldValue, mode), redact.Apply(e.NewValue, mode))
		case Unchanged:
			out += fmt.Sprintf("  %s = %s\n", e.Key, redact.Apply(e.OldValue, mode))
		}
	}
	return out
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
