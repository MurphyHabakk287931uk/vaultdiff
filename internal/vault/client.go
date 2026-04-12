package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// Config holds configuration for creating a Vault client.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client from the given config.
// If Address or Token are empty, it falls back to environment variables
// VAULT_ADDR and VAULT_TOKEN respectively.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr != "" {
		vcfg.Address = addr
	}

	vc, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to create client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault: no token provided (set VAULT_TOKEN or pass --token)")
	}
	vc.SetToken(token)

	return &Client{vc: vc}, nil
}

// ReadSecrets reads the key/value pairs at the given KV v2 path.
// path should be in the form "secret/data/myapp/config".
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.vc.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: path %q not found or empty", path)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 or non-data path — return raw data
		return flattenData(secret.Data)
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: unexpected data format at %q", path)
	}
	return flattenData(dataMap)
}

func flattenData(raw map[string]interface{}) (map[string]string, error) {
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			out[k] = val
		case nil:
			out[k] = ""
		default:
			out[k] = fmt.Sprintf("%v", val)
		}
	}
	return out, nil
}
