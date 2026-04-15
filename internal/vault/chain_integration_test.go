package vault_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

// TestChainClient_WithCacheAndFallback demonstrates a realistic layering:
// a cache wrapping a primary, chained with a secondary.
func TestChainClient_WithCacheAndFallback(t *testing.T) {
	primary := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "hunter2"},
	})
	secondary := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "fallback"},
	})

	cached := vault.NewCachedClient(primary)
	chain := vault.NewChainClient(cached, secondary)

	// First read – populates cache via primary.
	got, err := chain.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_PASS"] != "hunter2" {
		t.Fatalf("expected hunter2, got %s", got["DB_PASS"])
	}

	// Second read – served from cache, secondary never consulted.
	got2, err := chain.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got2["DB_PASS"] != "hunter2" {
		t.Fatalf("expected hunter2 from cache, got %s", got2["DB_PASS"])
	}
}

// TestChainClient_ErrorPropagationAcrossWrappers ensures that when every
// client in the chain fails the terminal error surfaces correctly.
func TestChainClient_ErrorPropagationAcrossWrappers(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	bad := vault.NewMockClient(nil)
	bad.SetError("secret/x", sentinel)

	chain := vault.NewChainClient(bad)
	_, err := chain.ReadSecrets("secret/x")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
}
