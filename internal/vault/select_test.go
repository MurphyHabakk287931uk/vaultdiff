package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func TestSelectClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner client")
		}
	}()
	vault.NewSelectClient(nil, "KEY")
}

func TestSelectClient_NoKeys_ReturnsInnerDirectly(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/data": {"A": "1", "B": "2"},
	})
	client := vault.NewSelectClient(mock)
	// Should be the mock itself (no wrapper)
	secrets, err := client.ReadSecrets(context.Background(), "secret/data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(secrets))
	}
}

func TestSelectClient_FiltersToAllowedKeys(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/cfg": {"DB": "postgres", "API": "key123", "SECRET": "hidden"},
	})
	client := vault.NewSelectClient(mock, "DB", "API")
	secrets, err := client.ReadSecrets(context.Background(), "secret/cfg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(secrets), secrets)
	}
	if _, ok := secrets["SECRET"]; ok {
		t.Error("SECRET should have been filtered out")
	}
	if secrets["DB"] != "postgres" {
		t.Errorf("DB: got %q, want %q", secrets["DB"], "postgres")
	}
}

func TestSelectClient_KeyNotPresent_DroppedSilently(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/cfg": {"A": "1"},
	})
	client := vault.NewSelectClient(mock, "A", "MISSING")
	secrets, err := client.ReadSecrets(context.Background(), "secret/cfg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 1 {
		t.Fatalf("expected 1 key, got %d", len(secrets))
	}
}

func TestSelectClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := vault.NewMockClient(nil)
	mock.SetError("secret/cfg", sentinel)
	client := vault.NewSelectClient(mock, "A")
	_, err := client.ReadSecrets(context.Background(), "secret/cfg")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got: %v", err)
	}
}

func TestSelectClient_EmptyStringKeysIgnored(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/cfg": {"A": "1", "B": "2"},
	})
	// Only empty keys supplied — should return inner directly
	client := vault.NewSelectClient(mock, "", "")
	secrets, err := client.ReadSecrets(context.Background(), "secret/cfg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 2 {
		t.Fatalf("expected 2 keys (passthrough), got %d", len(secrets))
	}
}

func TestSelectClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil)
	var _ vault.SecretReader = vault.NewSelectClient(mock, "K")
}
