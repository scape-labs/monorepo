# scape-labs monorepo

> Phase 2 of the scape-labs platform migration. A single repo containing every scape-labs deployable — Go backend services and the Nuxt frontends that pair with them.

## What lives here

```
monorepo/
├── libraries/      # shared Go libs (money, tenant, auditlog, idgen)
├── services/       # one directory per deployable: backend Go module OR Nuxt UI
│   ├── service-bulksms/      # tier 1 backend (Go)
│   ├── service-wirepay/      # tier 0 backend (Go)
│   ├── service-dura/         # tier 0 backend (Go)
│   ├── service-flow/         # tier 1 backend (Go)
│   ├── service-flow-ussd/    # tier 1 backend (Go)
│   ├── bulksms-ui/           # tier 1 frontend (Nuxt 3)
│   ├── bulksms-admin-ui/     # tier 1 frontend (Nuxt 3)
│   ├── wirepay-ui/           # tier 1 frontend (Nuxt 3)
│   ├── dura-ui/              # tier 1 frontend (Nuxt 3)
│   └── flow-ui/              # tier 1 frontend (Nuxt 3)
├── docs/           # workspace-wide architecture, ADRs, onboarding
└── .github/        # workspace-wide CI + release workflows
```

## Why UIs live under `services/`

Every directory under `services/` is a **deployable unit** with its own tier, its own `service.yaml` (backends) or its own `package.json` + Dockerfile (frontends), and its own CODEOWNERS. A Nuxt frontend is no different from a Go backend: it ships independently, has its own on-call, and owns its own tier. Putting them under one tree makes that explicit.

## Relationship to the platform repo

The monorepo **consumes** the [platform repo](../platform/) but does not contain it.

- Platform repo (`github.com/scape-labs/platform`) owns the framework:
  - `kit/` — server composition, DI, HTTP, messaging, observability, config.
  - `dsl/` — the `service.yaml` schema (v0.8).
  - `compiler/` — turns a `service.yaml` into Kubernetes manifests.
  - `ci-actions/` — reusable GitHub Actions for build/test/release.
- Monorepo (`github.com/scape-labs/monorepo`) owns the services:
  - Each `services/<name>/go.mod` imports `github.com/scape-labs/platform/kit/...`.
  - Each `services/<name>/service.yaml` is parsed by `platform/dsl`.
  - Releases consume `platform/ci-actions` workflows.

In production the two repos live side-by-side; in development a `replace` directive in each `go.mod` points at the local `../platform` checkout so cross-repo refactors work.

## Relationship to the (deleted) per-service repos

Before this monorepo existed, every service was a standalone repo under `~/workspace/scape-labs/<service>/`. Those repos have been removed; their git history has been (or will be) folded into this monorepo via `git subtree add`. Going forward, every commit to a service is a commit to this monorepo.

## Quick start

```bash
# build everything
make build

# test everything
make test

# run a single service
cd services/service-bulksms && make run

# bring up local dev deps (Postgres, RabbitMQ, Valkey, Mailpit, OpenTelemetry)
make dev-deps
```

## Adding a new service

See [`docs/onboarding/new-service.md`](docs/onboarding/new-service.md).

## Architecture

See [`docs/architecture/system-overview.md`](docs/architecture/system-overview.md) and the ADRs in [`docs/adrs/`](docs/adrs/).
