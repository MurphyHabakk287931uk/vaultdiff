package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestReadOnlyClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner, got none")
		}
	}()
	vault.NewReadOnlyClient(nil)
}

func TestReadOnlyClient_DelegatesRead(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	client := vault.NewReadOnlyClient(mock)

	got, err := client.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("expected value %q, got %q", "value", got["key"])
	}
}

func TestReadOnlyClient_RejectsWrite(t *testing.T) {
	mock := vault.NewMockClient(nil)
	client := vault.NewReadOnlyClient(mock)

	err := client.WriteSecrets(context.Background(), "secret/app", map[string]string{"k": "v"})
	if !errors.Is(err, vault.ErrReadOnly) {
		t.Errorf("expected ErrReadOnly, got %v", err)
	}
}

func TestReadOnlyClient_PropagatesReadError(t *testing.T) {
	sentinel := errors.New("backend unavailable")
	mock := vault.NewMockClient(nil)
	mock.SetError("secret/broken", sentinel)
	client := vault.NewReadOnlyClient(mock)

	_, err := client.ReadSecrets(context.Background(), "secret/broken")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestReadOnlyClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil)
	var _ vault.SecretReader = vault.NewReadOnlyClient(mock)
}

func TestReadOnlyClient_ImplementsSecretWriter(t *testing.T) {
	mock := vault.NewMockClient(nil)
	var _ vault.SecretWriter = vault.NewReadOnlyClient(mock)
}
