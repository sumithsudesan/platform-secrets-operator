package providers

import (
	"context"

	secretsv1alpha1 "github.com/sumithsudesan/platform-secrets-operator/api/v1alpha1"
)

// Aws secret provider implementation. 
// This is a stub implementation and does not actually connect to an AWS backend.
const ProviderTypeAws = "aws"

func init() {
	Register(ProviderTypeAws, &awsProvider{})
}

// awsProvider is a stub implementation of the Provider interface for AWS.
type awsProvider struct{}

// NewClient returns a new awsClient.
func (p *awsProvider) NewClient(ctx context.Context, store *secretsv1alpha1.SecretStore) (Client, error) {
	return nil, nil
}

// awsClient is a stub implementation of the Client interface for AWS.
type awsClient struct{}

// GetSecrets is a stub implementation that returns nil for both the map and error.
func (c *awsClient) GetSecrets(ctx context.Context, refs []secretsv1alpha1.RemoteRef) (map[string][]byte, error) {
	return nil, nil
}

// Close is a stub implementation that does nothing and returns nil.
func (c *awsClient) Close() error {
	return nil
}