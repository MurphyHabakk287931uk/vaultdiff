package vault_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestNewNamespaceClient_EmptyNamespace(t *testing.T) {
	mock := vault.NewMockClient(nil, nil)
	_, err := vault.NewNamespaceClient(mock, "")
	if err == nil {
		t.Fatal("expected error for empty namespace, got nil")
	}
}

func TestNewNamespaceClient_WhitespaceSlashes(t *testing.T) {
	mock := vault.NewMockClient(nil, nil)
	_, err := vault.NewNamespaceClient(mock, "/")
	if err == nil {
		t.Fatal("expected error for slash-only namespace, got nil")
	}
}

func TestNamespaceClient_PrependsSingleSegment(t *testing.T) {
	captured := ""
	mock := vault.NewMockClient(map[string]map[string]string{
		"acme/secret/db": {"password": "s3cr3t"},
	}, nil)

	nc, err := vault.NewNamespaceClient(mock, "acme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = captured

	secrets, err := nc.ReadSecrets("secret/db")
	if err != nil {
		t.Fatalf("ReadSecrets error: %v", err)
	}
	if secrets["password"] != "s3cr3t" {
		t.Errorf("expected 's3cr3t', got %q", secrets["password"])
	}
}

func TestNamespaceClient_StripsLeadingSlashFromPath(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"org/secret/api": {"key": "abc123"},
	}, nil)

	nc, _ := vault.NewNamespaceClient(mock, "org")
	secrets, err := nc.ReadSecrets("/secret/api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", secrets["key"])
	}
}

func TestNamespaceClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := vault.NewMockClient(nil, sentinel)

	nc, _ := vault.NewNamespaceClient(mock, "ns")
	_, err := nc.ReadSecrets("any/path")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestNamespaceClient_Namespace(t *testing.T) {
	mock := vault.NewMockClient(nil, nil)
	nc, _ := vault.NewNamespaceClient(mock, "/prod/")
	if nc.Namespace() != "prod" {
		t.Errorf("expected 'prod', got %q", nc.Namespace())
	}
}

func TestNamespaceClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil, nil)
	nc, _ := vault.NewNamespaceClient(mock, "ns")
	var _ vault.SecretReader = nc
}
