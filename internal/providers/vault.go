package providers

import (
	"context"

	secretsv1alpha1 "github.com/sumithsudesan/platform-secrets-operator/api/v1alpha1"
)

// Vault provider implementation. 
// This is a stub implementation and does not actually connect to a Vault backend.
const ProviderTypeVault = "vault"

func init() {
	Register(ProviderTypeVault, &vaultProvider{})
}

// vaultProvider is a stub implementation of the Provider interface for Vault.
type vaultProvider struct{}

// NewClient returns a new vaultClient.
func (p *vaultProvider) NewClient(ctx context.Context, store *secretsv1alpha1.SecretStore) (Client, error) {
	return nil, nil
}

// vaultClient is a stub implementation of the Client interface for Vault.
type vaultClient struct{}

// GetSecrets is a stub implementation that returns nil for both the map and error.
func (c *vaultClient) GetSecrets(ctx context.Context, refs []secretsv1alpha1.RemoteRef) (map[string][]byte, error) {
	return nil, nil
}

// Close is a stub implementation that does nothing and returns nil.
func (c *vaultClient) Close() error {
	return nil
}