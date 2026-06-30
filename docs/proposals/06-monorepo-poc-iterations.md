# 06 — Monorepo PoC Scaffold Evolution

**Turns:** 9, 10, 13
**Phase:** Scaffold execution + iteration

The user prompt (turn 9) said:

> spawn two subagents to create platform and monorepo repos, they
> should be mostly poc with minimal file content but full dir structure,
> for services add bulksms, wirepay, dura, flow, flow-ussd, as well as
> the corresponding ui's. i want to review the setups

This file documents the **three iterations** the monorepo went through
in response to subsequent user requests (turns 10 and 11). The naming
decisions (turn 11) are covered in
[07-naming-decisions.md](./07-naming-decisions.md).

---

## Iteration 1 — Initial scaffold (turn 9)

Built by subagent `abbe4ea9e`. **343 files, 2,058 LOC Go, 25 TS/Vue,
21 YAML, 159 Markdown.**

```
monorepo/
├── README.md                         # workspace orientation
├── go.work                           # 6 modules:
│                                       ./shared
│                                       ./services/bulksms
│                                       ./services/wirepay
│                                       ./services/dura
│                                       ./services/flow
│                                       ./services/flow-ussd
├── go.work.sum
├── Makefile                          # top-level build/test/lint
├── docker-compose.dev.yml            # postgres + rabbitmq + valkey + mailpit + otel
├── .gitignore / .golangci.yml
│
├── shared/                           # shared Go libs (later renamed → libraries/)
│   ├── README.md, go.mod (module: github.com/scape-labs/monorepo/shared)
│   ├── money/money.go                # Currency + Amount, int64 minor units
│   ├── money/money_test.go
│   ├── tenant/tenant.go              # ctx-based resolver
│   ├── auditlog/emitter.go           # Event + Emitter iface
│   └── idgen/idgen.go                # flake-style 64-bit ID
│
├── services/                         # 5 Go backends
│   ├── bulksms/                      # tier 1, 58 files
│   ├── wirepay/                      # tier 0, 58 files (engine/deciders + engine/effectors + engine/providers + fraud/rules)
│   ├── dura/                         # tier 0, 51 files (lending, ledger, notifications, operations, payments)
│   ├── flow/                         # tier 1, 48 files (workflow engine)
│   └── flow-ussd/                    # tier 1, 47 files (no cmd/seed; ussd session/menu/handler)
│
├── ui/                               # 5 Nuxt 3 scaffolds (separate sibling to services/)
│   ├── bulksms-ui/
│   ├── bulksms-admin-ui/
│   ├── wirepay-ui/
│   ├── dura-ui/
│   └── flow-ui/
│
├── docs/
│   ├── README.md
│   ├── architecture/
│   │   ├── system-overview.md
│   │   ├── bounded-contexts.md
│   │   └── monorepo-decision.md
│   ├── onboarding/
│   │   ├── README.md
│   │   └── new-service.md
│   └── adrs/{README.md, 0001-monorepo.md}
│
└── .github/
    ├── CODEOWNERS
    └── workflows/
        ├── ci.yml
        └── release.yml
```

**Per-UI structure (`services/bulksms-ui/` example):**

```
bulksms-ui/
├── README.md
├── package.json                      # Nuxt 3 deps stub
├── nuxt.config.ts
├── tsconfig.json
├── app.vue
├── pages/{index.vue, login.vue}
├── components/HelloWorld.vue
├── assets/css/main.css
├── public/favicon.ico
├── .gitignore
├── Dockerfile
└── docker-compose.dev.yml
```

### Per-Go-service structure

Same 18 layered packages in all 5 services:

```
internal/
├── app/                              # wire-DI composition
│   ├── service-<name>.go
│   └── routes.go
├── entity/                           # domain types
├── handler/http/                     # HTTP handlers
├── service/                          # business logic
├── repository/                       # data access
├── messaging/                        # AMQP wiring
├── observability/                    # Prometheus + OTel
├── logging/
├── middleware/
├── feature/                          # feature flags
├── config/                           # typed config
├── cron/                             # typed cron jobs
├── rpc/                              # typed RPC clients
├── analytics/                        # warehouse event emitter
├── validation/
├── apperror/                         # error types
├── leaderlock/
└── testpg/                           # test helpers
```

Plus service-specific packages:
- `service-bulksms`: `auth`, `billing`, `credits`, `discounts`, `channels`,
  `dispatch`, `messages`, `reconciliation`
- `service-wirepay`: `engine/{deciders,effectors,providers}`, `fraud/rules`,
  `settlement`, `payments`
- `service-dura`: `lending`, `ledger`, `notifications`, `reconciliation`
- `service-flow-ussd`: `ussd/{session,menu,handler}`
- `service-flow`: `workflow/engine`

---

## Iteration 2 — UIs moved under `services/` (turn 10)

User request:

> ui services fall under services as well

Reasoning: every directory under `services/` is a deployable unit.
Mixing Go backends and Nuxt frontends as siblings makes the deployable
boundary match the directory boundary.

60 file moves across 5 UI directories. `monorepo/ui/` deleted.

```
monorepo/
├── shared/           ← shared Go libs
├── docs/
├── services/         ← 10 deployable units
│   ├── bulksms/             tier 1  Go backend
│   ├── wirepay/             tier 0  Go backend
│   ├── dura/                tier 0  Go backend
│   ├── flow/                tier 1  Go backend
│   ├── flow-ussd/           tier 1  Go backend
│   ├── bulksms-ui/          tier 1  Nuxt 3
│   ├── bulksms-admin-ui/    tier 1  Nuxt 3
│   ├── wirepay-ui/          tier 0  Nuxt 3
│   ├── dura-ui/             tier 0  Nuxt 3
│   └── flow-ui/             tier 1  Nuxt 3
└── .github/
```

---

## Iteration 3 — Final state after rename (turn 13)

Final commit on top of iteration 2. Applies:
1. `shared/` → `libraries/` (covered in [07-naming-decisions.md](./07-naming-decisions.md))
2. `services/bulksms/` → `services/service-bulksms/` (kebab-case)
3. `name: bulksms` → `name: service.bulksms` (Monzo-dot)
4. UI tier downgrade: `wirepay-ui` tier 0 → tier 1, `dura-ui` tier 0 → tier 1

```
monorepo/
├── libraries/                          ← was shared/
│   ├── money/         module: monorepo/libraries/money
│   ├── tenant/
│   ├── auditlog/
│   └── idgen/
├── services/
│   ├── service-bulksms/      tier 1  name: service.bulksms
│   ├── service-wirepay/      tier 0  name: service.wirepay
│   ├── service-dura/         tier 0  name: service.dura
│   ├── service-flow/         tier 1  name: service.flow
│   ├── service-flow-ussd/    tier 1  name: service.flow-ussd
│   ├── bulksms-ui/           tier 1  (Nuxt 3)
│   ├── bulksms-admin-ui/     tier 1  (Nuxt 3)
│   ├── wirepay-ui/           tier 1  (Nuxt 3, downgraded from 0)
│   ├── dura-ui/              tier 1  (Nuxt 3, downgraded from 0)
│   └── flow-ui/              tier 1  (Nuxt 3)
└── docs/
```

### What changed (turn 13)

| Before | After |
|---|---|
| `monorepo/shared/` | `monorepo/libraries/` |
| `monorepo/services/bulksms/` | `monorepo/services/service-bulksms/` |
| `name: bulksms` in service.yaml | `name: service.bulksms` |
| `dependencies.upstream: [- bulksms]` | `dependencies.upstream: [- service.bulksms]` |
| `kit.New("bulksms")` in main.go | `kit.New("service.bulksms")` |
| `wirepay-ui` tier 0 | tier 1 (UI inherits backend risk) |
| `dura-ui` tier 0 | tier 1 (UI inherits backend risk) |

### Updated `go.work`

```go
use (
    ./libraries
    ./services/service-bulksms
    ./services/service-wirepay
    ./services/service-dura
    ./services/service-flow
    ./services/service-flow-ussd
)
```

### What was kept intentionally

- `services/service-bulksms/` directory name (kebab-case for shell/IDE/Go)
- `module github.com/scape-labs/monorepo/services/service-bulksms` (Go module paths must be filesystem-style)
- `replace ... => ../../libraries/money` (filesystem paths)
- `import "monorepo/libraries/money"` (Go imports are filesystem paths)
- CI path filters (filesystem)
- UIs keep unprefixed names (not RPC participants)

---

## PoC verification

From session end:

```
go.work use                                 # validates
libraries/                                  # builds + vets + tests (money_test.go passes)
file count                                  # 343
Go LOC                                      # 2,058
TS/Vue                                      # 25
working tree                                # clean
commits on main                             # 2 (POC + the rename/fix)
```

**Known caveats (intentional for PoC):**

- Services' `go.mod` reference `github.com/scape-labs/platform/kit/server`
  but those don't resolve yet — TODO `// require ... replace ...` blocks
  remain.
- `libraries/` is the only module that compiles today (zero external
  deps by design).
- No OpenAPI specs.
- No Nuxt deps installed.
- No workspace-level CODEOWNERS (CODEOWNERS lives per-service).
- Per-service `service.yaml` files exist but use a simpler schema than
  the proposal in [04-bulksms-v2-canonical.md](./04-bulksms-v2-canonical.md).
