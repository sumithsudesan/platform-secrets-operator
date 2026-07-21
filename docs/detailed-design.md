# Detailed Design

## 1. Overview

This document defines the implementation contract for the initial MVP of the secrets operator.

## 2. Custom Resource

### 2.1 ManagedSecret

Purpose:
- describe the desired Kubernetes Secret output
- carry provider selection and backend-specific config in one resource
- map one or more remote secret entries to one generated Secret

Suggested fields:
- `metadata.name`
- `metadata.namespace`
- `spec.providerType`
- `spec.targetSecretName`
- `spec.refreshInterval`
- `spec.remoteRefs`
- `spec.providerConfig`
- `spec.deletionPolicy`

Design intent:
- `ManagedSecret` expresses intent for a single desired Secret object
- the created Secret is derived from the external backend result
- the common resource shape stays stable
- provider-specific parsing of `remoteRefs` and `providerConfig` is handled inside the adapter selected by `spec.providerType`

## 3. Reconciliation Flow

1. Watch `ManagedSecret` resources.
2. Read `spec.providerType` and select the provider adapter.
3. Use the provider interface to fetch the remote secret payload.
4. Convert provider output into a Kubernetes Secret payload.
5. Compare the desired payload with the existing Secret.
6. Create or update the Secret if needed.
7. Update status conditions and emit events.

## 4. Provider Interface

The provider interface should be intentionally small:

```go
 type Provider interface {
     FetchSecret(ctx context.Context, managedSecret ManagedSecret) (map[string][]byte, error)
 }
```

Design rule:
- controller logic must not depend on Vault or AWS SDK types directly
- each provider adapter implements the same contract
- a new backend can be added by introducing a new provider adapter and registering it under a new `spec.providerType`

## 5. Secret Naming and Ownership

The generated Secret should follow a stable naming rule:
- `spec.targetSecretName` is the final Kubernetes Secret name
- Helm chart placeholders should reference that same name
- the operator sets owner references so the Secret is cleaned up when the `ManagedSecret` is removed

## 6. Status and Events

The controller should expose status conditions such as:
- `Ready`
- `SyncFailed`
- `AuthError`

Events should be emitted for:
- successful sync
- provider fetch failure
- auth failure
- manual Secret deletion and self-heal attempt

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

## 8. Demo Constraints

The first implementation should prioritize clarity and narrative over full production behavior.

Recommended MVP behavior:
- one mock provider implementation
- one `vault` adapter
- one `aws-secrets-manager` adapter
- one sample `ManagedSecret`
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
4. `feat: add ManagedSecret API type`
5. `feat: generate CRD manifests`
6. `feat: add mock provider interface`
7. `feat: implement reconcile and Secret sync logic`
8. `test: add basic validation and unit coverage`
