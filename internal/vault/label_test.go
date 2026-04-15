package vault_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestLabelClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner, got none")
		}
	}()
	vault.NewLabelClient(nil, "prod")
}

func TestLabelClient_ReturnsLabel(t *testing.T) {
	mc := vault.NewMockClient()
	c := vault.NewLabelClient(mc, "staging")
	if c.Label() != "staging" {
		t.Fatalf("expected label %q, got %q", "staging", c.Label())
	}
}

func TestLabelClient_PassesThroughSecrets(t *testing.T) {
	mc := vault.NewMockClient()
	mc.SetSecrets("secret/app", map[string]string{"key": "value"})

	c := vault.NewLabelClient(mc, "prod")
	got, err := c.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Fatalf("expected value %q, got %q", "value", got["key"])
	}
}

func TestLabelClient_WrapsErrorWithLabel(t *testing.T) {
	mc := vault.NewMockClient()
	mc.SetError("secret/missing", errors.New("not found"))

	c := vault.NewLabelClient(mc, "prod")
	_, err := c.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "[prod]") {
		t.Fatalf("expected label in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected original message in error, got: %v", err)
	}
}

func TestLabelClient_EmptyLabelNoPrefix(t *testing.T) {
	mc := vault.NewMockClient()
	mc.SetError("secret/x", errors.New("boom"))

	c := vault.NewLabelClient(mc, "")
	_, err := c.ReadSecrets("secret/x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if strings.Contains(err.Error(), "[") {
		t.Fatalf("empty label should not add prefix, got: %v", err)
	}
}

func TestLabelClient_ImplementsSecretReader(t *testing.T) {
	mc := vault.NewMockClient()
	var _ vault.SecretReader = vault.NewLabelClient(mc, "test")
}
