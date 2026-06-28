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
