# Workspace-level Makefile.
# Per-service Makefiles live in services/<svc>/Makefile.
# TODO: replace shell-out with proper targets once service count grows past 5.

SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

GO ?= go
MAKEFLAGS += --no-print-directory

SERVICES := service-bulksms service-wirepay service-dura service-flow service-flow-ussd
UIS := bulksms-ui bulksms-admin-ui wirepay-ui dura-ui flow-ui

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Targets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

## ─── top-level ────────────────────────────────────────────────────────────

.PHONY: build
build: ## Build every service binary.
	@for svc in $(SERVICES); do \
		echo "==> building $$svc"; \
		$(MAKE) -C services/$$svc build || exit 1; \
	done

.PHONY: test
test: ## Run tests in every service + libraries/.
	@for svc in $(SERVICES); do \
		echo "==> testing $$svc"; \
		$(MAKE) -C services/$$svc test || exit 1; \
	done
	cd shared && $(GO) test ./...

.PHONY: lint
lint: ## Run golangci-lint across the workspace.
	golangci-lint run ./...

.PHONY: tidy
tidy: ## Run `go mod tidy` in every module.
	@for svc in $(SERVICES); do \
		echo "==> tidying $$svc"; \
		cd services/$$svc && $(GO) mod tidy && cd ../..; \
	done
	cd libraries && $(GO) mod tidy && cd ..
	$(GO) work sync

.PHONY: fmt
fmt: ## gofmt -w the whole tree.
	$(GO) fmt ./...

## ─── local dev ────────────────────────────────────────────────────────────

.PHONY: dev-deps
dev-deps: ## Bring up postgres + rabbitmq + valkey + mailpit + otel via docker compose.
	docker compose -f docker-compose.dev.yml up -d

.PHONY: dev-deps-down
dev-deps-down: ## Tear down local dev deps.
	docker compose -f docker-compose.dev.yml down

## ─── UI ───────────────────────────────────────────────────────────────────

.PHONY: ui-install
ui-install: ## Install deps for every UI.
	@for ui in $(UIS); do \
		echo "==> installing $$ui"; \
		cd ui/$$ui && (npm install || pnpm install || yarn install) && cd ../..; \
	done

.PHONY: ui-build
ui-build: ## Build every UI.
	@for ui in $(UIS); do \
		echo "==> building $$ui"; \
		cd ui/$$ui && (npm run build || pnpm build || yarn build) && cd ../..; \
	done
ne
