# Why a monorepo (and not a polyrepo)

> TODO: turn this into a proper ADR-style writeup. For now it's a punch list of the arguments we used to convince ourselves.

## The problem

Before this monorepo, every scape-labs service was its own repo. Cross-service refactors were painful:

- Touch `kit/server` → open 5 PRs across 5 repos.
- Change a shared event shape → fan out a breaking change with no compiler help.
- Onboard a new engineer → `git clone` 5+ repos and wire up GOPATH.

## Why monorepo (with a Go workspace) wins

- **Atomic cross-service refactors.** One branch, one PR, one CI run.
- **Type-checked event contracts.** All event publishers/consumers live in the same workspace, so `go build` catches mismatches.
- **Single source of truth for shared libs.** `libraries/` replaces the old "go get github.com/scape-labs/money" pattern.
- **Local dev is one repo.** `git clone` once, get everything.
- **Go workspaces keep module boundaries honest.** Each service still has its own `go.mod`; we only "merge" them at build time via `go.work`.

## Why not Bazel / Nx / Turborepo?

- We're Go-first. Go workspaces give us 90% of the value with 10% of the toolchain.
- Bazel would force every service to declare its deps as BUILD files — high migration cost for low incremental value at our size (5 services).
- We can revisit if we hit >15 services.

## Risks we accepted

- **Larger clone.** Mitigated by sparse checkout / `--depth 1` for casual contributors.
- **CI gets slower as the tree grows.** Mitigated by changed-path detection (see `.github/workflows/release.yml`).
- **One bad `git push` breaks everyone.** Mitigated by CODEOWNERS per service.
