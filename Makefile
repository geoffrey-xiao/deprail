SHELL := /bin/sh

GO_VERSION := $(shell tr -d '[:space:]' < .go-version)
GO := go
BINARY := deprail
VERSION ?= development
TAG ?= unknown
COMMIT ?= unknown
LDFLAGS := -X github.com/geoffrey-xiao/deprail/internal/buildinfo.Version=$(VERSION) -X github.com/geoffrey-xiao/deprail/internal/buildinfo.Tag=$(TAG) -X github.com/geoffrey-xiao/deprail/internal/buildinfo.Commit=$(COMMIT)

.PHONY: bootstrap generate test lint build verify test-integration clean

bootstrap:
	@command -v $(GO) >/dev/null 2>&1 || { echo "error: Go $(GO_VERSION) is required" >&2; exit 1; }
	@test "$(shell $(GO) version 2>/dev/null | cut -d' ' -f3 | cut -c3-)" = "$(GO_VERSION)" || { echo "error: expected Go $(GO_VERSION)" >&2; $(GO) version >&2; exit 1; }

# No generators are committed yet; this remains the stable generation entry point.
generate:
	$(GO) generate ./...

test:
	$(GO) test ./...

lint:
	@test -z "$(shell $(GO) list -f '{{.Dir}}' ./... | xargs gofmt -l)" || { echo "error: gofmt required for:"; $(GO) list -f '{{.Dir}}' ./... | xargs gofmt -l; exit 1; }
	$(GO) vet ./...

build:
	$(GO) build -ldflags "$(LDFLAGS)" ./...

verify: bootstrap generate lint test build

test-integration:
	$(GO) test -tags=integration ./...

clean:
	$(GO) clean
