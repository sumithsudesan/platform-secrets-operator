// Package v1alpha1 contains API Schema definitions for the secrets v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=secrets.operator.io
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is the group and version used to register these types with Kubernetes.
	GroupVersion = schema.GroupVersion{Group: "secrets.operator.io", Version: "v1alpha1"}

	// SchemeBuilder collects our Go types so they can be registered together.
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme registers our types (SecretStore, ManagedSecret) with the manager.
	AddToScheme = SchemeBuilder.AddToScheme
)