# Detailed Design

## 1. Overview

This document defines the implementation contract for the initial MVP of the secrets operator.

## 2. Custom Resources

### 2.1 SecretStore

Purpose:
- describe the backend endpoint and authentication model
- act as a shared connection definition for many managed secrets

Suggested fields:
- `apiVersion`: `secrets.operator.io/v1alpha1`
- `kind`: `SecretStore`
- `metadata.name`
- `spec.provider`: `vault` or `aws-secrets-manager`
- `spec.endpoint`
- `spec.auth`
- `spec.namespace`

Design intent:
- one store is reusable across multiple ManagedSecret objects
- the store should not carry workload-specific secret mapping

### 2.2 ManagedSecret

Purpose:
- describe the desired Kubernetes Secret output
- map one or more remote secret entries to one generated Secret

Suggested fields:
- `metadata.name`
- `metadata.namespace`
- `spec.storeRef`
- `spec.targetSecretName`
- `spec.remoteKey`
- `spec.refreshInterval`
- `spec.deliveryMode`
- `spec.deletionPolicy`

Design intent:
- `ManagedSecret` expresses intent for a single desired Secret object
- the created Secret is derived from the external backend result

## 3. Reconciliation Flow

1. Watch `ManagedSecret` resources.
2. Resolve the referenced `SecretStore`.
3. Use the provider interface to fetch the remote secret payload.
4. Convert provider output into a Kubernetes Secret payload.
5. Compare the desired payload with the existing Secret.
6. Create or update the Secret if needed.
7. Update status conditions and emit events.

## 4. Provider Interface

The provider interface should be intentionally small:

```go
 type Provider interface {
     FetchSecret(ctx context.Context, ref SecretRef) (map[string][]byte, error)
 }
```

Design rule:
- controller logic must not depend on Vault or AWS SDK types directly
- each provider adapter implements the same contract

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
- one sample `SecretStore`
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
4. `feat: add SecretStore and ManagedSecret API types`
5. `feat: generate CRD manifests`
6. `feat: add mock provider interface`
7. `feat: implement reconcile and Secret sync logic`
8. `test: add basic validation and unit coverage`
