SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

GO ?= go
SVC := bulksms
BIN := bin/$(SVC)

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "bulksms Makefile — targets: build test run seed lint tidy"

.PHONY: build
build:
	@mkdir -p bin
	$(GO) build -o $(BIN) ./cmd/$(SVC)

.PHONY: test
test:
	$(GO) test ./...

.PHONY: run
run: build
	./$(BIN)

.PHONY: seed
seed:
	$(GO) run ./cmd/seed

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: docker
docker:
	docker build -t ghcr.io/scape-labs/monorepo/$(SVC):dev .
