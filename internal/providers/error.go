package providers

import (
	"fmt"

	secretsv1alpha1 "github.com/sumithsudesan/platform-secrets-operator/api/v1alpha1"
)

// AuthError is returned when a provider rejects the credentials 
// provided in a SecretStore.
type AuthError struct {
	Store string // name of the SecretStore whose credentials were rejected
	Err   error
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("provider auth error for store %q: %v", e.Store, e.Err)
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// NotFoundError is returned when a secret is not found in the backend
//  for a given RemoteRef.
type NotFoundError struct {
	Ref secretsv1alpha1.RemoteRef
	Err error
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("secret not found in backend for ref %+v: %v", e.Ref, e.Err)
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}