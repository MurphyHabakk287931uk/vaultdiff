package vault_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func TestNewPrefixClient_NilInner(t *testing.T) {
	_, err := vault.NewPrefixClient(nil, "prod")
	if err == nil {
		t.Fatal("expected error for nil inner client, got nil")
	}
}

func TestNewPrefixClient_EmptyPrefix_ReturnsInner(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "val"},
	})
	client, err := vault.NewPrefixClient(mock, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should be the original mock, not a PrefixClient wrapper.
	if _, ok := client.(*vault.PrefixClient); ok {
		t.Fatal("expected raw inner client for empty prefix, got PrefixClient")
	}
}

func TestPrefixClient_PrependsSingleSegment(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"prod/secret/app": {"db_pass": "s3cr3t"},
	})
	client, err := vault.NewPrefixClient(mock, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("ReadSecrets error: %v", err)
	}
	if got["db_pass"] != "s3cr3t" {
		t.Errorf("expected 's3cr3t', got %q", got["db_pass"])
	}
}

func TestPrefixClient_StripsLeadingSlashFromPath(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"prod/secret/app": {"key": "value"},
	})
	client, err := vault.NewPrefixClient(mock, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := client.ReadSecrets("/secret/app")
	if err != nil {
		t.Fatalf("ReadSecrets error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("expected 'value', got %q", got["key"])
	}
}

func TestPrefixClient_NormalisesSlashyPrefix(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"prod/secret/app": {"token": "abc"},
	})
	client, err := vault.NewPrefixClient(mock, "/prod/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pc := client.(*vault.PrefixClient)
	if pc.Prefix() != "prod" {
		t.Errorf("expected normalised prefix 'prod', got %q", pc.Prefix())
	}
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("ReadSecrets error: %v", err)
	}
	if got["token"] != "abc" {
		t.Errorf("expected 'abc', got %q", got["token"])
	}
}

func TestPrefixClient_PropagatesError(t *testing.T) {
	expected := errors.New("vault unavailable")
	mock := vault.NewMockClient(nil)
	mock.SetError("prod/secret/broken", expected)
	client, err := vault.NewPrefixClient(mock, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, got := client.ReadSecrets("secret/broken")
	if !errors.Is(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestPrefixClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil)
	client, _ := vault.NewPrefixClient(mock, "ns")
	var _ vault.SecretReader = client
}
