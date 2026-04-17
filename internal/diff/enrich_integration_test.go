package diff_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/redact"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// TestRunner_EnrichClientAddsKeysToComparison verifies that when one side is
// wrapped with EnrichClient the injected keys appear in the diff output.
func TestRunner_EnrichClientAddsKeysToComparison(t *testing.T) {
	srcData := map[string]map[string]string{
		"secret/app": {"password": "abc"},
	}
	dstData := map[string]map[string]string{
		"secret/app": {"password": "abc"},
	}

	src := vault.NewEnrichClient(
		vault.NewMockClient(srcData),
		map[string]string{"env": "staging"},
	)
	dst := vault.NewMockClient(dstData)

	r := diff.NewRunner(src, dst, redact.ModeNone, false)
	results, err := r.Run(nil, "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var found bool
	for _, entry := range results {
		if entry.Key == "env" {
			found = true
			if entry.Status != diff.StatusRemoved {
				t.Errorf("expected env to be removed (only in src), got %v", entry.Status)
			}
		}
	}
	if !found {
		t.Error("expected 'env' key in diff results from enriched src")
	}
}

// TestRunner_BothEnrichedSameKeyUnchanged checks that when both sides are
// enriched with the same key/value pair the key is unchanged in the diff.
func TestRunner_BothEnrichedSameKeyUnchanged(t *testing.T) {
	base := map[string]map[string]string{
		"secret/app": {"token": "xyz"},
	}
	extra := map[string]string{"region": "us-east-1"}

	src := vault.NewEnrichClient(vault.NewMockClient(base), extra)
	dst := vault.NewEnrichClient(vault.NewMockClient(base), extra)

	r := diff.NewRunner(src, dst, redact.ModeNone, true)
	results, err := r.Run(nil, "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no diff entries when both sides identical, got %d", len(results))
	}
}
