package vault

import (
	"errors"
	"fmt"
	"os"
)

// AuthMethod represents a Vault authentication method.
type AuthMethod string

const (
	AuthToken      AuthMethod = "token"
	AuthAppRole    AuthMethod = "approle"
	AuthKubernetes AuthMethod = "kubernetes"
)

// AuthConfig holds credentials for a given auth method.
type AuthConfig struct {
	Method   AuthMethod
	Token    string
	RoleID   string
	SecretID string
	JWTPath  string
	Role     string
}

// ParseAuthMethod parses a string into an AuthMethod.
func ParseAuthMethod(s string) (AuthMethod, error) {
	switch AuthMethod(s) {
	case AuthToken:
		return AuthToken, nil
	case AuthAppRole:
		return AuthAppRole, nil
	case AuthKubernetes:
		return AuthKubernetes, nil
	default:
		return "", fmt.Errorf("unknown auth method %q: must be one of token, approle, kubernetes", s)
	}
}

// AuthConfigFromEnv builds an AuthConfig by reading well-known environment
// variables for the given method.
func AuthConfigFromEnv(method AuthMethod) (AuthConfig, error) {
	cfg := AuthConfig{Method: method}
	switch method {
	case AuthToken:
		cfg.Token = os.Getenv("VAULT_TOKEN")
		if cfg.Token == "" {
			return cfg, errors.New("VAULT_TOKEN must be set for token auth")
		}
	case AuthAppRole:
		cfg.RoleID = os.Getenv("VAULT_ROLE_ID")
		cfg.SecretID = os.Getenv("VAULT_SECRET_ID")
		if cfg.RoleID == "" || cfg.SecretID == "" {
			return cfg, errors.New("VAULT_ROLE_ID and VAULT_SECRET_ID must be set for approle auth")
		}
	case AuthKubernetes:
		cfg.JWTPath = os.Getenv("VAULT_K8S_JWT_PATH")
		if cfg.JWTPath == "" {
			cfg.JWTPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		}
		cfg.Role = os.Getenv("VAULT_K8S_ROLE")
		if cfg.Role == "" {
			return cfg, errors.New("VAULT_K8S_ROLE must be set for kubernetes auth")
		}
	default:
		return cfg, fmt.Errorf("unsupported auth method: %s", method)
	}
	return cfg, nil
}
