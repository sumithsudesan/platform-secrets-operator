# Test Strategy

## 1. Current State

This repository is currently in a design-first stage. The documentation describes the architecture, the `ManagedSecret` resource model, and the provider adapter model, but the actual operator controller, CRD manifests, and runtime implementation are not yet scaffolded.

Because of that, the test strategy must be split into two layers:

1. architecture and contract validation
2. runtime validation after implementation

## 2. Design Validation

The design can be validated by reviewing the repo intent directly:

- one stable `ManagedSecret` resource shape
- `spec.providerType` selects the backend provider
- provider-specific semantics for `remoteRefs` and `providerConfig` live behind the provider adapter boundary
- the controller stays generic and backend-agnostic

This is the correct validation level for the current documentation-only milestone.

## 3. Runtime Validation Plan

Once the operator implementation is added, the runtime validation plan should include:

### 3.1 CRD Validation
- generate CRD manifests
- verify the schema for `ManagedSecret`
- confirm that the expected fields are accepted by Kubernetes

### 3.2 Mock Provider Validation
- use a mock provider for the first end-to-end test
- verify that a `ManagedSecret` produces the expected derived Kubernetes `Secret`
- verify reconcile behavior for create, update, and delete cases

### 3.3 Provider Adapter Validation
- confirm that each provider implementation can read its own `remoteRefs` and `providerConfig`
- verify backend-specific validation is handled in the provider layer

### 3.4 Integration Validation
- apply a sample `ManagedSecret` to a test cluster
- check that the derived `Secret` is created or updated
- modify the remote secret source and confirm the reconcile loop reflects the new state

## 4. What Is Not Valid Yet

The current repository does not yet contain the runtime pieces needed to execute a real operator test flow. Therefore, a test scenario that assumes a running controller and generated CRD is not a valid proof point for the current documentation-only state.

The test plan should be considered a future implementation milestone, not an immediate validation step.
