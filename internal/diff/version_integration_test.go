package diff_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/redact"
	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestRunner_VersionedClientsCompareVersions(t *testing.T) {
	v1Secrets := map[string]string{"db_pass": "old", "api_key": "same"}
	v2Secrets := map[string]string{"db_pass": "new", "api_key": "same"}

	base := vault.NewMockClient(nil, nil)
	base.SetSecrets("secret/myapp/v1", v1Secrets)
	base.SetSecrets("secret/myapp/v2", v2Secrets)

	src := vault.NewVersionedClient(base, "v1")
	dst := vault.NewVersionedClient(base, "v2")

	runner := diff.NewRunner(src, dst, redact.ModeNone, true)
	results, err := runner.Run("secret/myapp", "secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	changed := 0
	for _, r := range results {
		if r.Status == diff.StatusModified {
			changed++
			if r.Key != "db_pass" {
				t.Errorf("unexpected modified key %q", r.Key)
			}
		}
	}
	if changed != 1 {
		t.Errorf("got %d modified keys, want 1", changed)
	}
}

func TestRunner_VersionedClient_AddedKey(t *testing.T) {
	base := vault.NewMockClient(nil, nil)
	base.SetSecrets("secret/svc/v1", map[string]string{"x": "1"})
	base.SetSecrets("secret/svc/v2", map[string]string{"x": "1", "y": "2"})

	src := vault.NewVersionedClient(base, "v1")
	dst := vault.NewVersionedClient(base, "v2")

	runner := diff.NewRunner(src, dst, redact.ModeNone, false)
	results, err := runner.Run("secret/svc", "secret/svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var added int
	for _, r := range results {
		if r.Status == diff.StatusAdded {
			added++
		}
	}
	if added != 1 {
		t.Errorf("got %d added keys, want 1", added)
	}
}
