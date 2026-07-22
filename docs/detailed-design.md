# Detailed Design

## 1. Overview

This document defines the implementation contract for the initial MVP of the secrets operator.

## 2. Custom Resources

### 2.1 SecretStore

Purpose:
- describe where a backend lives and how to authenticate to it
- act as the reusable, RBAC-scoped connection object referenced by one or more `ManagedSecret`s

Suggested fields:
- `metadata.name`
- `metadata.namespace`
- `spec.providerType` — `vault`, `aws-secrets-manager`, or `mock`
- `spec.providerConfig` — provider-specific connection/auth config (e.g. Vault
  `endpoint`/`role`, AWS `region`/`roleArn`)

Design intent:
- credentials and endpoints live here, not in `ManagedSecret`, so rotation touches one object
  instead of every CR that uses it
- `providerConfig` shape is owned by the provider adapter selected by `spec.providerType`,
  same extension model as `ManagedSecret.spec.remoteRefs`

### 2.2 ManagedSecret

Purpose:
- describe the desired Kubernetes Secret output
- reference a `SecretStore` for connection/auth, and carry only fetch + delivery intent
- map one or more remote secret entries to one generated Secret

Suggested fields:
- `metadata.name`
- `metadata.namespace`
- `spec.storeRef` — name of a `SecretStore` in the same namespace
- `spec.targetSecretName`
- `spec.refreshInterval`
- `spec.remoteRefs`
- `spec.deletionPolicy` — `Retain` (default) or `Delete`

Design intent:
- `ManagedSecret` expresses intent for a single desired Secret object
- the created Secret is derived from the external backend result
- the common resource shape stays stable
- provider-specific parsing of `remoteRefs` is handled inside the adapter selected by the
  referenced `SecretStore.spec.providerType`

## 3. Reconciliation Flow

1. Watch `ManagedSecret` resources.
2. Resolve `spec.storeRef` to a `SecretStore` in the same namespace. If not found, set
   `StoreNotFound` and requeue.
3. Read `SecretStore.spec.providerType` and select the provider adapter.
4. Check for a same-namespace conflict: another `ManagedSecret` already targeting
   `spec.targetSecretName`. If found, set `ConflictError` on this (later) CR and stop —
   do not touch the Secret.
5. Use the provider interface to fetch the remote secret payload.
6. Convert provider output into a Kubernetes Secret payload.
7. Compare the desired payload with the existing Secret.
8. Create or update the Secret if needed, with an owner reference back to the `ManagedSecret`.
9. Apply `deletionPolicy` handling (see section 7).
10. Update status conditions and emit events.

## 4. Provider Implementation Model

The provider layer should follow a backend-specific adapter model.

Design rule:
- controller logic must not depend on Vault or AWS SDK types directly
- each provider implementation owns its own remote reference format and connection configuration
- a new backend can be added by introducing a separate provider implementation and registering it under a new `SecretStore.spec.providerType`

This allows the operator to stay generic while still supporting provider-specific semantics for `remoteRefs`, `providerConfig`, and auth behavior.

## 5. Secret Naming and Ownership

The generated Secret should follow a stable naming rule:
- `spec.targetSecretName` is the final Kubernetes Secret name
- Helm chart placeholders should reference that same name
- the operator sets owner references so the Secret is cleaned up when the `ManagedSecret` is removed
- `targetSecretName` must be unique per namespace across `ManagedSecret`s; the second CR to
  claim a name gets `ConflictError`, not a silent overwrite (see section 7)

## 6. Status and Events

The controller should expose status conditions such as:
- `Ready`
- `SyncFailed`
- `AuthError`
- `StoreNotFound` — `spec.storeRef` does not resolve to an existing `SecretStore`
- `ConflictError` — another `ManagedSecret` in the namespace already owns `targetSecretName`

Events should be emitted for:
- successful sync
- provider fetch failure
- auth failure
- manual Secret deletion and self-heal attempt
- conflict detected

## 7. Failure Handling

### Backend outage
- keep the last synced Secret
- mark status as `SyncFailed`
- retry with backoff

### Auth error
- do not interpret as not found
- mark status as `AuthError`
- emit a distinct event

### Secret deleted manually in cluster
- next reconcile recreates the Secret
- status resets to `Ready` if sync succeeds

### Remote key missing in backend (`deletionPolicy`)
- `Retain` (default): keep the existing Kubernetes Secret intact, set `SecretNotFoundInBackend`,
  alert — never delete on `Retain`
- `Delete`: track a `status.notFoundCount` per `ManagedSecret`; only delete the Kubernetes
  Secret once the backend has returned a confirmed NotFound for **N consecutive reconciles
  spanning at least 2 refresh intervals** (a single transient NotFound must not trigger
  deletion — resets `notFoundCount` to 0 on any successful fetch)

### Two ManagedSecrets target the same Secret name
- first CR (by resource creation timestamp) owns the Secret
- second CR is set to `ConflictError` and the controller does not create, update, or delete
  the Secret on its behalf
- resolved automatically once the conflicting `targetSecretName` is changed or the earlier CR
  is deleted

### CR deleted while backend is down
- finalizer must not block deletion indefinitely
- timeout after 5 minutes and remove the finalizer, logging that cleanup could not be confirmed

## 8. Demo Constraints

The first implementation should prioritize clarity and narrative over full production behavior.

Recommended MVP behavior:
- one mock provider implementation (`vault` and `aws-secrets-manager` adapters are a
  follow-up, after the mock-backed reconcile loop is proven end-to-end)
- one sample `SecretStore` (type `mock`)
- one sample `ManagedSecret` referencing it
- one sample workload using the generated Secret name

## 9. Suggested Repository Structure

- `api/` for CRD types
- `controllers/` for controller logic
- `internal/providers/` for provider implementations
- `config/samples/` for demo YAML
- `docs/` for architecture and design notes

## 10. Recommended Initial Commit Sequence

1. `chore: scaffold design-first repo`
2. `docs: add high-level architecture design`
3. `docs: add detailed implementation design`
4. `feat: add SecretStore and ManagedSecret API types`
5. `feat: generate CRD manifests`
6. `feat: add mock provider interface`
7. `feat: implement reconcile and Secret sync logic`
8. `test: add basic validation and unit coverage`
