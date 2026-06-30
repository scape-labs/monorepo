# 01 — Monzo Layout (the inspiration)

**Turn:** 2
**Phase:** Monzo-analysis
**Source:** Rebuilt from Monzo engineering blog posts + two open-source Monzo
repos (`typhon`, `terrors`, `egress-operator`, `response`, `aws-nitro-util`).
Full detail in `/tmp/monzo-analysis/06-service-structure.md` (deleted after
the session but the gist is captured here).

## Monorepo top-level

```
monorepo/
├── go.mod, go.sum                  # single root module; vendored deps
├── vendor/
├── services/                       # ~3,000 services, one dir each
│   ├── service.payments/
│   ├── service.fraud/
│   ├── service.ledger/
│   ├── service.distrate/           # co-operative rate limiter
│   ├── service.crons/              # cron execution
│   ├── service.karpenter/
│   ├── service.mpa/                # multi-party-auth broker (Nitro Enclave)
│   ├── service.software-facts/
│   ├── service.software-excellence/
│   ├── service.iapi.software-excellence/
│   ├── service.vpn-session/
│   ├── service.github/
│   └── … (~2,990 more)
├── libraries/
│   ├── typhon/                     # RPC framework [open-sourced]
│   ├── terrors/                    # typed errors [open-sourced]
│   ├── distrate/                   # rate-limit client lib
│   ├── cronfig/                    # cron.Config{} DSL
│   ├── distsync/                   # etcd-backed distributed mutex
│   ├── kafka/                      # "queue on top of Kafka" client
│   ├── cassandra/
│   ├── config/                     # 60s-refresh typed config
│   ├── auth/                       # ServiceAccount → Vault
│   ├── analytics/                  # event emitter → BigQuery
│   └── …
├── web/
│   └── backstage/                  # Software Excellence UI
├── protos/                         # shared .proto definitions
├── egress-operator/                # K8s operator [open-sourced]
├── response/                       # Slack incident bot [open-sourced]
├── aws-nitro-util/                 # EIF builder [open-sourced]
├── tools/                          # codegen, depdiff, semgrep-rules
├── infra/                          # Terraform + Concourse
└── .github/                        # CODEOWNERS, PR templates
```

## A typical backend service — `service.payments/`

```
service.payments/
├── descriptor.yaml                 # owner team, business function, tier, on-call
│
├── main.go                         # entrypoint; registers Typhon Service
│
├── handler/                        # typhon.Service implementations
│   ├── card_authorization.go       # auth → fraud → ledger → response
│   ├── fps_outbound.go
│   ├── bacs_inbound.go
│   ├── international.go
│   └── refund.go
│
├── domain/                         # business logic, framework-free
│   ├── payment.go
│   ├── decline_reason.go           # 5 documented decline points
│   ├── authorisation.go
│   └── presentment.go
│
├── repository/                     # data-access layer
│   ├── cassandra/
│   │   ├── payment.go
│   │   └── schema.cql
│   └── ledger/
│       └── entry_set_writer.go
│
├── kafka/                          # event pub/sub
│   ├── publisher.go                # typed event structs
│   ├── subscription.go             # partition-key serialization
│   └── topics.go
│
├── cron/                           # cronfig DSL in Go
│   ├── settlement.go               # var SettlementJob = cron.Config{
│   ├── end_of_day.go               #     CronName: "settlement",
│   └── fx_reconciliation.go        #     Request:  ledger.SettleRequest{…},
│                                   #     Schedule: cron.Schedule{
│                                   #         Crontab:  "0 16 * * 1-5",
│                                   #         Timezone: "Europe/London",
│                                   #     },
│                                   #     FailureSemantics: &cron.FailureSemantics{
│                                   #         PageOnFailure: true,
│                                   #         RunbookURL:    "...",
│                                   #     },
│                                   # }
│
├── feature/flags.go                # feature flags read at runtime
├── config/                         # typed config structs, refreshed every 60s
│
├── manifests/                      # code-as-config for ops
│   ├── egress/
│   │   ├── external/               # one .rule per external dependency
│   │   │   ├── github.com:443.rule
│   │   │   ├── api.openbanking.org:443.rule
│   │   │   ├── mastercard-iso:443.rule
│   │   │   ├── fps-scheme:443.rule
│   │   │   └── bacs-swift:443.rule
│   │   └── internal/               # one .rule per internal allowed caller/callee
│   │       ├── service.fraud.rule
│   │       ├── service.ledger.rule
│   │       ├── service.crons.rule
│   │       └── service.distrate.rule
│   ├── deploy/                     # Kustomize base + Argo Rollouts
│   │   ├── kustomization.yaml
│   │   ├── deployment.yaml         # 2 containers: service + envoy sidecar
│   │   ├── service.yaml
│   │   ├── rollout.yaml            # canary + Prometheus analysis
│   │   ├── serviceaccount.yaml     # K8s SA → Vault dynamic creds
│   │   └── prometheusrule.yaml
│   └── prometheus/alerts.yaml
│
├── rpc/                            # proto-generated clients for downstream services
│   ├── service.ledger/client.go
│   ├── service.fraud/client.go
│   └── …
│
├── middleware/                     # Typhon Filters
│   ├── auth.go                     # mTLS / staff identity
│   ├── logging.go
│   ├── tracing.go                  # OpenTelemetry propagation
│   ├── recovery.go                 # panic → 500 + alert
│   ├── metrics.go
│   └── rate_limit.go               # distrate wrapper
│
├── observability/                  # Prometheus metrics + analytics events
├── test/{integration,contract,fixtures}/
│
├── Dockerfile                      # multi-stage; static binary in scratch
├── README.md
└── CODEOWNERS                      # team ownership
```

## How it gets built and deployed

```
Edit code in service.payments/
  ↓
PR raised
  ├── descriptor.yaml validated
  ├── *.rule files → Calico NetworkPolicy + pod labels
  ├── cronfig structs in cron/*.go → JSON for service.crons
  ├── CODEOWNERS auto-assigns reviewers
  └── semgrep runs global convention checks
  ↓
PR merged → CI runs unit + integration + CDC tests
  → build → static Go binary in scratch Docker image
  ↓
Engineer runs `monzo deploy service.payments`
  ↓
Argo Rollouts canary + Prometheus auto-rollback
  ↓
Within ~60s, downstream services see:
  ├── service.crons picks up new cron definitions
  ├── Calico picks up new egress labels
  └── Vault picks up new dynamic DB credentials (via K8s SA annotations)
```

## Patterns identified

From the 5 thematic analyses (architecture, data/ML, security, payments,
frontend) the session surfaced:

1. **Single Go module** for the whole monorepo (~3,000 services in one module).
2. **Per-service directory with consistent shape** — `descriptor.yaml` +
   `handler/`, `domain/`, `repository/`, `kafka/`, `cron/`, `manifests/`,
   `rpc/`, `middleware/`, `observability/`.
3. **Cronfig DSL in Go** — typed `cron.Config{}` structs compile to JSON for
   the dedicated cron executor service.
4. **Egress as `.rule` files** — code-as-config; the egress-operator compiles
   them to Calico `NetworkPolicy` + pod labels.
5. **Two-container pod** — service + Envoy sidecar (for mTLS, retries,
   circuit-breaking).
6. **Vault dynamic creds via K8s SA annotations** — no static DB passwords.
7. **Argo Rollouts with Prometheus analysis** — canary with auto-rollback.
8. **Typhon RPC framework** — open-sourced; powers all inter-service calls.

## Cross-references

- ADR 0001 (adopt Go-workspace monorepo) — derived from this analysis.
- `scape-engineering-model/02-go-service-standard.md` — scape-labs's per-service
  standard, derived from this template.
