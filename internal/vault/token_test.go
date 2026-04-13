package vault

import (
	"os"
	"testing"
)

func TestTokenAuthFromEnv_Success(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.testtoken123")

	token, err := TokenAuthFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "s.testtoken123" {
		t.Errorf("expected s.testtoken123, got %q", token)
	}
}

func TestTokenAuthFromEnv_Missing(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")

	_, err := TokenAuthFromEnv()
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is unset")
	}
}

func TestTokenAuthFromEnv_Whitespace(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "  s.trimmed  ")

	token, err := TokenAuthFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "s.trimmed" {
		t.Errorf("expected trimmed token, got %q", token)
	}
}

func TestTokenAuthFromEnv_EmptyString(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")

	_, err := TokenAuthFromEnv()
	if err == nil {
		t.Fatal("expected error for empty VAULT_TOKEN")
	}
}

func TestLoginWithToken_Empty(t *testing.T) {
	err := LoginWithToken(nil, "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestLoginWithToken_SetsToken(t *testing.T) {
	client, err := newTestVaultClient()
	if err != nil {
		t.Skipf("could not create vault client: %v", err)
	}

	if err := LoginWithToken(client, "s.mytoken"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := client.Token(); got != "s.mytoken" {
		t.Errorf("expected token s.mytoken, got %q", got)
	}
}
