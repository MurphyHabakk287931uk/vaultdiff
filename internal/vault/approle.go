package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// AppRoleAuth holds credentials for AppRole authentication.
type AppRoleAuth struct {
	RoleID   string
	SecretID string
}

// AppRoleAuthFromEnv reads AppRole credentials from environment variables.
// It expects VAULT_ROLE_ID and VAULT_SECRET_ID to be set.
func AppRoleAuthFromEnv() (AppRoleAuth, error) {
	roleID := os.Getenv("VAULT_ROLE_ID")
	if roleID == "" {
		return AppRoleAuth{}, fmt.Errorf("VAULT_ROLE_ID is not set")
	}
	secretID := os.Getenv("VAULT_SECRET_ID")
	if secretID == "" {
		return AppRoleAuth{}, fmt.Errorf("VAULT_SECRET_ID is not set")
	}
	return AppRoleAuth{RoleID: roleID, SecretID: secretID}, nil
}

// LoginWithAppRole authenticates the Vault client using AppRole and returns
// a new client configured with the resulting token.
func LoginWithAppRole(cfg *vaultapi.Config, auth AppRoleAuth, mountPath string) (*vaultapi.Client, error) {
	if mountPath == "" {
		mountPath = "approle"
	}

	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	data := map[string]interface{}{
		"role_id":   auth.RoleID,
		"secret_id": auth.SecretID,
	}

	path := fmt.Sprintf("auth/%s/login", mountPath)
	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return nil, fmt.Errorf("approle login at %q: %w", path, err)
	}
	if secret == nil || secret.Auth == nil {
		return nil, fmt.Errorf("approle login returned no auth info")
	}

	client.SetToken(secret.Auth.ClientToken)
	return client, nil
}
