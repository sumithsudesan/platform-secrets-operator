# High-Level Design

## 1. Purpose

The operator provides a GitOps-friendly way to materialize secrets from external backends into native Kubernetes Secrets. The external backend remains the source of truth for secret values. The Kubernetes Secret is a derived output owned by the operator.

## 2. Problem Statement

Applications need secret data from systems such as Vault or AWS Secrets Manager, but those systems are not directly the desired state model that Kubernetes teams operate with. The operator bridges that gap by reconciling a custom resource into a native Secret object.

## 3. Core Architecture

External Backend
  -> SecretStore CR (connection + auth)
  -> ManagedSecret CR (fetch + delivery intent, references a SecretStore)
  -> Provider adapter selected from `SecretStore.spec.providerType`
  -> Reconcile loop in the operator
  -> Native Kubernetes Secret objects

## 4. Main Components

### SecretStore
Represents where a backend lives and how to authenticate to it.

Responsibilities:
- declare the provider type to use
- carry provider-specific connection and auth configuration
- act as a reusable, RBAC-scoped object: one `SecretStore` can be referenced by many
  `ManagedSecret`s, so credential/endpoint changes happen in one place

Why this is a separate CRD rather than a field on `ManagedSecret`: Kubernetes RBAC is
enforced per resource type, not per field. Splitting connection/auth into its own CRD lets a
platform/security team own `SecretStore` (who can read/write backend credentials) while app
teams own `ManagedSecret` (what secret they want) — this is what FR7's "per-team stores,
per-team auth, RBAC-isolated" actually requires structurally.

### ManagedSecret
Represents desired cluster output.

Responsibilities:
- reference a `SecretStore` via `spec.storeRef`
- declare the remote secret reference format required by the selected provider
- define the target Kubernetes Secret name and sync rules

### Provider Layer
A pluggable adapter model used by the controller.

Responsibilities:
- fetch a secret from the configured backend
- return structured secret data or errors
- isolate provider-specific logic from controller logic
- allow a new backend to be added as a separate provider implementation without changing the common CRD contract

The design intentionally keeps the top-level `ManagedSecret` shape stable while letting each provider implementation define the backend-specific form of `remoteRefs` and `providerConfig`.

### Controller
The reconcile engine.

Responsibilities:
- load `ManagedSecret`, resolve `spec.storeRef` to a `SecretStore`
- select the provider adapter from `SecretStore.spec.providerType`
- call the provider interface
- create or update the generated Secret
- detect conflicts when two `ManagedSecret`s target the same Secret name in a namespace
- emit status conditions and events

## 5. Important Design Boundaries

- The operator never writes back into the external backend.
- The generated Secret is derived and operator-owned.
- The operator must self-heal when the Secret is deleted manually.
- The operator must keep the last known Secret value when a backend error occurs.

## 6. MVP Scope

For the first implementation pass, the repo should focus on:
- the `SecretStore` and `ManagedSecret` CRDs
- a provider adapter model with backend-specific behavior
- one mock provider adapter for local demo (`vault` and `aws-secrets-manager` adapters are a
  follow-up, once the mock-backed reconcile loop is proven end-to-end)
- Kubernetes Secret reconciliation, including same-namespace target-name conflict detection

## 7. Success Criteria

The implementation is considered successful when:
- a `ManagedSecret` resolves its `SecretStore` and creates a Kubernetes Secret
- changing the backend data can be reflected through the reconcile loop
- a deleted Secret is recreated on the next reconcile
- two `ManagedSecret`s targeting the same Secret name in a namespace produce a
  `ConflictError` on the second one instead of silently overwriting
- provider-specific behavior is isolated behind a stable interface
