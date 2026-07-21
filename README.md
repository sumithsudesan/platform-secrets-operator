# secrets-operator

Kubernetes operator that syncs secrets from external backends into the cluster as native Kubernetes Secrets.

Status: Under active development — not production ready

## Overview
The operator is designed to be backend-agnostic. It watches `ManagedSecret` resources, reads secret data from an external provider such as Vault or AWS Secrets Manager, and materializes the result as a native Kubernetes `Secret` object.

### Why two CRDs?
- `SecretStore` defines how to connect to the backend and authenticate.
- `ManagedSecret` defines what should be synced into the cluster and how it should be managed.

This separates connection/auth concerns from desired secret output.

### Example use cases
- Multi-secret orchestration: split a single remote secret into multiple Kubernetes Secrets or key sets.
- Cross-namespace sync: push the same external secret into multiple namespaces.
- Combined flow: split one backend secret and replicate the derived outputs across environments.

## High-level architecture

```mermaid
flowchart LR
    A[External Backend<br/>Vault / AWS Secrets Manager / Provider] --> B[ManagedSecret CR<br/>Desired state in GitOps]
    B --> C[Operator Reconcile Loop]
    C --> D[Kubernetes Secret<br/>Derived output]
    D --> E[Workloads / Pods]
```

## Reconcile flow

```mermaid
flowchart TD
    R1[Read ManagedSecret] --> R2[Resolve SecretStore]
    R2 --> R3[Call Provider Interface]
    R3 --> R4[Build desired Secret payload]
    R4 --> R5[Create or update Kubernetes Secret]
    R5 --> R6[Update status + events]
```

## Design docs
- [docs/high-level-design.md](docs/high-level-design.md)
- [docs/detailed-design.md](docs/detailed-design.md)

## Status
Under active development. Not production ready.