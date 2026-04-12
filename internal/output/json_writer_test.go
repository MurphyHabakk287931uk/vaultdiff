package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/vaultdiff/internal/diff"
)

func TestWriteJSON_BasicOutput(t *testing.T) {
	entries := []diff.Entry{
		{Key: "API_KEY", Status: diff.Added, SrcValue: "", DstValue: "newval"},
		{Key: "DB_PASS", Status: diff.Removed, SrcValue: "old", DstValue: ""},
		{Key: "HOST", Status: diff.Modified, SrcValue: "a", DstValue: "b"},
		{Key: "PORT", Status: diff.Unchanged, SrcValue: "8080", DstValue: "8080"},
	}
	var buf bytes.Buffer
	err := WriteJSON(&buf, "secret/src", "secret/dst", entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result JSONResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result.Src != "secret/src" {
		t.Errorf("src mismatch: got %q", result.Src)
	}
	if result.Dst != "secret/dst" {
		t.Errorf("dst mismatch: got %q", result.Dst)
	}
	if len(result.Changes) != 4 {
		t.Errorf("expected 4 changes, got %d", len(result.Changes))
	}
}

func TestWriteJSON_Summary(t *testing.T) {
	entries := []diff.Entry{
		{Key: "A", Status: diff.Added},
		{Key: "B", Status: diff.Added},
		{Key: "C", Status: diff.Removed},
		{Key: "D", Status: diff.Modified},
		{Key: "E", Status: diff.Unchanged},
	}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, "src", "dst", entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result JSONResult
	json.Unmarshal(buf.Bytes(), &result)
	if result.Summary.Added != 2 {
		t.Errorf("added: want 2, got %d", result.Summary.Added)
	}
	if result.Summary.Removed != 1 {
		t.Errorf("removed: want 1, got %d", result.Summary.Removed)
	}
	if result.Summary.Modified != 1 {
		t.Errorf("modified: want 1, got %d", result.Summary.Modified)
	}
	if result.Summary.Unchanged != 1 {
		t.Errorf("unchanged: want 1, got %d", result.Summary.Unchanged)
	}
}

func TestStatusString(t *testing.T) {
	cases := map[diff.Status]string{
		diff.Added:     "added",
		diff.Removed:   "removed",
		diff.Modified:  "modified",
		diff.Unchanged: "unchanged",
	}
	for status, want := range cases {
		if got := statusString(status); got != want {
			t.Errorf("statusString(%v): got %q, want %q", status, got, want)
		}
	}
}
