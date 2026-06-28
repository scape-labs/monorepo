# System overview

> TODO: fill in once the platform repo's kit/ abstractions are stable enough to describe without hand-waving.

## At a glance

```
                ┌───────────────────────────────────────────────┐
                │                  scape-labs                  │
                └───────────────────────────────────────────────┘
                                       │
        ┌──────────────────────────────┼──────────────────────────────┐
        │                              │                              │
   ┌────▼─────┐                  ┌─────▼─────┐                  ┌─────▼─────┐
   │ platform │                  │ monorepo  │                  │  infra    │
   │  repo    │                  │   repo    │                  │   repo    │
   ├──────────┤                  ├───────────┤                  ├───────────┤
   │ kit/     │ ◄── consumed ──► │ services/ │ ── deployed ──►  │ argocd/   │
   │ dsl/     │                  │ (backends │                  │ terraform/│
   │ compiler/│                  │  + UIs)   │                  │  k8s/     │
   │ ci-actions│                 │libraries/ │                  │           │
   └──────────┘                  └───────────┘                  └───────────┘
```

## Services

UIs are co-located with their backends under `services/<name>-ui/` (Nuxt 3). Both backends and frontends are first-class "services" — they ship independently, have their own tier, and own their own CODEOWNERS.

### Backends (Go)

| Service                 | Tier | Owns                                          |
| ----------------------- | ---- | --------------------------------------------- |
| `service.bulksms`       | 1    | Bulk messaging (SMS / email / WhatsApp).      |
| `service.wirepay`       | 0    | Money movement, system of record for funds.   |
| `service.dura`          | 0    | Lender of record (credit, disbursement).      |
| `service.flow`          | 1    | Workflow orchestration.                       |
| `service.flow-ussd`     | 1    | USSD text-based interface to `service.flow`.  |

### Frontends (Nuxt 3)

| Service             | Tier | Owns                              |
| ------------------- | ---- | --------------------------------- |
| `bulksms-ui`        | 1    | Customer-facing Sendai dashboard. |
| `bulksms-admin-ui`  | 1    | Operator / admin dashboard.       |
| `wirepay-ui`        | 1    | Merchant payment console.         |
| `dura-ui`           | 1    | Loan officer / borrower portal.   |
| `flow-ui`           | 1    | Workflow builder / monitor.       |

Tier 0 = system of record. Tier 1 = user-facing. See `bounded-contexts.md`.

## ⚠️ Services that own money

`service-wirepay` and `service-dura` are **tier 0**. They are the system of
record for funds movement and credit respectively. **Every deploy, schema
migration, config change, or CODEOWNERS edit to either of them touches PII and
the ledger.**

Before touching them:

- **Two-person review is mandatory on every PR** — the workspace
  `.github/CODEOWNERS` enforces this; do not bypass it.
- Schema changes to `services/service-wirepay/{ledger,payments,settlement,auth,billing}`
  or `services/service-dura/{ledger,lending,billing,auth,notifications}` must
  ship with a written **backout plan** (forward + reverse migration both
  exercised in staging before merge).
- **No direct DB writes** outside the owning service. If another service needs
  the data, subscribe to its AMQP events (see `bounded-contexts.md`).
- Treat **migration order as load-bearing** — a half-applied forward migration
  on a money table is an incident, not a bug.
- Anything money- or PII-shaped must use the helpers in `libraries/money/` and
  `libraries/auditlog/` — no bespoke types, no `float64` for currency.

If you're not sure whether your change qualifies as "touching the ledger",
ask in `#scape-platform` before opening the PR.
