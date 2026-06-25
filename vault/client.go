package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with environment context.
type Client struct {
	api  *vaultapi.Client
	Env  string
}

// NewClient creates a new Vault client configured from environment variables
// or explicit address/token overrides.
func NewClient(env, address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()

	if address != "" {
		cfg.Address = address
	} else if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Address = addr
	} else {
		return nil, fmt.Errorf("vault address not set: use --address or VAULT_ADDR")
	}

	if err := cfg.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("reading vault environment: %w", err)
	}

	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	} else if t := os.Getenv("VAULT_TOKEN"); t != "" {
		client.SetToken(t)
	} else {
		return nil, fmt.Errorf("vault token not set: use --token or VAULT_TOKEN")
	}

	return &Client{api: client, Env: env}, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// Pass version=0 to read the latest version.
func (c *Client) ReadSecretVersion(mountPath, secretPath string, version int) (map[string]interface{}, int, error) {
	kvClient := c.api.KVv2(mountPath)

	var (
		secret *vaultapi.KVSecret
		err    error
	)

	if version == 0 {
		secret, err = kvClient.Get(nil, secretPath)
	} else {
		secret, err = kvClient.GetVersion(nil, secretPath, version)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("reading secret %s/%s: %w", mountPath, secretPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, 0, fmt.Errorf("secret %s/%s not found", mountPath, secretPath)
	}

	ver := 0
	if secret.VersionMetadata != nil {
		ver = secret.VersionMetadata.Version
	}

	return secret.Data, ver, nil
}
