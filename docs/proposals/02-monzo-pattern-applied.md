# 02 — Monzo Pattern Applied to scape-labs Repos

**Turn:** 3
**Phase:** Monzo-pattern (with separate `descriptor.yaml` + `manifests/`)

The Monzo service layout from proposal 01, applied to the three
scape-labs repos. This was the first attempt at mapping Monzo onto
scape-labs's existing codebase. It still used **separate** files —
`descriptor.yaml` (service metadata) and `manifests/egress/*.rule`
files (network egress rules). In proposal 04, both are folded into
the single `service.yaml` v0.8.

## 1. `~/workspace/scape-labs/bulksms-v2` — Sendai SMS Gateway

**Today:** monolith with 30+ `internal/*` packages, 5 CLI binaries, 3 staging
environments (k6/litmus/wiremock), email templates. No `descriptor.yaml`, no
`manifests/`, no Prometheus alerts, no cron DSL.

**Recommendation:** keep as a single service, add the Monzo meta-files.

```
bulksms-v2/
│
├── descriptor.yaml                       ← NEW: owner team, tier, on-call, business function
│                                            (powers SRE/Scorecard/CODEOWNERS)
├── service.yaml                          ← keep (already scape-labs convention; bridges to K8s)
│
├── go.mod                                # module bulksms
├── go.sum
├── docker-compose.yaml                   # dev deps
├── Dockerfile                            # multi-stage → static binary in scratch
├── Makefile                              ← add: descriptor-validate, manifests-generate, cron-lint
├── CLAUDE.md
│
├── cmd/                                  # keep all 5 binaries
│   ├── bulksms/main.go                   # primary service entrypoint
│   ├── seed/main.go                      # data seeding
│   ├── csvgenerator/                     # bulk-upload CSV → batches
│   ├── holdctl/main.go                   # operator hold/release CLI
│   └── chaosctl/main.go                  # chaos engineering controller
│
├── internal/                             # REORGANISE into Monzo split
│   │
│   ├── app/                              # wiring only (WIRE providers, kit.New, lifecycle)
│   ├── observability/                    # Prometheus metrics, OTel spans
│   ├── logging/                          # structured slog
│   ├── auth/                             # JWT/API token validation
│   ├── authz/                            # RBAC + tenant scoping
│   ├── validation/                       # input validation
│   │
│   ├── entity/                           # framework-free domain types
│   │   ├── message.go, batch.go, account.go, tenant.go
│   │   ├── operator.go                   # Econet / NetOne / Telecel
│   │   ├── billing.go, money.go
│   │
│   ├── handler/                          # NEW: typhon-style handler layer
│   │   ├── http/                         # moves from internal/api/*
│   │   │   ├── middleware.go, respond.go, routes.go
│   │   │   ├── pagination.go, logger.go
│   │   │   ├── messages.go, batches.go, campaigns.go
│   │   │   ├── accounts.go, billing.go, invoices.go
│   │   │   ├── credits.go, discounts.go, sender_ids.go
│   │   │   ├── users.go, api_tokens.go, auth.go
│   │   │   ├── tariffs.go, admin.go
│   │   └── webhook/
│   │       ├── consumer.go
│   │       ├── retry_scheduler.go
│   │       └── handler.go
│   │
│   ├── service/                          # NEW: business-logic services (transport-agnostic)
│   │   ├── messages/, batches/, accounts/, billing/, credits/, discounts/
│   │   ├── overrides/, invoicing/, reconciliation/, audit/
│   │
│   ├── channels/                         # operator SMS/WhatsApp/Email adapters
│   │   ├── sms/{econet, netone, telecel}
│   │   ├── whatsapp/
│   │   ├── email/                        # SES/SMTP adapter
│   │   └── registry.go
│   │
│   ├── dispatch/, delivery/, store/, storage/, payments/
│   ├── leaderlock/, preflight/, webhook/, oauthclient/, testpg/
│   │
│   ├── messaging/                        # AMQP wiring
│   │
│   ├── cron/                             ← NEW: cron DSL
│   │   ├── settlement.go, batch_flush.go, reconciliation_daily.go
│   │   ├── billing_period_close.go, operator_health_check.go, webhook_retry_sweeper.go
│   │
│   ├── kafka/                            ← NEW: typed events
│   │   └── topics.go
│   │
│   ├── feature/                          ← NEW: feature flags
│   │   └── flags.go
│   │
│   ├── config/                           ← NEW: typed config
│   │   ├── config.go, defaults.go
│   │
│   ├── middleware/                       # relocated from auth/authz
│   │
│   ├── rpc/                              ← NEW: typed clients to other scape-labs services
│   │   ├── wirepay/client.go
│   │   └── dura/client.go
│   │
│   └── analytics/                        ← NEW: analytics event emitter
│       └── events.go
│
├── migrations/                           # Goose SQL (keep)
│
├── manifests/                            ← NEW: code-as-config for ops
│   ├── egress/
│   │   ├── external/
│   │   │   ├── econet-smpp:2775.rule
│   │   │   ├── netone-smpp:2775.rule
│   │   │   ├── telecel-http:443.rule
│   │   │   ├── whatsapp-business:443.rule
│   │   │   ├── mailpit:1025.rule
│   │   │   └── otel-collector:4317.rule
│   │   └── internal/
│   │       ├── service.billing.rule
│   │       ├── service.payments.rule
│   │       └── service.reconciliation.rule
│   ├── deploy/
│   │   ├── kustomization.yaml
│   │   ├── deployment.yaml               # 2 containers: service + envoy sidecar
│   │   ├── service.yaml
│   │   ├── serviceaccount.yaml
│   │   └── rollout.yaml
│   └── prometheus/
│       ├── alerts.yaml                   ← NEW
│       └── rules.yaml
│
├── staging/                              # KEEP — already Monzo-style
├── templates/                            # email templates (Maizzle)
├── tests/                                # KEEP — multi-layered test rig
│
├── docs/                                 ← NEW
│   ├── architecture.md
│   ├── runbooks/{sms-stuck, billing-mismatch, webhook-dlq-spike, operator-outage}.md
│   └── adrs/
│
└── .github/
    ├── CODEOWNERS                        # ← ensure team ownership
    └── workflows/
        ├── descriptor-validate.yml       ← NEW
        ├── manifests-validate.yml        ← NEW
        └── cron-lint.yml                 ← NEW
```

## 2. `~/workspace/scape-labs/dura` — Nano-loans Platform

**Today:** already the cleanest of the three. Modular monolith on
`scape-labs/kit`. Recommendation: minimal changes. Mostly add the Monzo
meta-files, lift the cron DSL to a `cron/` package, add `manifests/`,
add alerts.

```
dura/
│
├── descriptor.yaml                       ← NEW: owner=lending, tier=0 (lender-of-record)
├── service.yaml                          # already here — keep
├── go.mod                                # module github.com/scape-labs/dura
├── Dockerfile, Dockerfile.runtime
├── Makefile                              # add: descriptor-validate, cron-lint, manifests-generate
├── README.md
│
├── cmd/
│   └── dura/main.go                      # 38 lines → keep
│
├── api/                                  # OpenAPI 3.1 contract
│   ├── paths/, schemas/, examples/
│   └── openapi.yaml
│
├── internal/                             # KEEP bounded domains
│   ├── app/                              # wiring
│   ├── api/                              # HTTP handlers
│   ├── entity/                           # framework-free types (KEEP)
│   │   ├── borrower.go, kyc.go, product.go, loan.go
│   │   ├── application.go, payment.go, notification.go, money.go
│   │
│   ├── lending/                          # domain service
│   ├── ledger/                           # double-entry postings
│   ├── notifications/                    # SMS/email job runner
│   ├── operations/                       # support + reconciliation read models
│   ├── payments/                         # gateway adapter (includes wirepay)
│   │
│   ├── handler/                          ← NEW
│   │
│   ├── cron/                             ← NEW: cron DSL
│   │   ├── payment_overdue_scan.go, kyc_expiry_sweep.go
│   │   ├── statement_generation.go, reconciliation_daily.go, ledger_close_of_day.go
│   │
│   ├── observability/                    ← NEW: metrics, spans, audit events
│   ├── messaging/                        ← NEW: AMQP wiring
│   ├── config/                           ← NEW
│   ├── feature/flags.go                  ← NEW
│   ├── rpc/                              ← NEW: typed clients
│   │   ├── wirepay/client.go, bulksms/client.go
│   │
│   ├── analytics/                        ← NEW
│   └── middleware/
│
├── pkg/database/migrations/              # Goose migrations
│
├── manifests/                            ← NEW
│   ├── egress/
│   │   ├── external/{wirepay:443, bulksms:443, otel-collector:4317}.rule
│   │   └── internal/service.notifications.rule
│   ├── deploy/
│   └── prometheus/{alerts.yaml, rules.yaml}
│
├── docs/{discovery,architecture,product,integration,finance,mvp}/
├── research/pocs/ledger-go/main.go       # POC, could be its own repo
├── tests/{k3d, venom/{lib,scenarios}}/
│
└── .github/
    ├── CODEOWNERS                        # ← NEW
    └── workflows/
        ├── descriptor-validate.yml       ← NEW
        ├── manifests-validate.yml        ← NEW
        └── cron-lint.yml                 ← NEW
```

## 3. `~/workspace/scape-labs/wirepay` — Payments Gateway

**Today:** the most Monzo-aligned already. Has `cmd/wirepay/main.go`,
`cmd/seedadmin/main.go`; `internal/deciders/...` matches Monzo's
International Payments pattern exactly. Recommendation: smallest set
of additions.

```
wirepay/
│
├── descriptor.yaml                       ← NEW: owner=payments, tier=0
├── service.yaml                          # already here — keep
├── alerts.rules.yaml                     # already here (48 lines) — keep
├── go.mod                                # module github.com/scape-labs/wirepay
├── go.sum
├── Dockerfile
├── Makefile                              # add: descriptor-validate, cron-lint, manifests-generate
├── CLAUDE.md
│
├── cmd/
│   ├── wirepay/main.go                   # 71 lines — keep
│   └── seedadmin/main.go                 # 97 lines — keep
│
├── internal/                             # KEEP — already Monzo-like
│   ├── app/, observability/, logging/, auth/
│   ├── domain/                           # framework-free types
│   │   ├── payment.go, merchant.go, settlement_account.go, money.go
│   │
│   ├── engine/                           # Adaptor → Deciders → Effectors pipeline (KEEP)
│   │   ├── pipeline.go, outcome.go, idempotency.go
│   │
│   ├── deciders/                         # pure decision units — Monzo IP verbatim
│   │   ├── card_active.go, balance_ok.go, fraud_check.go
│   │   ├── pin_verify.go, sca_check_ussd.go
│   │   └── decider.go
│   │
│   ├── effectors/                        # action executors
│   │   ├── ledger_post.go, webhook_dispatch.go, notification.go
│   │
│   ├── providers/                        # Adaptors
│   │   ├── ecocash.go, mpesa.go, visanet.go, zimswitch.go
│   │   ├── bank_transfer.go
│   │   └── registry.go
│   │
│   ├── fraud/                            # reactive fraud platform
│   │   ├── engine.go, checker.go, helpers.go
│   │   └── rules/
│   │       ├── amount_anomaly.go, high_amount.go, new_card.go
│   │       ├── night_time.go, velocity.go (Detectors)
│   │       ├── warning_screen.go, block_high_risk.go, escalate_to_review.go (← NEW: Action Controls)
│   │
│   ├── settlement/, rails/, merchant/, kycstorage/, webhooks/
│   │
│   ├── handler/http/                     ← NEW: HTTP layer
│   │   ├── routes.go, payments.go, merchants.go, settlement.go
│   │   ├── providers.go, kyc.go, team.go, admin.go, ecocash_webhooks.go
│   │   └── middleware.go
│   │
│   ├── cron/                             ← NEW: cron DSL
│   │   ├── settlement_daily.go, microdeposit_verify_sweep.go
│   │   ├── reconciliation_hourly.go, webhook_retry_sweeper.go
│   │   ├── provider_health_check.go, merchant_payout.go
│   │
│   ├── messaging/, config/, feature/, rpc/{bulksms,dura}/client.go
│   ├── analytics/, testpg/
│
├── pkg/database/migrations/              # Goose migrations
│
├── manifests/                            ← NEW
│   ├── egress/
│   │   ├── external/{ecocash, mpesa, visanet, zimswitch, bank-transfer-api, otel-collector}.rule
│   │   └── internal/{service.fraud, service.ledger, service.deciders}.rule
│   ├── deploy/
│   └── prometheus/{alerts.yaml, rules.yaml}
│
├── docs/{architecture,api,runbooks/{provider-outage,webhook-dlq-spike,ledger-imbalance,kyc-document-failure}.md,adrs}/
├── tests/{k3d, venom/{config,fixtures,lib,scenarios}}/
├── scripts/
│
└── .github/
    ├── CODEOWNERS                        # ← NEW
    └── workflows/
        ├── descriptor-validate.yml, manifests-validate.yml, cron-lint.yml
        └── alerts-validate.yml           # ← NEW
```

## What changed vs. current

Per-repo deltas:

| Repo | Additions | Rationale |
|---|---|---|
| bulksms-v2 | `descriptor.yaml`, `manifests/{egress,deploy,prometheus}/`, `internal/{handler,service,cron,kafka,feature,config,middleware,rpc,analytics}/`, `docs/{runbooks,adrs}/`, `.github/workflows/{descriptor-validate,manifests-validate,cron-lint}.yml` | Largest gap to Monzo shape |
| dura | `descriptor.yaml`, `manifests/`, `internal/{cron,observability,messaging,config,feature,rpc,analytics}/`, runbooks, ADRs | Minimal — already cleanest |
| wirepay | `descriptor.yaml`, `manifests/`, `internal/{handler,cron,config,feature,rpc,analytics}/`, runbooks, ADRs, `alerts-validate.yml` | Smallest — already Monzo-like |

## What this proposal gets wrong (and proposal 04 fixes)

- **Three files for one concept** — `descriptor.yaml`, `service.yaml`, and
  `manifests/egress/*.rule` files all describe one service. The user
  later requested consolidating them into a single `service.yaml` v0.8.
- **No tier rules** — descriptors are uniform regardless of whether a
  service is system-of-record or internal-only.
- **No dsl validation hooks** — descriptor + egress rules are documented
  but not enforced by the compiler.

See [04-bulksms-v2-canonical.md](./04-bulksms-v2-canonical.md) for the
resolution.
