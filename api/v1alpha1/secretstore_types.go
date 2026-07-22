package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VaultConfig holds connection details for a HashiCorp Vault backend.
type VaultConfig struct {
	Address string `json:"address"`
	Role    string `json:"role"`
}

// AWSConfig holds connection details for an AWS Secrets Manager backend.
type AWSConfig struct {
	Region  string `json:"region"`
	RoleArn string `json:"roleArn"`
}

// MockConfig holds settings for the local mock backend used for testing.
type MockConfig struct {
	Endpoint string `json:"endpoint,omitempty"`
}

// SecretStoreSpec defines how to connect and authenticate to a backend.
type SecretStoreSpec struct {
	// +kubebuilder:validation:Enum=vault;aws-secrets-manager;mock
	ProviderType string `json:"providerType"`

	Vault *VaultConfig `json:"vault,omitempty"`
	AWS   *AWSConfig   `json:"aws,omitempty"`
	Mock  *MockConfig  `json:"mock,omitempty"`
}

// SecretStoreStatus reports the observed state of a SecretStore.
type SecretStoreStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced

// SecretStore describes where a backend lives and how to authenticate to it.
type SecretStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretStoreSpec   `json:"spec,omitempty"`
	Status SecretStoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecretStoreList contains a list of SecretStore.
type SecretStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretStore{}, &SecretStoreList{})
}