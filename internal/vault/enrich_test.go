package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestEnrichClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	vault.NewEnrichClient(nil, map[string]string{})
}

func TestEnrichClient_PanicsOnNilExtra(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil extra")
		}
	}()
	mock := vault.NewMockClient(map[string]map[string]string{})
	vault.NewEnrichClient(mock, nil)
}

func TestEnrichClient_InjectsExtra(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"db_pass": "s3cr3t"},
	})
	client := vault.NewEnrichClient(mock, map[string]string{"env": "prod"})
	secrets, err := client.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", secrets["env"])
	}
	if secrets["db_pass"] != "s3cr3t" {
		t.Errorf("expected db_pass=s3cr3t, got %q", secrets["db_pass"])
	}
}

func TestEnrichClient_ExtraOverridesInner(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"": "staging"},
	})
	client := vault.NewEnrichClient(mock, map[string]string{"env": "prod"})
	secrets, err := client.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["env"] != "prod" {
		t.Errorf("extra should override inner: got %q", secrets["env"])
	}
}

func TestEnrichClient_PropagatesError(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{})
	mock.SetError("secret/missing", errors.New("not found"))
	client := vault.NewEnrichClient(mock, map[string]string{"env": "prod"})
	_, err := client.ReadSecrets(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestEnrichClient_DoesNotMutateExtraMap(t *testing.T) {
	extra := map[string]string{"env": "prod"}
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {},
	})
	client := vault.NewEnrichClient(mock, extra)
	_, _ = client.ReadSecrets(context.Background(), "secret/app")
	extra["injected"] = "yes"
	secrets, _ := client.ReadSecrets(context.Background(), "secret/app")
	if _, ok := secrets["injected"]; ok {
		t.Error("mutation of original extra map affected client")
	}
}

func TestEnrichClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{})
	var _ vault.SecretReader = vault.NewEnrichClient(mock, map[string]string{})
}
