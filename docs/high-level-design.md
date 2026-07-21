# High-Level Design

## 1. Purpose

The operator provides a GitOps-friendly way to materialize secrets from external backends into native Kubernetes Secrets. The external backend remains the source of truth for secret values. The Kubernetes Secret is a derived output owned by the operator.

## 2. Problem Statement

Applications need secret data from systems such as Vault or AWS Secrets Manager, but those systems are not directly the desired state model that Kubernetes teams operate with. The operator bridges that gap by reconciling a custom resource into a native Secret object.

## 3. Core Architecture

External Backend
  -> ManagedSecret CR
  -> Provider adapter selected from `spec.providerType`
  -> Reconcile loop in the operator
  -> Native Kubernetes Secret objects

## 4. Main Components

### ManagedSecret
Represents desired cluster output.

Responsibilities:
- declare the provider type to use
- carry provider-specific connection and auth configuration
- declare the remote secret reference format required by the selected provider
- define the target Kubernetes Secret name and sync rules

### Provider Interface
A pluggable abstraction used by the controller.

Responsibilities:
- fetch a secret from the configured backend
- return structured secret data or errors
- isolate provider-specific logic from controller logic
- allow a new backend to be added as a separate adapter without changing the common CRD contract

### Controller
The reconcile engine.

Responsibilities:
- load `ManagedSecret`
- select the provider adapter from `spec.providerType`
- call the provider interface
- create or update the generated Secret
- emit status conditions and events

## 5. Important Design Boundaries

- The operator never writes back into the external backend.
- The generated Secret is derived and operator-owned.
- The operator must self-heal when the Secret is deleted manually.
- The operator must keep the last known Secret value when a backend error occurs.

## 6. MVP Scope

For the first implementation pass, the repo should focus on:
- one `ManagedSecret` CRD
- provider interface abstraction
- one adapter for `vault`
- one adapter for `aws-secrets-manager`
- one mock provider adapter for local demo
- Kubernetes Secret reconciliation

## 7. Success Criteria

The implementation is considered successful when:
- a `ManagedSecret` creates a Kubernetes Secret
- changing the backend data can be reflected through the reconcile loop
- a deleted Secret is recreated on the next reconcile
- provider-specific behavior is isolated behind a stable interface
