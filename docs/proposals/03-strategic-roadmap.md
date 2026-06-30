# 03 — Strategic Roadmap

**Turn:** 4
**Phase:** Monzo-pattern → scape-labs-pattern

The user's prompt at this turn asked four orthogonal questions in one
message:

> scape-labs has many services should they form a monorepo? is the
> descriptor + egress rules better than dsl's service.yaml? our services
> use ci-actions reusable workflows assess that as well... in cases where
> there are generic Service struct with many methods how would a monzo
> alternative/breakdown/setup look like?

The four answers are presented here in order.

---

## Answer 1 — Monorepo: yes, staged

### Phase 1 — "Platform monorepo" (low risk, immediate value, ~1 week)

```
github.com/scape-labs/platform          # new repo
├── kit/                                # ← from scape-labs/kit
├── dsl/                                # ← from scape-labs/dsl
├── protos/                             # shared .proto files
│   ├── payments/v1/payment_event.proto
│   ├── ledger/v1/posting.proto
│   ├── common/v1/money.proto
│   └── …
├── codegen/                            # proto → Go, cron DSL → JSON, etc.
├── ci-actions/                         # ← from scape-labs/ci-actions
└── manifests/                          # shared egress / k8s manifests
```

**Why first:** pure move of existing code, no consumer change yet, no
service merges. Unblocks centralisation of platform primitives without
touching the service layer.

### Phase 2 — Go workspace monorepo (medium risk, huge value, ~1 month)

```
github.com/scape-labs/monorepo          # new repo, absorbs services
├── go.work                             # references every module
├── platform/                           # from Phase 1
├── services/
│   ├── wirepay/
│   ├── dura/
│   ├── bulksms-v2/
│   ├── cnp-ussd/
│   ├── cronos/
│   └── …
├── protos/                             # generated code lands in each consumer
└── tools/
```

**Why second:** forces the platform/services split to be physical
(`monorepo/services/*` cannot import another service's `internal/`).
Cross-service type-checking on event contracts. Atomic cross-service
refactors in one PR.

---

## Answer 2 — descriptor + egress vs `service.yaml`

**Verdict:** orthogonal — keep all three, add `descriptor` + `egress`.

`service.yaml` (existing dsl v0.7) describes **deployment**. It's the
input to Kustomize/Argo and lives close to the `deploy/` directory.

`descriptor.yaml` (new) describes **service metadata** — owner team,
business function, tier, on-call, SLOs. It powers SRE scorecards,
CODEOWNERS, Backstage catalog. Owned by the platform team.

`manifests/egress/*.rule` files (new) describe **network egress** —
what domains this service can talk to. Owned by security + the service
team jointly. Compiled to Calico `NetworkPolicy`.

**Why keep all three:** each has a different audience and a different
review cadence. `descriptor` changes quarterly (on-call rotation,
ownership). `service.yaml` changes weekly (replicas, env vars). Egress
rules change rarely but require security review.

> *Note:* This was the proposal-3 verdict. The user overruled it in
> proposal 4 ("adopt descriptor, egress rules etc into dsl") and folded
> all three into a single `service.yaml` v0.8. See
> [04-bulksms-v2-canonical.md](./04-bulksms-v2-canonical.md).

---

## Answer 3 — Reusable CI workflows

`scape-labs/ci-actions` already exists. Expand it to add:

```
ci-actions/.github/workflows/
├── reusable-go-ci.yml                  # go test, lint, vet, build
├── reusable-dsl-validate.yml           # validates service.yaml schema
├── reusable-egress-validate.yml        # checks egress coverage
├── reusable-deploy-staging.yml         # Argo Rollouts canary to staging
├── reusable-deploy-prod.yml            # Argo Rollouts canary to prod
├── reusable-cron-lint.yml              # validates internal/cron/* AllJobs exhaustive
├── reusable-alerts-validate.yml        # validates PrometheusRule YAML
└── reusable-descriptor-validate.yml    # validates descriptor.yaml
```

Each consumer service's `.github/workflows/ci.yml` calls these via
`uses: scape-labs/ci-actions/.github/workflows/reusable-X.yml@v1`.

**Why:** single source of truth for CI policy. Adding a new validation
rule (e.g. "every tier-0 service must have a DR section") requires
updating one workflow, not N repos.

---

## Answer 4 — Generic Service struct breakdown (the "Monzo alternative")

The textbook case is `bulksms-v2/internal/account/service.go` —
**27 methods, 1,084 LOC**, 5+ bounded contexts crammed into one struct.

### The 6 triggers that say "split this service"

1. **>15 methods** on one `Service` struct.
2. **Methods cluster around 2+ distinct nouns** (account, user, sender-id,
   document — not "account CRUD").
3. **Different consistency profiles** — subaccount transfers must be
   ledger-consistent; account onboarding can be eventually consistent.
4. **Different change cadences** — sender ID approval is regulator-driven
   (quarterly); user CRUD is daily.
5. **Different teams** would want to own different methods — payments
   vs. compliance vs. support.
6. **>2 distinct downstream dependencies** — subaccount depends on
   ledger; user depends on RBAC; senderid depends on regulator.

### The 6-service decomposition

```
Before:
  internal/account/service.go    1,084 LOC, 27 methods, 1 struct

After:
  internal/account/service.go      220 LOC,  6 methods   (lifecycle only)
  internal/subaccount/service.go   180 LOC,  4 methods   (transfers; depends on ledger via RPC)
  internal/user/service.go         200 LOC,  6 methods   (own RBAC concept)
  internal/senderid/service.go     220 LOC,  5 methods   (regulator approval)
  internal/document/service.go      80 LOC,  2 methods   (KYC)
  internal/credits/service.go       60 LOC,  1 method
  ─────────────────────────────────────────
  total                            960 LOC, 24 methods, 6 structs
```

Each new service keeps the existing `internal/app/<component>.go` Wire
pattern but narrows it: `<Component>` now wraps ~200 LOC of lifecycle
wiring per bounded context, not one god-component.

### Same treatment for `messages.Service` (18 methods, 871 LOC)

```
internal/messages/compose/        Create, CreateBulk, CreateCampaign
internal/messages/validate/       loadAndValidate, validateIdentity, validateTemplate
internal/messages/filter/         filterSuppressions, filterCampaignSuppressions
internal/messages/price/          calculate, reserveMessageHolds, releaseHolds
internal/messages/dispatch/       publishDispatchJobs, GetCampaign, List*...
```

In the final scaffold (proposal 4 tree), `internal/messages/` is
preserved as a single package with `service.go` reduced to ~500 LOC
and a `validate/` sub-package added — a partial decomposition that
keeps the public surface stable.

### HTTP handler split for the god-Service breakdown

The HTTP layer also narrows naturally — split
`internal/api/accounts.go` into:

```
internal/handler/http/
├── accounts.go          # /api/v1/accounts/* (uses account.Service)
├── subaccounts.go       # /api/v1/accounts/:id/subaccounts/* (uses subaccount.Service)
├── users.go             # /api/v1/users/* (uses user.Service)
├── senderids.go         # /api/v1/sender-ids/* (uses senderid.Service)
├── documents.go         # /api/v1/accounts/:id/documents/* (uses document.Service)
└── accounts_routes.go   # mounts the 5 sub-routers on the Echo instance
```

---

## Recommended order

> The ordering matters: do #4 first (it's pure refactoring, no infra
> change), then #2 (it lights up the descriptor/egress workflow), then
> #3 (now those workflows have real things to validate), then #1
> phase 1 (the kit monorepo).

The user then executed (a) the god-service decomposition in `bulksms`
PR #104, (b) dsl issue #36 (descriptor+egress folded in, superseding
answer 2), (c) the `ci-actions` expansion, then (d) the platform +
monorepo scaffolds.
