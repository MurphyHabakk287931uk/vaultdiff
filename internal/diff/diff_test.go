package diff_test

import (
	"strings"
	"testing"

	"github.com/user/vaultdiff/internal/diff"
	"github.com/user/vaultdiff/internal/redact"
)

func TestCompare_AllUnchanged(t *testing.T) {
	src := map[string]string{"key": "value"}
	dst := map[string]string{"key": "value"}
	r := diff.Compare(src, dst)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
	if len(r.Entries) != 1 || r.Entries[0].Change != diff.Unchanged {
		t.Errorf("expected Unchanged, got %v", r.Entries)
	}
}

func TestCompare_Added(t *testing.T) {
	src := map[string]string{}
	dst := map[string]string{"new_key": "new_val"}
	r := diff.Compare(src, dst)
	if !r.HasChanges() {
		t.Error("expected changes")
	}
	if r.Entries[0].Change != diff.Added {
		t.Errorf("expected Added, got %v", r.Entries[0].Change)
	}
}

func TestCompare_Removed(t *testing.T) {
	src := map[string]string{"old_key": "old_val"}
	dst := map[string]string{}
	r := diff.Compare(src, dst)
	if r.Entries[0].Change != diff.Removed {
		t.Errorf("expected Removed, got %v", r.Entries[0].Change)
	}
}

func TestCompare_Modified(t *testing.T) {
	src := map[string]string{"k": "v1"}
	dst := map[string]string{"k": "v2"}
	r := diff.Compare(src, dst)
	e := r.Entries[0]
	if e.Change != diff.Modified {
		t.Errorf("expected Modified, got %v", e.Change)
	}
	if e.OldValue != "v1" || e.NewValue != "v2" {
		t.Errorf("unexpected values: %v", e)
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	src := map[string]string{"z": "1", "a": "2", "m": "3"}
	dst := map[string]string{"z": "1", "a": "2", "m": "3"}
	r := diff.Compare(src, dst)
	if r.Entries[0].Key != "a" || r.Entries[1].Key != "m" || r.Entries[2].Key != "z" {
		t.Errorf("keys not sorted: %v", r.Entries)
	}
}

func TestFormat_RedactMode(t *testing.T) {
	src := map[string]string{"secret": "plaintext"}
	dst := map[string]string{"secret": "plaintext"}
	r := diff.Compare(src, dst)
	out := diff.Format(r, redact.ModeRedact)
	if strings.Contains(out, "plaintext") {
		t.Errorf("expected redacted output, got: %s", out)
	}
}

func TestFormat_NoneMode(t *testing.T) {
	src := map[string]string{"k": "v"}
	dst := map[string]string{"k": "v2"}
	r := diff.Compare(src, dst)
	out := diff.Format(r, redact.ModeNone)
	if !strings.Contains(out, "v") || !strings.Contains(out, "v2") {
		t.Errorf("expected plain output, got: %s", out)
	}
}

func TestFormat_AddedPrefix(t *testing.T) {
	r := diff.Compare(map[string]string{}, map[string]string{"x": "y"})
	out := diff.Format(r, redact.ModeNone)
	if !strings.HasPrefix(out, "+") {
		t.Errorf("expected '+' prefix for added, got: %s", out)
	}
}
