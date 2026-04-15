package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func TestQuotaClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	vault.NewQuotaClient(nil, vault.DefaultQuotaConfig())
}

func TestQuotaClient_ZeroQuota_Unlimited(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/a": {"k": "v"},
	})
	client := vault.NewQuotaClient(mock, vault.QuotaConfig{MaxReads: 0})
	ctx := context.Background()
	for i := 0; i < 50; i++ {
		if _, err := client.ReadSecrets(ctx, "secret/a"); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i+1, err)
		}
	}
}

func TestQuotaClient_AllowsReadsUpToLimit(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/a": {"k": "v"},
	})
	client := vault.NewQuotaClient(mock, vault.QuotaConfig{MaxReads: 3})
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if _, err := client.ReadSecrets(ctx, "secret/a"); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i+1, err)
		}
	}
}

func TestQuotaClient_ExceedsLimit(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/a": {"k": "v"},
	})
	client := vault.NewQuotaClient(mock, vault.QuotaConfig{MaxReads: 2})
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		client.ReadSecrets(ctx, "secret/a") //nolint:errcheck
	}
	_, err := client.ReadSecrets(ctx, "secret/a")
	if !errors.Is(err, vault.ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestQuotaClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil)
	var _ vault.SecretReader = vault.NewQuotaClient(mock, vault.DefaultQuotaConfig())
}
