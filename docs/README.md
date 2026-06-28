# docs

Workspace-wide documentation. Service-specific docs live in `services/<svc>/docs/`.

## Layout

- [`architecture/`](architecture/) — system overview, bounded contexts, decision rationale.
- [`onboarding/`](onboarding/) — how to add a new service or UI to the monorepo.
- [`adrs/`](adrs/) — Architecture Decision Records.

## How this doc tree evolves

1. Anything workspace-wide → here.
2. Anything about a single service → `services/<svc>/docs/`.
3. Anything reusable across services → promote to `libraries/<pkg>/README.md`.
