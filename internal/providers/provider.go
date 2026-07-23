package providers

import (
	"context"

	secretsv1alpha1 "github.com/sumithsudesan/platform-secrets-operator/api/v1alpha1"
)

// Interface for a provider that can create a Client for a specific backend.
// vault, AWS, and mock providers implement this interface.
type Provider interface {
	NewClient(ctx context.Context, store *secretsv1alpha1.SecretStore) (Client, error)
}

// Client is a backend-specific client that can resolve RemoteRefs to their
// decoded values. Each provider (vault, AWS, mock) implements this interface.
type Client interface {
	// GetSecrets resolves every entry in refs against the backend
	GetSecrets(ctx context.Context, refs []secretsv1alpha1.RemoteRef) (map[string][]byte, error)

	// Close disconnects and releases resources held by the client.
	Close() error
}

// Returns the local key for a RemoteRef. 
// The local key is used to store the resolved secret in the SecretStore's status. 
// The order of precedence is:
// 1. LocalKey (if set)
// 2. Key (if set)
// 3. SecretID (if set)
func LocalKeyFor(ref secretsv1alpha1.RemoteRef) string {
	switch {
	case ref.LocalKey != "":
		return ref.LocalKey
	case ref.Key != "":
		return ref.Key
	default:
		return ref.SecretID
	}
}