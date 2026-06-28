# dura

> TODO: one-paragraph description of dura — what problem it solves, who its
> primary users are, and what tier it sits at.

## Quick start

```bash
make dev-deps        # in repo root
make run             # runs cmd/dura
```

## Layout

- `cmd/dura/` — entrypoint.
- `cmd/seed/` — one-shot data seeder (CI fixtures).
- `internal/app/` — wire-DI composition root.
- `internal/service/` — business logic + Postgres store.
- `internal/<feature>/` — feature-specific package(s).
- `migrations/` — SQL migrations (run by `platform/migrator`).
- `docs/` — service-level docs.

## See also

- Workspace docs: `../../../docs/architecture/system-overview.md`
- service.yaml: `./service.yaml`
