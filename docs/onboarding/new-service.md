# Adding a new service

> TODO: this is the v0.1 checklist. Once we add the 6th service we'll know what to flesh out.

## 1. Scaffold the directory

```bash
cd ~/workspace/scape-labs/monorepo
mkdir -p services/service-<name>
# copy the per-service template:
cp -R docs/templates/service-template/. services/service-<name>/
```

## 2. Wire it into the Go workspace

Add the new module to `go.work`:

```go
use (
    ./libraries
    ./services/service-bulksms
    ...
    ./services/service-<name>   # ← add this
)
```

## 3. Initialize the module

```bash
cd services/service-<name>
go mod init github.com/scape-labs/monorepo/services/service-<name>
```

Add the platform `replace` directive for local dev:

```go
replace github.com/scape-labs/platform => ../../platform
```

(Remove before tagging a release — CI substitutes the upstream version.)

## 4. Write `service.yaml`

Use the v0.8 schema from `platform/dsl`. See [`docs/templates/service.yaml`](../../docs/templates/service.yaml) for a minimal example.

## 5. Add CODEOWNERS

Append a line to `.github/CODEOWNERS` and create `services/<name>/.github/CODEOWNERS` with the team handle.

## 6. CI

The reusable workflow `scape-labs/ci-actions/.github/workflows/go-test.yml@v1` discovers services automatically as long as your service lives under `services/<name>/`. No workflow changes required.

## 7. First PR

The first PR should:
- introduce the scaffold (this is mostly mechanical);
- include an empty `internal/service/service.go` with a `TODO: implement`;
- pass CI (lint + test + build).
