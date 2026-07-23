package providers

import (
	"fmt"
	"sort"
	"sync"
)

// Registry of providers. Providers register themselves in their init() function.
var (
	mu       sync.RWMutex
	registry = map[string]Provider{}
)

// Register registers a provider for a given providerType. 
// It panics if a provider is already registered for the given providerType.
func Register(providerType string, p Provider) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[providerType]; exists {
		panic(fmt.Sprintf("providers: Register called twice for providerType %q", providerType))
	}
	registry[providerType] = p
}

// Get returns the provider registered for the given providerType.
func Get(providerType string) (Provider, error) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := registry[providerType]
	if !ok {
		return nil, fmt.Errorf("providers: no provider registered for providerType %q", providerType)
	}
	return p, nil
}

// Registered returns a sorted list of all registered provider types.
func Registered() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// FetchFromStore fetches secrets from the backend specified in the SecretStore.
// It creates a new client for the provider, fetches the secrets, and then closes the client.
func FetchFromStore(ctx context.Context, store *secretsv1alpha1.SecretStore, refs []secretsv1alpha1.RemoteRef) (map[string][]byte, error) {
	provider, err := Get(store.Spec.ProviderType)
	if err != nil {
		return nil, err
	}

	client, err := provider.NewClient(ctx, store)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.GetSecrets(ctx, refs)
}