# 07 — Naming Decisions

**Turn:** 11
**Phase:** Scaffold refinement

Two naming questions came up after the initial monorepo scaffold:
1. Should we use the Monzo `service.X/` naming scheme?
2. Should we use `shared/` or rename to something else?

The user approved both decisions in turn 11 ("yes") with the structure
preserved in [06-monorepo-poc-iterations.md](./06-monorepo-poc-iterations.md).

---

## Decision 1 — `service.*` naming scheme

**Adopt it for service identity, not for directory names.**

The split:

- **Directory:** kebab-case (`service-bulksms/`) — Go module paths,
  shells, IDEs all prefer kebab.
- **Service name (`name:` field in `service.yaml`):** Monzo-dot
  (`service.bulksms`) — runtime identity for PagerDuty, Slack, DNS,
  egress rules, observability dashboards.

```
services/
├── service-bulksms/         # directory: kebab-case (Go-friendly)
├── service-wirepay/
├── service-dura/
├── service-flow/
├── service-flow-ussd/
├── bulksms-ui/              # UIs stay unprefixed — they're not RPC participants
├── wirepay-ui/
└── ...
```

Companion `service.yaml` example:

```yaml
name: service.bulksms           # ← Monzo-style runtime identity
team: messaging
descriptor:
  oncall:
    pagerduty: service-bulksms-primary
    slack: "#oncall-messaging"
egress:
  internal:
    - service: service.wirepay  # ← what you reference in egress rules
```

### Decision rationale

| | Directory (kebab-case) | Service name (dot) |
|---|---|---|
| Goes into | `go.mod` module path, shell paths, IDE imports | PagerDuty schedule, Slack channel, DNS, egress rules, observability dashboards |
| Tooling friendliness | Go, shells, IDEs all happy | Some Go tooling gets confused |
| Monzo precedent | Their open-source repos use kebab-case (`monzo/typhon`, `monzo/egress-operator`) | Their internal services use `service.X` (their runtime identity) |

> You want the **runtime identity** to match Monzo (so on-call,
> observability, dsl v0.8 validation all line up), but the
> **filesystem** to stay conventional (so every Go tool works without
> surprises).

### Special case: UIs

Don't prefix UIs. A Nuxt frontend isn't an RPC service:

- Doesn't participate in egress rules.
- Doesn't have a PagerDuty schedule.
- Doesn't show up in DNS as a service endpoint.

Putting `service.` in front would be cargo-culting.

### UI tier inheritance

Tier assignment per UI also needs revisiting. UIs shouldn't be tier 0
(system-of-record) just because they display money. The tier reflects
the backend's risk; the UI inherits. Both `wirepay-ui` and `dura-ui`
were downgraded from tier 0 to tier 1.

---

## Decision 2 — `shared/` → `libraries/`

**Yes — rename to `libraries/`.**

Reasoning:
- Monzo uses `libraries/`.
- The platform monorepo uses `kit/` so the parallel is `kit/`
  (platform) → `libraries/` (monorepo).
- `shared/` is vague.

```
monorepo/
├── services/           (unchanged)
├── libraries/          ← was shared/
│   ├── money/
│   ├── tenant/
│   ├── auditlog/
│   └── idgen/
└── docs/
```

### Renames

| Before | After |
|---|---|
| `services/bulksms/` | `services/service-bulksms/` |
| `services/wirepay/` | `services/service-wirepay/` |
| `services/dura/` | `services/service-dura/` |
| `services/flow/` | `services/service-flow/` |
| `services/flow-ussd/` | `services/service-flow-ussd/` |
| `shared/` | `libraries/` |

### Per-file changes

- `go.mod` module path
- `service.yaml` `name:` field
- `cmd/<svc>/main.go` imports
- `.github/CODEOWNERS` paths

### Companion commit

`cd1f827` in `~/workspace/scape-labs/platform/` updates `services.md`
to point at the new `services/service-*` paths.

---

## Verification

From session end:

| Check | Result |
|---|---|
| `go.work use` validates | ✓ |
| `libraries/` builds + vets + tests (money_test.go passes) | ✓ |
| File count | 343 (unchanged) |
| Go LOC | 2,058 (unchanged) |
| TS/Vue | 25 (unchanged) |
| Working tree | clean |
| Commits on main | 2 (POC + the rename/fix) |
