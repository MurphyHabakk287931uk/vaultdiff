package diff_test

import (
	"context"
	"testing"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/redact"
	"github.com/your-org/vaultdiff/internal/vault"
)

// TestRunner_RejectsOutOfScopeSource verifies that when the source client
// is wrapped with a ScopeClient, an out-of-scope read surfaces as an error
// through the Runner and is not silently swallowed.
func TestRunner_RejectsOutOfScopeSource(t *testing.T) {
	srcData := map[string]map[string]string{
		"secret/prod/db": {"pass": "prod-pass"},
	}
	dstData := map[string]map[string]string{
		"secret/staging/db": {"pass": "stage-pass"},
	}

	srcMock := vault.NewMockClient(srcData)
	dstMock := vault.NewMockClient(dstData)

	// Restrict source to prod only.
	scopedSrc := vault.NewScopeClient(srcMock, "secret/prod")

	r := diff.NewRunner(diff.RunnerConfig{
		Src:        scopedSrc,
		Dst:        dstMock,
		SrcPath:    "secret/staging/db", // intentionally out of scope
		DstPath:    "secret/staging/db",
		RedactMode: redact.ModeNone,
		ShowAll:    true,
	})

	_, err := r.Run(context.Background())
	if err == nil {
		t.Fatal("expected error when source path is out of scope")
	}
}

// TestRunner_ScopedClientsCompareSamePath confirms that two scoped clients
// pointing at the same allowed path produce a normal diff result.
func TestRunner_ScopedClientsCompareSamePath(t *testing.T) {
	srcData := map[string]map[string]string{
		"secret/prod/app": {"key": "v1"},
	}
	dstData := map[string]map[string]string{
		"secret/prod/app": {"key": "v2"},
	}

	srcScoped := vault.NewScopeClient(vault.NewMockClient(srcData), "secret/prod")
	dstScoped := vault.NewScopeClient(vault.NewMockClient(dstData), "secret/prod")

	r := diff.NewRunner(diff.RunnerConfig{
		Src:        srcScoped,
		Dst:        dstScoped,
		SrcPath:    "secret/prod/app",
		DstPath:    "secret/prod/app",
		RedactMode: redact.ModeNone,
		ShowAll:    true,
	})

	results, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(results))
	}
	if results[0].Status != diff.StatusModified {
		t.Errorf("expected Modified, got %v", results[0].Status)
	}
}
