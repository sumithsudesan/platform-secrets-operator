package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RemoteRef points to one value in the backend, and where it should land locally.
type RemoteRef struct {
	// Path is the Vault secret path (used when the SecretStore's providerType is "vault").
	Path string `json:"path,omitempty"`

	// Key is the field name to read from the remote secret.
	Key string `json:"key,omitempty"`

	// SecretID is the AWS secret name (used when providerType is "aws-secrets-manager").
	SecretID string `json:"secretId,omitempty"`

	// VersionStage is the AWS version stage to read, e.g. AWSCURRENT.
	VersionStage string `json:"versionStage,omitempty"`

	// LocalKey is the key name to use in the generated Kubernetes Secret.
	// If empty, it defaults to Key (or SecretID).
	LocalKey string `json:"localKey,omitempty"`
}

// ManagedSecretSpec describes what to fetch and how to deliver it.
type ManagedSecretSpec struct {
	// StoreRef is the name of a SecretStore in the same namespace.
	StoreRef string `json:"storeRef"`

	// TargetSecretName is the Kubernetes Secret this ManagedSecret will create/update.
	TargetSecretName string `json:"targetSecretName"`

	// RefreshInterval controls how often the operator re-checks the backend.
	RefreshInterval metav1.Duration `json:"refreshInterval"`

	// RemoteRefs lists the remote values to fetch.
	RemoteRefs []RemoteRef `json:"remoteRefs"`

	// +kubebuilder:validation:Enum=Retain;Delete
	// +kubebuilder:default=Retain
	DeletionPolicy string `json:"deletionPolicy,omitempty"`
}

// ManagedSecretStatus reports what the operator has actually done.
type ManagedSecretStatus struct {
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
	LastSyncTime  *metav1.Time       `json:"lastSyncTime,omitempty"`
	NotFoundCount int32              `json:"notFoundCount,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced

// ManagedSecret describes a desired Kubernetes Secret,sourced from a SecretStore.
type ManagedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedSecretSpec   `json:"spec,omitempty"`
	Status ManagedSecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ManagedSecretList contains a list of ManagedSecret.
type ManagedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedSecret{}, &ManagedSecretList{})
}

