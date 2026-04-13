package vault

import (
	"testing"
)

func TestParseAuthMethod_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     AuthMethod
	}{
		{"token", AuthToken},
		{"approle", AuthAppRole},
		{"kubernetes", AuthKubernetes},
	}
	for _, tc := range cases {
		got, err := ParseAuthMethod(tc.input)
		if err != nil {
			t.Errorf("ParseAuthMethod(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseAuthMethod(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseAuthMethod_Invalid(t *testing.T) {
	_, err := ParseAuthMethod("ldap")
	if err == nil {
		t.Error("expected error for unknown auth method, got nil")
	}
}

func TestAuthConfigFromEnv_Token(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.abc123")
	cfg, err := AuthConfigFromEnv(AuthToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Token != "s.abc123" {
		t.Errorf("expected token s.abc123, got %q", cfg.Token)
	}
}

func TestAuthConfigFromEnv_Token_Missing(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := AuthConfigFromEnv(AuthToken)
	if err == nil {
		t.Error("expected error when VAULT_TOKEN is empty")
	}
}

func TestAuthConfigFromEnv_AppRole(t *testing.T) {
	t.Setenv("VAULT_ROLE_ID", "role-id")
	t.Setenv("VAULT_SECRET_ID", "secret-id")
	cfg, err := AuthConfigFromEnv(AuthAppRole)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RoleID != "role-id" || cfg.SecretID != "secret-id" {
		t.Errorf("unexpected approle config: %+v", cfg)
	}
}

func TestAuthConfigFromEnv_AppRole_Missing(t *testing.T) {
	t.Setenv("VAULT_ROLE_ID", "")
	t.Setenv("VAULT_SECRET_ID", "")
	_, err := AuthConfigFromEnv(AuthAppRole)
	if err == nil {
		t.Error("expected error when approle env vars are empty")
	}
}

func TestAuthConfigFromEnv_Kubernetes(t *testing.T) {
	t.Setenv("VAULT_K8S_ROLE", "my-role")
	t.Setenv("VAULT_K8S_JWT_PATH", "/tmp/jwt")
	cfg, err := AuthConfigFromEnv(AuthKubernetes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Role != "my-role" || cfg.JWTPath != "/tmp/jwt" {
		t.Errorf("unexpected k8s config: %+v", cfg)
	}
}

func TestAuthConfigFromEnv_Kubernetes_DefaultJWT(t *testing.T) {
	t.Setenv("VAULT_K8S_ROLE", "my-role")
	t.Setenv("VAULT_K8S_JWT_PATH", "")
	cfg, err := AuthConfigFromEnv(AuthKubernetes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	if cfg.JWTPath != want {
		t.Errorf("expected default JWT path %q, got %q", want, cfg.JWTPath)
	}
}

func TestAuthConfigFromEnv_Kubernetes_MissingRole(t *testing.T) {
	t.Setenv("VAULT_K8S_ROLE", "")
	_, err := AuthConfigFromEnv(AuthKubernetes)
	if err == nil {
		t.Error("expected error when VAULT_K8S_ROLE is empty")
	}
}
