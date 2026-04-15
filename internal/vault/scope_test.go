package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func TestScopeClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner client")
		}
	}()
	vault.NewScopeClient(nil, "secret/prod")
}

func TestScopeClient_EmptyScopeReturnsInner(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/prod": {"key": "val"},
	})
	client := vault.NewScopeClient(mock, "")
	if client == mock {
		// acceptable — empty scope returns inner directly
	}
	secrets, err := client.ReadSecrets(context.Background(), "secret/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "val" {
		t.Errorf("expected val, got %q", secrets["key"])
	}
}

func TestScopeClient_AllowsMatchingPath(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/prod/db": {"pass": "s3cr3t"},
	})
	client := vault.NewScopeClient(mock, "secret/prod")
	secrets, err := client.ReadSecrets(context.Background(), "secret/prod/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["pass"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", secrets["pass"])
	}
}

func TestScopeClient_AllowsExactScopePath(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/prod": {"token": "abc"},
	})
	client := vault.NewScopeClient(mock, "secret/prod")
	secrets, err := client.ReadSecrets(context.Background(), "secret/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["token"] != "abc" {
		t.Errorf("expected abc, got %q", secrets["token"])
	}
}

func TestScopeClient_RejectsOutOfScopePath(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{})
	client := vault.NewScopeClient(mock, "secret/prod")
	_, err := client.ReadSecrets(context.Background(), "secret/staging/db")
	if err == nil {
		t.Fatal("expected error for out-of-scope path")
	}
	var scopeErr *vault.ErrOutOfScope
	if !errors.As(err, &scopeErr) {
		t.Errorf("expected ErrOutOfScope, got %T: %v", err, err)
	}
}

func TestScopeClient_StripsScopeLeadingSlash(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/prod/api": {"key": "v"},
	})
	client := vault.NewScopeClient(mock, "/secret/prod")
	_, err := client.ReadSecrets(context.Background(), "secret/prod/api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScopeClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil)
	var _ vault.SecretReader = vault.NewScopeClient(mock, "secret")
}
