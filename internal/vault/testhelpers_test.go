package vault

import (
	hashivault "github.com/hashicorp/vault/api"
)

// newTestVaultClient creates a minimal Vault API client pointed at a
// dummy address, suitable for unit tests that do not make real HTTP calls.
func newTestVaultClient() (*hashivault.Client, error) {
	cfg := hashivault.DefaultConfig()
	cfg.Address = "http://127.0.0.1:8200"
	return hashivault.NewClient(cfg)
}
