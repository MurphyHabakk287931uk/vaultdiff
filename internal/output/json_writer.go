package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/vaultdiff/internal/diff"
)

// JSONResult is the top-level JSON output structure.
type JSONResult struct {
	Src     string        `json:"src"`
	Dst     string        `json:"dst"`
	Changes []JSONChange  `json:"changes"`
	Summary JSONSummary   `json:"summary"`
}

// JSONChange represents a single diff entry in JSON output.
type JSONChange struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	Src    string `json:"src,omitempty"`
	Dst    string `json:"dst,omitempty"`
}

// JSONSummary holds counts of each change type.
type JSONSummary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Unchanged int `json:"unchanged"`
}

// WriteJSON encodes the diff results as JSON to w.
func WriteJSON(w io.Writer, src, dst string, entries []diff.Entry) error {
	result := JSONResult{
		Src:     src,
		Dst:     dst,
		Changes: make([]JSONChange, 0, len(entries)),
	}
	for _, e := range entries {
		change := JSONChange{
			Key:    e.Key,
			Status: statusString(e.Status),
			Src:    e.SrcValue,
			Dst:    e.DstValue,
		}
		result.Changes = append(result.Changes, change)
		switch e.Status {
		case diff.Added:
			result.Summary.Added++
		case diff.Removed:
			result.Summary.Removed++
		case diff.Modified:
			result.Summary.Modified++
		case diff.Unchanged:
			result.Summary.Unchanged++
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("encoding JSON output: %w", err)
	}
	return nil
}

func statusString(s diff.Status) string {
	switch s {
	case diff.Added:
		return "added"
	case diff.Removed:
		return "removed"
	case diff.Modified:
		return "modified"
	default:
		return "unchanged"
	}
}
