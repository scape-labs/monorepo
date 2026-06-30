# Monorepo Proposal History

Chronological record of every package-structure proposal produced during the
design session that bootstrapped this repo.

**Source session:** Kimi CLI session `2b37db29235e4f3b3390565530c3f734` /
sub-session `6d11f4c3-b55e-405a-8b01-c8b73b9aafea` (Jun 28, 2026).

The session went through four design phases — Monzo-analysis → Monzo-pattern →
Monzo-applied-to-scape-labs → dsl-folding — each producing a more concrete
proposal than the last. The proposals are preserved here in their original
form so the reasoning behind the current structure is auditable, and so the
ADRs in `docs/adrs/` can cite the exact proposal that motivated them.

## Reading order

| # | File | Turn | What it is |
|---|---|---|---|
| 1 | [01-monzo-layout.md](./01-monzo-layout.md) | 2 | The Monzo inspiration — monorepo top-level + `service.payments/` layout + build/deploy flow |
| 2 | [02-monzo-pattern-applied.md](./02-monzo-pattern-applied.md) | 3 | Monzo pattern applied to `bulksms-v2`, `dura`, `wirepay` (with separate `descriptor.yaml` + `manifests/`) |
| 3 | [03-strategic-roadmap.md](./03-strategic-roadmap.md) | 4 | The four strategic answers — phased monorepo roadmap + god-service breakdown |
| 4 | [04-bulksms-v2-canonical.md](./04-bulksms-v2-canonical.md) | 6 | **The chosen design** — `bulksms-v2` with single `service.yaml` v0.8 (descriptor + egress folded into dsl) |
| 5 | [05-platform-poc.md](./05-platform-poc.md) | 9 | `platform/` scaffold spec (kit, dsl, protos, codegen, ci-actions, manifests) |
| 6 | [06-monorepo-poc-iterations.md](./06-monorepo-poc-iterations.md) | 9, 10, 13 | `monorepo/` scaffold evolution — initial scaffold → UIs-under-services → final state |
| 7 | [07-naming-decisions.md](./07-naming-decisions.md) | 11 | `service-*` directory scheme + `service.X` runtime identity + `shared/` → `libraries/` rename |

## Vocabulary used

The session stayed close to Monzo's own nomenclature. The vocabulary that
**did not** appear (it came later, in `main-omp`):
- "bounded context" — used only in Monzo's sense (4× total)
- "hexagonal / ports & adapters"
- "clean architecture / onion"
- "vertical slice"
- "package-by-layer"

What the session used:
- "monorepo with services" / "service-with-monorepo" / "platform + monorepo"
- "bounded context" (Monzo sense) — domain-clustering inside a service
- "god-service" / "fat Service" / "Service struct with many methods"
- "descriptor" / "egress rules" / "manifests"
- "dsl v0.8" / `service.yaml` / "deployment descriptor"
- "Adaptor → Deciders → Effectors" (Monzo's international payments pattern)
- "reactive control network" (Monzo's fraud platform)
- "cronfig" / "typed Go declarations" / "code-as-config"

## Key decisions that emerged

| Turn | Decision |
|---|---|
| 4 | Monorepo = yes, staged. Phase 1 = platform repo. Phase 2 = Go-workspace services monorepo. |
| 4 | `descriptor` + `egress` + `service.yaml` are orthogonal, keep all three. |
| 4 | Reusable CI workflows in `ci-actions` (descriptor-validate, egress-validate, cron-lint, alerts-validate, deploy-staging, deploy-prod). |
| 4 | Decompose god-services (the `bulksms-v2/internal/account/service.go` example: 1,084 LOC, 27 methods → 6 bounded contexts × ~200 LOC). |
| 6 | **Adopt descriptor + egress into dsl** — single `service.yaml` v0.8 with three blocks. Delete separate `descriptor.yaml` + `manifests/`. |
| 9 | Tier rules: 0 = system-of-record (24×7 oncall, all descriptor fields), 1 = user-facing, 2 = internal-only. |
| 10 | UIs are deployables → they live under `services/` (not a sibling `ui/`). |
| 11 | Kebab-case directories (`service-bulksms/`) + Monzo-dot runtime identity (`name: service.bulksms`). |
| 11 | `shared/` → `libraries/` (matches Monzo + parallels platform's `kit/`). |
| 13 | UI tier inheritance — UI cannot be tier 0 even if it displays money. |

## Authoring notes

- Annotations in the trees (`← NEW`, `← MOVED`, `PRESERVE`, `DELETE`) come
  from the original session and are preserved verbatim.
- Tree 9 is the most detailed (full `bulksms-v2` with every annotation).
  It is the canonical proposal — every other file either precedes or
  scaffolds toward it.
- The PoC scaffolds (Trees 10–15) used simplified stubs, not the full
  annotations. The differences between proposal and scaffold are real and
  mostly intentional simplifications for PoC scope.

## Source

The full transcript is preserved in
`/Users/jamesdube/.kimi/sessions/2b37db29235e4f3b3390565530c3f734/6d11f4c3-b55e-405a-8b01-c8b73b9aafea/`,
captured to `/tmp/kimi_turns.txt`. Cross-references from this history
back to specific turns use the turn numbers listed above.
