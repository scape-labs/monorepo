# 05 — Platform PoC Scaffold

**Turn:** 9
**Phase:** Scaffold execution

User prompt:

> spawn two subagents to create platform and monorepo repos, they should
> be mostly poc with minimal file content but full dir structure, for
> services add bulksms, wirepay, dura, flow, flow-ussd, as well as the
> corresponding ui's. i want to review the setups

This is the **first** of two scaffold specs — the platform repo. Built
by subagent `ab7f528a0`. **146 files, 1,248 LOC Go, 24 YAML, 60 Markdown.**

```
platform/
├── README.md, go.mod, Makefile, docker-compose.dev.yml
│
├── kit/                              # ← from scape-labs/kit (16 sub-packages)
│   ├── README.md, go.mod (module: github.com/scape-labs/platform/kit)
│   ├── server/                       # HTTP server (Echo wrapper)
│   │   ├── server.go                 # type Server struct + NewServer() + Start()
│   │   ├── server_test.go
│   │   └── middleware.go             # RequestID, Recovery
│   ├── database/                     # sqlx wrapper
│   │   ├── database.go               # Component with DB *sqlx.DB
│   │   └── migrations.go             # Migrate(ctx, db)
│   ├── broker/                       # AMQP wrapper
│   │   ├── broker.go                 # Publisher/Consumer
│   │   ├── rabbitmq.go
│   │   └── outbox.go                 # outbox pattern
│   ├── auth/                         # JWT/API key auth
│   │   ├── authn.go                  # Middleware, Identity, HeaderExtractor
│   │   ├── authz.go                  # Enforcer, AddPolicy
│   │   └── jwt.go
│   ├── observability/                # OTel + Prometheus
│   │   ├── observability.go
│   │   ├── metrics.go
│   │   └── trace.go
│   ├── logging/                      # slog setup
│   ├── cache/                        # Redis/Valkey wrapper
│   ├── cron/                         # cronfig (typed cron DSL)
│   │   ├── config.go                 # type Config struct (CronName, Request, Schedule, FailureSemantics)
│   │   ├── schedule.go
│   │   └── failure_semantics.go
│   ├── feature/                      # typed feature flags (Flag[T any])
│   ├── ratelimit/                    # distributed rate limiter
│   ├── migrations/                   # goose migration helpers
│   ├── store/                        # generic entity store (Store[T any])
│   ├── oauth2/                       # OAuth2 client/server helpers
│   ├── redis/                        # kitredis wrapper
│   ├── ticker/                       # leader-elected ticker primitive
│   ├── config/                       # typed config (cleanenv)
│   └── internal/                     # kit-internal helpers
│       └── service.go                # type Service interface { Setup/Run/Alive/Ready/Close }
│
├── dsl/                              # ← from scape-labs/dsl
│   ├── README.md, go.mod (module: github.com/scape-labs/platform/dsl), Makefile
│   ├── cmd/dsl/main.go               # cobra root command "dsl"
│   ├── internal/
│   │   ├── parser/                   # service.yaml parser
│   │   │   ├── service.go            # type Service struct; Parse([]byte)
│   │   │   ├── descriptor.go         # ParseDescriptor block (v0.8)
│   │   │   ├── egress.go             # ParseEgress block (v0.8)
│   │   │   └── schema.go             # JSON schema for v0.8
│   │   ├── compiler/                 # service.yaml → K8s manifests
│   │   │   ├── deployment.go         # CompileDeployment
│   │   │   ├── service.go            # CompileService
│   │   │   ├── networkpolicy.go      # CompileNetworkPolicy (from egress block)
│   │   │   ├── prometheusrule.go     # CompilePrometheusRule (from descriptor.slos)
│   │   │   └── codeowners.go         # CompileCODEOWNERS (from team + reviewers)
│   │   ├── validate/                 # CI-side validation
│   │   │   ├── validate.go
│   │   │   ├── pagerduty.go          # checkScheduleExists
│   │   │   ├── slack.go              # checkChannelExists
│   │   │   ├── coverage.go           # checkEgressCoverage
│   │   │   └── runbook.go            # checkRunbookExists
│   │   └── ast/                      # typed AST
│   ├── pkg/
│   │   ├── manifest/                 # deployment.go, service.go, networkpolicy.go
│   │   └── schema/                   # v0.8.go, v0.7.go (legacy compat)
│   ├── schemas/                      # JSON Schema files
│   │   ├── v0.7.schema.json
│   │   └── v0.8.schema.json
│   ├── examples/                     # reference service.yaml files
│   │   ├── wirepay.yaml              # tier-0 reference
│   │   ├── bulksms.yaml              # tier-1 reference
│   │   └── admin-tool.yaml           # tier-2 reference
│   ├── docs/{schema-v0.7.md, schema-v0.8.md, migration-v0.7-to-v0.8.md, tier-rules.md}
│   └── tests/{conformance/, golden/}
│
├── protos/                           # shared .proto files
│   ├── buf.yaml, buf.gen.yaml
│   ├── payments/v1/{payment_event.proto, refund.proto}
│   ├── ledger/v1/{posting.proto, balance_definition.proto}
│   ├── common/v1/{money.proto, pagination.proto}
│   ├── auth/v1/identity.proto
│   └── observability/v1/event.proto
│
├── codegen/                          # proto → Go, cronfig → JSON, etc.
│   ├── README.md, go.mod
│   ├── cmd/
│   │   ├── protoc-gen-go/main.go
│   │   ├── cronfig-gen/main.go
│   │   └── dsl-gen/main.go
│   └── internal/{loader/, render/}
│
├── ci-actions/                       # ← from scape-labs/ci-actions
│   ├── README.md
│   └── .github/workflows/
│       ├── reusable-go-ci.yml
│       ├── reusable-dsl-validate.yml
│       ├── reusable-deploy-staging.yml
│       ├── reusable-deploy-prod.yml
│       ├── reusable-egress-validate.yml
│       ├── reusable-cron-lint.yml
│       ├── reusable-alerts-validate.yml
│       └── reusable-descriptor-validate.yml
│
├── manifests/                        # shared Kubernetes manifests
│   ├── base/{namespaces.yaml, network-policies.yaml, rbac/}
│   ├── monitoring/{prometheus.yaml, grafana.yaml, alertmanager.yaml}
│   └── secrets/vault-agent.yaml
│
└── docs/{architecture/, onboarding/{new-service.md, new-skill.md}, runbooks/}
```

## Key kit types

```go
// platform/kit/internal/service.go
type Service interface {
    Setup(ctx context.Context) error
    Run(ctx context.Context) error
    Alive() error
    Ready() error
    Close(ctx context.Context) error
}
```

Each consumer service's `cmd/<name>/main.go` is:

```go
func main() {
    slog.Info("starting service.bulksms")
    // TODO: kit.New("service.bulksms").Register(...).Run()
    _ = server.Component{}
}
```

## Known PoC limitations

- No actual implementation files for kit sub-packages beyond `server.go`.
  Most are interface declarations + stubs.
- `dsl/internal/compiler/*.go` files declare `Compile*` functions but
  don't have working bodies — the compiler outputs are not yet produced.
- `codegen/` is just stubs (`protoc-gen-go/main.go` returns an error).
- `protos/` has `.proto` files but no generated `.pb.go` checked in.
- `ci-actions/.github/workflows/reusable-*.yml` exist as placeholder
  workflow files; their `uses:` patterns reference `@v1` which doesn't
  exist yet.

These are intentional for a PoC. The structure is the deliverable, not
the implementations.
