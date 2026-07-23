# providers

Backend-agnostic adapter layer used by the `ManagedSecret` controller to fetch
secret values. Controller code depends only on the types in this package —
never on a specific backend SDK (Vault, AWS, etc.) directly.

## Structure

- `provider.go` — the `Provider` and `Client` interfaces every backend must
  implement, plus `LocalKeyFor` (resolves a `RemoteRef`'s output key: `LocalKey`
  → `Key` → `SecretID`).
- `error.go` — `AuthError` and `NotFoundError`, typed errors so callers can tell
  "credentials rejected" apart from "value doesn't exist" instead of matching on
  error message text. Anything else is treated as a generic failure.
- `registry.go` — `Register`/`Get`/`Registered`, a lookup table keyed by
  `SecretStore.spec.providerType` (`vault`, `aws-secrets-manager`, `mock`), plus
  `FetchFromStore`, a convenience wrapper for the connect → fetch → close
  sequence.

## How it fits together

```
SecretStore.spec.providerType (string)
        |
        v
   registry.Get(providerType) -> Provider
        |
        v
   provider.NewClient(ctx, store) -> Client   (this is the "connect" step)
        |
        v
   client.GetSecrets(ctx, refs) -> map[string][]byte
        |
        v
   client.Close()                              (disconnect)
```

`FetchFromStore(ctx, store, refs)` does all four steps in one call.

## Adding a new provider

To add a new backend (e.g. `vault`, `aws-secrets-manager`):

1. Create a new file in this package, e.g. `vault.go`.
2. Define a struct implementing `Provider` (`NewClient`) and a second struct
   implementing `Client` (`GetSecrets`, `Close`).
3. In `NewClient`, read the backend-specific config off
   `store.Spec.<YourBackend>` (e.g. `store.Spec.Vault`), connect/authenticate,
   and return a `Client` bound to that connection.
4. In `GetSecrets`, resolve each `RemoteRef` against the backend. Wrap failures
   as `*AuthError` or `*NotFoundError` where the cause is known; return a plain
   error otherwise.
5. Self-register in the file's own `init()`:

   ```go
   func init() {
       Register("vault", &vaultProvider{})
   }
   ```

That's the entire extension point — no changes needed to `registry.go`,
`provider.go`, or the controller. Compiling the new file into the binary is
enough to make the provider available.

**Note:** each provider must register under a unique `providerType` string —
`Register` panics on a duplicate key, since that can only mean two files are
registering the same value (a bug in our own code, not a runtime condition).
