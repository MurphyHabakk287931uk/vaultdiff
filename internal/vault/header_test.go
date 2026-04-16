package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nicholasgasior/vaultdiff/internal/vault"
)

func TestHeaderClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	vault.NewHeaderClient(nil, map[string]string{}, "")
}

func TestHeaderClient_PanicsOnNilHeaders(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	mock := vault.NewMockClient(nil, nil)
	vault.NewHeaderClient(mock, nil, "")
}

func TestHeaderClient_InjectsHeaders(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t"},
	}, nil)

	c := vault.NewHeaderClient(mock, map[string]string{"source": "prod", "region": "us-east-1"}, "_meta")

	got, err := c.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %q", got["DB_PASS"])
	}
	if got["_meta.source"] != "prod" {
		t.Errorf("expected _meta.source=prod, got %q", got["_meta.source"])
	}
	if got["_meta.region"] != "us-east-1" {
		t.Errorf("expected _meta.region=us-east-1, got %q", got["_meta.region"])
	}
}

func TestHeaderClient_DefaultPrefix(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"KEY": "val"},
	}, nil)
	c := vault.NewHeaderClient(mock, map[string]string{"env": "staging"}, "")
	got, err := c.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["_header.env"] != "staging" {
		t.Errorf("expected _header.env=staging, got %q", got["_header.env"])
	}
}

func TestHeaderClient_PropagatesError(t *testing.T) {
	expected := errors.New("vault unavailable")
	mock := vault.NewMockClient(nil, expected)
	c := vault.NewHeaderClient(mock, map[string]string{"x": "y"}, "_h")
	_, err := c.ReadSecrets(context.Background(), "secret/app")
	if !errors.Is(err, expected) {
		t.Errorf("expected propagated error, got %v", err)
	}
}

func TestHeaderClient_ImplementsSecretReader(t *testing.T) {
	mock := vault.NewMockClient(nil, nil)
	var _ vault.SecretReader = vault.NewHeaderClient(mock, map[string]string{}, "_h")
}
