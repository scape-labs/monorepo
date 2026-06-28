# shared

Internal Go libs shared across the services in this monorepo.

This is **not** the same as the platform repo. The platform repo owns the
framework (`kit/`, `dsl/`, `compiler/`). This `shared/` tree owns domain types
that are common to multiple services but not generic enough to belong in `kit/`.

## Layout

| Package       | Purpose                                              |
| ------------- | ---------------------------------------------------- |
| `money/`      | Money type with currency + minor units (no float).   |
| `tenant/`     | Tenant resolver — every request carries a tenant.    |
| `auditlog/`   | Emits audit events to the AMQP audit exchange.       |
| `idgen/`      | Flake-style 64-bit ID generator.                     |

## When to add to `shared/` vs `services/<svc>/`

- **Add to `shared/`** when ≥2 services need the same type/function and the dependency would otherwise be a copy-paste.
- **Keep in the service** when it's clearly local to one bounded context.
- **Promote to platform `kit/`** when it's generic enough to be reusable by non-scape-labs consumers of the platform.
