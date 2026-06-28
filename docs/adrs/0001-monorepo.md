# 0001 — Adopt a Go-workspace monorepo

**Status:** Accepted
**Date:** 2026-06-28

## Context

We have 5 Go services plus 5 Nuxt UIs. Until now every service lived in its own repo. Cross-service refactors required 5 PRs and no compiler could catch broken event contracts.

## Decision

Adopt a single monorepo at `github.com/scape-labs/monorepo` that uses the Go 1.24 workspace feature (`go.work`) to group the service modules. The monorepo also hosts the Nuxt UIs. The platform repo (`kit/`, `dsl/`, `compiler/`, `ci-actions/`) remains a separate repo and is consumed via `go.mod` `replace` directives in dev.

## Consequences

- **+** Atomic cross-service refactors.
- **+** Type-checked event contracts.
- **+** Single clone, single CI entry point.
- **−** Larger repo. Mitigated with sparse checkout.
- **−** Risk of `services/<x>/` becoming the de-facto kitchen sink. Mitigated by CODEOWNERS per service.

## Alternatives considered

- **Stay polyrepo + versioned `kit/` releases.** Rejected — coordination cost grows linearly with service count.
- **Bazel / Pants.** Rejected — high upfront cost at our current size.
