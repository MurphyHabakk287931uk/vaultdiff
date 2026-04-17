package vault_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func TestAliasClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	vault.NewAliasClient(nil, map[string]string{})
}

func TestAliasClient_PanicsOnNilAliases(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil aliases")
		}
	}()
	mock := vault.NewMockClient(nil, nil)
	vault.NewAliasClient(mock, nil)
}

func TestAliasClient_NoAlias_ForwardsOriginalPath(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/real": {"key": "val"},
	}, nil)
	client := vault.NewAliasClient(mock, map[string]string{})

	got, err := client.ReadSecrets("secret/real")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "val" {
		t.Errorf("expected val, got %q", got["key"])
	}
}

func TestAliasClient_ResolvesAlias(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/prod/db": {"password": "s3cr3t"},
	}, nil)
	client := vault.NewAliasClient(mock, map[string]string{
		"db": "secret/prod/db",
	})

	got, err := client.ReadSecrets("db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["password"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", got["password"])
	}
}

func TestAliasClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault down")
	mock := vault.NewMockClient(nil, map[string]error{
		"secret/prod/db": sentinel,
	})
	client := vault.NewAliasClient(mock, map[string]string{
		"db": "secret/prod/db",
	})

	_, err := client.ReadSecrets("db")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestAliasClient_AliasMethod(t *testing.T) {
	mock := vault.NewMockClient(nil, nil)
	client := vault.NewAliasClient(mock, map[string]string{
		"app": "secret/app/config",
	})

	if got := client.Alias("app"); got != "secret/app/config" {
		t.Errorf("expected secret/app/config, got %q", got)
	}
	if got := client.Alias("other"); got != "other" {
		t.Errorf("expected other, got %q", got)
	}
}

func TestAliasClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil, nil)
	var _ vault.SecretReader = vault.NewAliasClient(mock, map[string]string{})
}
