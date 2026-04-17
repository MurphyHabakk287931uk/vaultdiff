package diff_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/redact"
	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestRunner_SplitClientExposesNamespacedKeys(t *testing.T) {
	left := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"TOKEN": "abc", "HOST": "prod.example.com"},
	}, nil)
	right := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"TOKEN": "xyz", "HOST": "prod.example.com"},
	}, nil)
	split := vault.NewSplitClient(left, "prod", right, "staging")

	// Use the split client as both src and dst with the same path so we can
	// observe the namespaced keys flowing through the diff runner.
	static := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {},
	}, nil)

	r := diff.NewRunner(split, static, redact.ModeNone)
	result, err := r.Run("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var foundProd, foundStaging bool
	for _, e := range result {
		if strings.HasPrefix(e.Key, "prod/") {
			foundProd = true
		}
		if strings.HasPrefix(e.Key, "staging/") {
			foundStaging = true
		}
	}
	if !foundProd {
		t.Error("expected keys prefixed with 'prod/'")
	}
	if !foundStaging {
		t.Error("expected keys prefixed with 'staging/'")
	}
}
