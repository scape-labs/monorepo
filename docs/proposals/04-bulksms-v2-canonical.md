# 04 — Canonical `bulksms-v2` Proposal (`service.yaml` v0.8)

**Turn:** 6
**Phase:** dsl-folding
**Status:** **Adopted** — this is the design that ships.

After seeing proposals 1–3, the user requested:

> ...my thinking was **adopting descriptor, egress rules etc into dsl** —
> thoughts on that?

The assistant agreed: a single `service.yaml` v0.8 with three blocks
(deployment + descriptor + egress) is better than three separate files
because:

1. **One schema, one compiler, one validator.**
2. **One PR** changes service shape — no fan-out across `descriptor.yaml`
   + `service.yaml` + `manifests/egress/*.rule`.
3. **Cross-references** (`egress.internal[].service` ↔ another service's
   `name:`) live in the same file, so the validator can resolve them.

The trade-off accepted: dsl grows. But it grows in one direction
(declarative service spec), not three.

---

## The `service.yaml` v0.8 schema

Replaces `descriptor.yaml` + `manifests/egress/*.rule` files.

```yaml
---
version: "0.8"
name: sendai                       # bulksms-v2's product name in code
team: messaging                    # current team-of-record

# ── Deployment (unchanged from v0.7) ──
build:
  image: ghcr.io/scape-labs/bulksms

deployment:
  base:
    replicas: 2
    ports: { http: 5389 }
    env:
      REDIS_HOST: redis
      LOG_LEVEL: info
    health:
      liveness: /healthz
      readiness: /readyz
  staging:
    env:
      DATABASE_HOST: ${resources.db.host}
      DATABASE_PASSWORD: ${resources.db.password}
      RABBITMQ_HOST: ${resources.rabbitmq.host}

# ── NEW: descriptor block (was standalone descriptor.yaml) ──
descriptor:
  description: |
    Multi-tenant SMS gateway for Zimbabwean telecom operators
    (Econet, NetOne, Telecel). Owns the multichannel message
    pipeline (SMS + Email + WhatsApp) and the billing ledger
    integration with WirePay.

  tier: 1                          # user-facing but not system-of-record for money

  oncall:
    pagerduty: sendai-primary
    slack: "#oncall-messaging"
    rotation:
      handoff: "Monday 10:00 Africa/Harare"
      coverage: business-hours

  slos:
    availability: 99.9%
    sms_p95_dispatch_latency_ms: 5000

  dependencies:
    upstream:                       # services this one calls
      - wirepay                     # billing integration
    downstream:                     # services that call this one
      - bulksms-admin-ui
      - bulksms-ui
    databases:
      - postgres://sendai-primary
      - redis://sendai-cache
      - rabbitmq://sendai-broker

  compliance:
    pii_fields: [phone, email, sender_id]
    audit_log_retention_days: 365
    regulators: [potraz]

  runbooks:
    - docs/runbooks/sms-stuck.md
    - docs/runbooks/billing-mismatch.md
    - docs/runbooks/webhook-dlq-spike.md

# ── NEW: egress block (replaces manifests/egress/*.rule files) ──
egress:
  external:
    - domain: api.wirepay.svc.cluster.local
      port: 443
      protocol: TCP
      description: WirePay billing integration
      owner: payments@scape-labs.dev
      approved_by: [secops@scape-labs.dev]
      review_by: 2027-01-15
      in_cluster: true

    - domain: smtp.resend.com
      port: 587
      protocol: TCP
      description: Resend SMTP for transactional email
      owner: messaging@scape-labs.dev
      review_by: 2027-01-15

    - domain: graph.facebook.com
      port: 443
      protocol: TCP
      description: WhatsApp Business API (Meta)
      owner: messaging@scape-labs.dev
      data_classification: confidential
      approved_by: [secops@scape-labs.dev, dpo@scape-labs.dev]
      review_by: 2027-01-15

    - domain: otel-collector.observability.svc.cluster.local
      port: 4317
      protocol: TCP
      description: OTel collector (in-cluster)
      owner: platform@scape-labs.dev
      review_by: 2027-01-15
      in_cluster: true

    - domain: rabbitmq.sendai.svc.cluster.local
      port: 5672
      protocol: TCP
      description: RabbitMQ broker (in-cluster)
      owner: messaging@scape-labs.dev
      in_cluster: true

  internal:
    - service: service.wirepay
      rpcs: [InitiatePayment, GetPayment]
      owner: messaging@scape-labs.dev
      reviewers: [team/payments]

    - service: service.billing
      rpcs: [ReserveHold, SettleHold]
      owner: messaging@scape-labs.dev

    - service: service.audit
      rpcs: [RecordEvent]
      owner: messaging@scape-labs.dev
```

### What the compiler produces

| Output | Source block | Format |
|---|---|---|
| Kubernetes Deployment | `deployment` | `apps/v1` Deployment |
| Kubernetes Service | `deployment.ports` | `core/v1` Service |
| NetworkPolicy (default deny + allow per egress rule) | `egress` | `networking.k8s.io/v1` NetworkPolicy |
| PrometheusRule (SLO alerts) | `descriptor.slos` | `monitoring.coreos.com/v1` PrometheusRule |
| CODEOWNERS entries | `team`, `egress.internal[].reviewers` | `.github/CODEOWNERS` |
| Backstage catalog | `descriptor`, `team` | `catalog-info.yaml` |
| PagerDuty schedule reference | `descriptor.oncall.pagerduty` | validated, alert if missing |

### Tier rules

| Tier | Required `descriptor` fields | Required oncall | Egress validation |
|---|---|---|---|
| **0** (system-of-record) | all `descriptor` + `dr` + `compliance` (if PII) | 24×7 | `review_by` ≤ 12 months; `approved_by` required |
| **1** (user-facing) | all `descriptor` (no `dr`) | business-hours or 24×7 | `review_by` ≤ 18 months |
| **2** (internal-only) | `description`, `tier`, `team`, `owner` only | best-effort | `review_by` ≤ 24 months |

---

## Full proposed `bulksms-v2` tree

Every line annotated with what it maps to or what's new. Replaces
proposals 1–3's separate `descriptor.yaml` + `manifests/` files.

```
bulksms-v2/
│
├── service.yaml                                  ← EXTENDED to v0.8: deployment + descriptor + egress
│
├── descriptor.yaml                               ← DELETED (folded into service.yaml)
│
├── manifests/                                    ← DELETED (egress rules live in service.yaml)
│
├── go.mod / go.sum
├── Dockerfile                                    # multi-stage Go 1.24.1 → distroless
├── Makefile                                      # add: cron-list, descriptor-validate, egress-validate
├── docker-compose.yaml                           # postgres, rabbitmq, valkey, wiremock, mailpit, otel
├── .golangci.yml / .spectral.yaml / .venomrc
├── .gitignore                                    # ← ADD: bulksms, *.log, /manifests/
├── .dockerignore
├── .env / .env.example
├── CLAUDE.md                                     # ← update to reference service.yaml v0.8
├── README.md
│
├── bulksms                                       ← DELETED (build artifact, add to .gitignore)
│
├── cmd/                                          # entrypoint binaries (PRESERVE — 5 binaries)
│   ├── bulksms/main.go                           # boots kit.New("sendai")
│   ├── chaosctl/main.go
│   ├── csvgenerator/generator.go
│   ├── holdctl/main.go
│   └── seed/
│       ├── main.go
│       └── initiate_payment.sh                   # → MOVE to tests/stepci/
│
├── internal/
│   │
│   ├── app/                                      # wire-DI composition (PRESERVE PATTERN)
│   │   ├── account.go / audit_component.go / auth.go / authorization.go
│   │   ├── batch_component.go / billing.go / credits.go / delivery.go
│   │   ├── discounts.go / dispatch.go / email.go / hold_resolver.go
│   │   ├── invoicing.go / messages.go / messaging.go / migration.go
│   │   ├── oauth.go / oauth_client.go / overrides.go / payment.go
│   │   ├── reconciliation_component.go / routes.go
│   │   ├── sms_component.go                      # ← DELETED
│   │   ├── subaccount.go / user.go / senderid.go / document.go   # NEW (bounded contexts)
│   │   ├── webhook.go
│   │   ├── cron.go / feature.go                  # NEW
│   │
│   ├── entity/                                   # framework-free domain types
│   │   ├── account.go / account_test.go
│   │   ├── api_token.go, balance.go, base.go, batch.go
│   │   ├── client.go, currency.go, document.go, event.go
│   │   ├── invoice.go, msisdn.go, multichannel.go (+ _test)
│   │   ├── onboarding.go, password_reset.go, payment.go
│   │   ├── refund.go, sender_id.go, sms.go, sms_log.go
│   │   ├── subaccount.go / user.go              # NEW: split from account.go
│   │   ├── tariff.go, transaction.go, verification.go
│   │
│   ├── handler/                                  # NEW top-level (replaces internal/api/)
│   │   └── http/
│   │       ├── routes.go / middleware.go / pagination.go
│   │       ├── respond.go (+ _test) / requests.go (+ _test)
│   │       ├── authz_helpers.go
│   │       ├── legacy_v02.go (+ _test)          # PRESERVE (1363 LOC, wire-format stability)
│   │       │
│   │       ├── accounts/{handler.go, requests.go, legacy_test.go}
│   │       ├── api_tokens/{handler.go, requests.go}
│   │       ├── auth/{handler.go, requests.go}
│   │       ├── batch/{handler.go, query.go, request.go}
│   │       ├── campaigns/handler.go
│   │       ├── credits/{handler.go, requests.go}
│   │       ├── delivery/{handler.go, requests.go}
│   │       ├── discounts/{handler.go, requests.go, requests_test.go}
│   │       ├── invoices/handler.go
│   │       ├── message_webhooks/{handler.go, requests.go}
│   │       ├── messages/{handler.go, requests.go}
│   │       ├── multichannel/{handler.go, requests.go, requests_test.go}
│   │       ├── payments/{handler.go, requests.go}
│   │       ├── sender_ids/{handler.go, requests.go}
│   │       ├── sms/{handler.go, query.go, requests.go}    # LEGACY (mark deprecated)
│   │       ├── sms_query.go / sse.go
│   │       ├── tariff_overrides/{handler.go, requests.go}
│   │       ├── tariffs/{handler.go, requests.go}
│   │       ├── users/{handler.go, requests.go}
│   │       └── subaccounts/{handler.go, requests.go}     # NEW bounded-context sub-router
│   │
│   ├── account/                                  # FOCUSED: lifecycle only
│   │   ├── service.go                            # ← FROM 1084 LOC, 27 methods → ~220 LOC, 6 methods
│   │   ├── store.go, pgstore.go
│   │   ├── consumer_welcome.go
│   │   ├── consumer_senderid.go / email.go       # ← MOVED to internal/senderid/
│   │   ├── errors.go, types.go
│   │   ├── service_credits_test.go / service_senderid_test.go   # ← MOVED
│   │   └── (remove) bulksms.log                  # ← DELETE
│   │
│   ├── subaccount/                               # NEW: split from account/
│   │   ├── service.go                            # CreateSubaccount, ListSubaccounts, TransferFunds, GetBalances
│   │   ├── store.go, pgstore.go, ledger_client.go, *_test.go
│   │
│   ├── user/                                     # NEW: split from account/
│   │   ├── service.go                            # GetUser, ListUsers, ListUsersByAccount, UpdateUser, CreateUserForAccount, AssignRole
│   │   ├── store.go, pgstore.go, rbac.go, *_test.go
│   │
│   ├── senderid/                                 # NEW: split from account/
│   │   ├── service.go                            # CreateSenderId, ListAccountSenderIds, ApproveSenderId, UpdateSenderIdStatus, ListAllSenderIds
│   │   ├── store.go, pgstore.go, consumer.go, email.go, regulator.go, *_test.go
│   │
│   ├── document/                                 # NEW: split from account/
│   │   ├── service.go                            # SubmitDocuments, ReviewDocument
│   │   ├── store.go, pgstore.go, blobstore.go, *_test.go
│   │
│   ├── credits/                                  # PRESERVE (already focused: 5 methods)
│   ├── sms/                                      # PRESERVE for now (legacy, marked deprecated)
│   ├── messages/                                 # PRESERVE (god-service: 18 methods, but bounded by workflow)
│   │   └── validate/                             # NEW: private validation helpers
│   │       ├── identity.go
│   │       └── template.go
│   │
│   ├── messaging/                                # AMQP wiring (PRESERVE)
│   │   ├── account_events.go, consumer_helpers.go (+ _test)
│   │   ├── dispatch_job.go, envelope.go, message_dispatch_job.go
│   │   ├── payment_events.go, publisher.go, sms_events.go, topics.go
│   │
│   ├── channels/                                 # CONSOLIDATED — merges channels/* + operators/* + emails/*
│   │   ├── provider.go                          # Provider interface (unified)
│   │   ├── registry.go                           # NEW
│   │   ├── sms/{provider.go, provider_test.go, operators/{bulkit, d7networks(+_test), econet, esolutions, netone, provider}.go}
│   │   ├── email/{provider.go, provider_test.go, adapters/{console, logger, mailpit, provider, resend}.go}
│   │   └── whatsapp/{common.go, d7.go(+_test), twilio.go(+_test)}
│   │
│   ├── billing/                                  # PRESERVE (already well-factored)
│   ├── batch/                                    # PRESERVE (already focused: 3 methods)
│   ├── auth/, authz/, discounts/, overrides/,    # PRESERVE
│   ├── payments/, delivery/, dispatch/, invoicing/, webhook/,
│   ├── audit/, reconciliation/, oauthclient/, observability/,
│   ├── validation/, preflight/, leaderlock/, logging/, testpg/, storage/
│   │
│   ├── http/                                     # DELETED (was empty)
│   ├── store/                                    # DELETED (dead code)
│   ├── emails/                                   # DELETED (moved to channels/email/adapters/)
│   ├── operators/                                # DELETED (moved to channels/sms/operators/)
│   ├── api/                                      # DELETED (split into handler/http/)
│   │
│   ├── cron/                                     # NEW: typed cron jobs
│   │   ├── config.go                             # cronfig.Config + Schedule/FailureSemantics types
│   │   ├── all_jobs.go                           # var AllJobs = []cronfig.Config{...} (CI-enforced exhaustive)
│   │   ├── settlement_daily.go, billing_period_close.go
│   │   ├── webhook_retry_sweep.go, delivery_portal_token_refresh.go
│   │   ├── batch_sweep_stuck.go, credits_expiry.go
│   │   ├── discounts_validate.go, reconciliation_audit.go
│   │
│   ├── config/                                   # NEW: typed config (replaces .env reads)
│   ├── feature/                                  # NEW: feature flags
│   ├── middleware/                               # NEW: Typhon-style filters
│   │   ├── auth.go, rate_limit.go, recovery.go, tracing.go, metrics.go
│   │
│   ├── rpc/                                      # NEW: typed clients to other scape-labs services
│   │   ├── wirepay/client.go
│   │   ├── dura/client.go
│   │   └── bulksms_sms/client.go
│   │
│   └── analytics/                                # NEW: warehouse event emitter
│       ├── events.go, emitter.go
│
├── migrations/                                   # goose SQL (PRESERVE — all 39 files)
├── templates/                                    # PRESERVE (paired 1:1)
│   ├── generated/                                # compiled output (committed)
│   └── maizzle/                                  # Maizzle sources (dev-only npm)
├── tests/                                        # PRESERVE (harnesses)
│   ├── chaos/, k3d/, k6/, otel/, stepci/, venom/, wiremock/
├── docs/                                         # REORGANISE by purpose
│   ├── api/{openapi.yaml, openapi-v1.yaml, openapi-internal.yaml,
│   │       api-route-ownership.md, legacy-v0.2.md, index.html, internal.html, …}
│   ├── architecture/{system-overview.md, data-flow.md, bounded-contexts.md}
│   ├── design/{batch-query-optimization.md, whatsapp-email-bulk-messaging-support.md}
│   ├── runbooks/{sms-stuck.md, billing-mismatch.md, webhook-dlq-spike.md,
│   │             econet-portal-down.md, batch-finalizer-stuck.md}
│   └── adrs/{0001-monorepo-keep-services-split.md, 0002-egress-in-dsl-v0.8.md,
│             0003-account-service-decomposition.md, 0004-handler-http-subdirs.md,
│             0005-sms-vs-messages-merger.md}
│
└── .github/
    ├── CODEOWNERS                               # NEW: derived from service.yaml descriptor.team
    └── workflows/
        ├── ci.yml                               # ← UPDATED to use scape-labs/ci-actions reusable workflows
        ├── descriptor-validate.yml              # ← service.yaml v0.8 validator
        ├── egress-validate.yml                  # ← dsl egress coverage check
        ├── cron-lint.yml                        # ← internal/cron/* AllJobs exhaustive check
        └── golangci.yml
```

---

## Mapping summary

| Action | Files | Lines affected |
|---|---|---|
| **PRESERVE as-is** | Most internal/* packages (32 of 33) | ~40,000 LOC |
| **DECOMPOSE** | `internal/account/` → 6 packages | 1,084 LOC god-service → 6 focused services |
| **SPLIT** | `internal/api/` (45 .go + 7 _test.go, 7,520 LOC) → `internal/handler/http/<resource>/` | 52 files reorganised |
| **CONSOLIDATE** | `internal/emails/` + `internal/channels/email/` → `internal/channels/email/adapters/` | 10 files merged |
| **CONSOLIDATE** | `internal/operators/` + `internal/channels/sms/` → `internal/channels/sms/operators/` | 11 files merged |
| **DELETE** | `internal/http/` (empty), `internal/store/` (dead code), `internal/emails/`, `internal/operators/`, `internal/api/`, checked-in `bulksms` binary, 4 stray `.log` files | -52 files |
| **CREATE** | `internal/subaccount/`, `internal/user/`, `internal/senderid/`, `internal/document/` | ~4 new packages |
| **CREATE** | `internal/cron/`, `internal/config/`, `internal/feature/`, `internal/middleware/`, `internal/rpc/`, `internal/analytics/` | ~6 new packages |
| **EXTEND** | `service.yaml` v0.7 → v0.8 (add `descriptor` + `egress` blocks) | 1 file |
| **REPLACE** | `descriptor.yaml` → embedded `descriptor:` block | 1 file deleted |
| **REPLACE** | `manifests/egress/*.rule` → embedded `egress:` block | 1 directory deleted |
| **REORGANISE** | `docs/` split into `api/`, `architecture/`, `design/`, `runbooks/`, `adrs/` | 7 files moved |

---

## What shipped vs. what didn't

The PoC scaffolds in proposals 5–7 used simplified stubs and did not
implement all of this proposal. Specifically:

| Proposal | PoC scaffold reality |
|---|---|
| `service.yaml` v0.8 with nested `descriptor.slos` and `descriptor.dependencies` | Scaffold flattens `slos` and `dependencies` to top level |
| `egress.external[]` and `egress.internal[]` | Scaffold reshapes `egress` around `smtp`, `http_clients`, `secrets` |
| `descriptor.oncall.{pagerduty, slack, rotation}` | Scaffold keeps `pagerduty` + `slack` only |
| `internal/{subaccount,user,senderid,document}/` decomposition | Not scaffolded — services are at proposal-2 level |
| `internal/handler/http/<resource>/` split | Scaffold keeps the older `handler/http/` flat structure |
| `internal/cron/`, `internal/rpc/`, `internal/analytics/`, `internal/feature/`, `internal/config/` | All scaffolded as README stubs |

These are intentional PoC simplifications. The proposal above is the
target; the scaffold is a step toward it.
