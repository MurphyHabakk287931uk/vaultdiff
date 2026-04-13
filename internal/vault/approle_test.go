package vault

import (
	"testing"
)

func TestAppRoleAuthFromEnv_Success(t *testing.T) {
	t.Setenv("VAULT_ROLE_ID", "my-role-id")
	t.Setenv("VAULT_SECRET_ID", "my-secret-id")

	auth, err := AppRoleAuthFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.RoleID != "my-role-id" {
		t.Errorf("expected RoleID %q, got %q", "my-role-id", auth.RoleID)
	}
	if auth.SecretID != "my-secret-id" {
		t.Errorf("expected SecretID %q, got %q", "my-secret-id", auth.SecretID)
	}
}

func TestAppRoleAuthFromEnv_MissingRoleID(t *testing.T) {
	t.Setenv("VAULT_ROLE_ID", "")
	t.Setenv("VAULT_SECRET_ID", "some-secret")

	_, err := AppRoleAuthFromEnv()
	if err == nil {
		t.Fatal("expected error for missing VAULT_ROLE_ID, got nil")
	}
}

func TestAppRoleAuthFromEnv_MissingSecretID(t *testing.T) {
	t.Setenv("VAULT_ROLE_ID", "some-role")
	t.Setenv("VAULT_SECRET_ID", "")

	_, err := AppRoleAuthFromEnv()
	if err == nil {
		t.Fatal("expected error for missing VAULT_SECRET_ID, got nil")
	}
}

func TestAppRoleAuthFromEnv_BothMissing(t *testing.T) {
	t.Setenv("VAULT_ROLE_ID", "")
	t.Setenv("VAULT_SECRET_ID", "")

	_, err := AppRoleAuthFromEnv()
	if err == nil {
		t.Fatal("expected error when both env vars are missing, got nil")
	}
}

func TestAppRoleAuth_Fields(t *testing.T) {
	auth := AppRoleAuth{
		RoleID:   "role-abc",
		SecretID: "secret-xyz",
	}
	if auth.RoleID != "role-abc" {
		t.Errorf("unexpected RoleID: %q", auth.RoleID)
	}
	if auth.SecretID != "secret-xyz" {
		t.Errorf("unexpected SecretID: %q", auth.SecretID)
	}
}
