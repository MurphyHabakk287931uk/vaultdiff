package vault

import (
	"fmt"
	"os"
	"strings"

	hashivault "github.com/hashicorp/vault/api"
)

// TokenAuthFromEnv reads a Vault token from the environment.
// It checks VAULT_TOKEN first, then falls back to ~/.vault-token via the SDK.
func TokenAuthFromEnv() (string, error) {
	token := strings.TrimSpace(os.Getenv("VAULT_TOKEN"))
	if token != "" {
		return token, nil
	}
	return "", fmt.Errorf("VAULT_TOKEN is not set")
}

// LoginWithToken configures the provided Vault client to use a static token.
func LoginWithToken(client *hashivault.Client, token string) error {
	if token == "" {
		return fmt.Errorf("token must not be empty")
	}
	client.SetToken(token)
	return nil
}

// RenewToken attempts a token self-renewal and returns the new TTL in seconds.
// Returns an error if the token is not renewable or the renewal fails.
func RenewToken(client *hashivault.Client) (int, error) {
	secret, err := client.Auth().Token().RenewSelf(0)
	if err != nil {
		return 0, fmt.Errorf("token renewal failed: %w", err)
	}
	if secret == nil || secret.Auth == nil {
		return 0, fmt.Errorf("renewal response contained no auth data")
	}
	return secret.Auth.LeaseDuration, nil
}
