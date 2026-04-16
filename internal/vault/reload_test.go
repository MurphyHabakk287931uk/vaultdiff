package vault_test

import (
	"context"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func TestReloadClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	vault.NewReloadClient(nil)
}

func TestReloadClient_DelegatesRead(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	c := vault.NewReloadClient(mock)
	secrets, err := c.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "value" {
		t.Errorf("expected value, got %q", secrets["key"])
	}
}

func TestReloadClient_ReloadsInner(t *testing.T) {
	old := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "old"},
	})
	new_ := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "new"},
	})
	c := vault.NewReloadClient(old)
	if err := c.Reload(new_); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets, err := c.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "new" {
		t.Errorf("expected new, got %q", secrets["key"])
	}
}

func TestReloadClient_ReloadNilReturnsError(t *testing.T) {
	mock := vault.NewMockClient(nil)
	c := vault.NewReloadClient(mock)
	if err := c.Reload(nil); err == nil {
		t.Fatal("expected error when reloading with nil")
	}
}

func TestReloadClient_Current(t *testing.T) {
	mock := vault.NewMockClient(nil)
	c := vault.NewReloadClient(mock)
	if c.Current() != mock {
		t.Error("Current() should return the inner client")
	}
}

func TestReloadClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil)
	var _ vault.SecretReader = vault.NewReloadClient(mock)
}
