package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestFallbackClient_PrimarySucceeds(t *testing.T) {
	primary := vault.NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "primary-value"},
	})
	secondary := vault.NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "secondary-value"},
	})

	client := vault.NewFallbackClient(primary, secondary, nil)
	got, err := client.ReadSecrets(context.Background(), "secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "primary-value" {
		t.Errorf("expected primary-value, got %q", got["key"])
	}
}

func TestFallbackClient_FallsBackOnAnyError(t *testing.T) {
	primary := vault.NewMockClient(nil)
	primary.SetError("secret/missing", vault.ErrNotFound)
	secondary := vault.NewMockClient(map[string]map[string]string{
		"secret/missing": {"key": "fallback-value"},
	})

	client := vault.NewFallbackClient(primary, secondary, nil)
	got, err := client.ReadSecrets(context.Background(), "secret/missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "fallback-value" {
		t.Errorf("expected fallback-value, got %q", got["key"])
	}
}

func TestFallbackClient_IsNotFound_Predicate(t *testing.T) {
	permanent := errors.New("permission denied")
	primary := vault.NewMockClient(nil)
	primary.SetError("secret/x", permanent)
	secondary := vault.NewMockClient(map[string]map[string]string{
		"secret/x": {"key": "should-not-reach"},
	})

	client := vault.NewFallbackClient(primary, secondary, vault.IsNotFound)
	_, err := client.ReadSecrets(context.Background(), "secret/x")
	if !errors.Is(err, permanent) {
		t.Errorf("expected permanent error to propagate, got %v", err)
	}
}

func TestFallbackClient_IsNotFound_TriggersOnNotFound(t *testing.T) {
	primary := vault.NewMockClient(nil)
	primary.SetError("secret/y", vault.ErrNotFound)
	secondary := vault.NewMockClient(map[string]map[string]string{
		"secret/y": {"env": "staging"},
	})

	client := vault.NewFallbackClient(primary, secondary, vault.IsNotFound)
	got, err := client.ReadSecrets(context.Background(), "secret/y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["env"] != "staging" {
		t.Errorf("expected staging, got %q", got["env"])
	}
}

func TestFallbackClient_PanicsOnNilPrimary(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil primary")
		}
	}()
	secondary := vault.NewMockClient(nil)
	vault.NewFallbackClient(nil, secondary, nil)
}

func TestFallbackClient_PanicsOnNilSecondary(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil secondary")
		}
	}()
	primary := vault.NewMockClient(nil)
	vault.NewFallbackClient(primary, nil, nil)
}

func TestFallbackClient_ImplementsSecretReader(t *testing.T) {
	primary := vault.NewMockClient(nil)
	secondary := vault.NewMockClient(nil)
	var _ vault.SecretReader = vault.NewFallbackClient(primary, secondary, nil)
}
