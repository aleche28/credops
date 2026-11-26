BIN     := bin/credops
MAIN    := ./cmd/credops
VERSPKG := github.com/aleche28/credops/internal/version

# Build metadata, injected with -ldflags -X (ADR-0002). Each can be overridden on the command
# line, e.g. `make build VERSION=v0.1.0`. DATE is the commit time in UTC, not the wall clock,
# so the same commit always builds the same binary (reproducible, verifiable releases).
# Outside a git checkout they come out empty, and the binary falls back to the build info Go
# embeds (runtime/debug.ReadBuildInfo).
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null)
COMMIT  ?= $(shell git rev-parse HEAD 2>/dev/null)
DATE    ?= $(shell TZ=UTC git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ 2>/dev/null)

LDFLAGS := -s -w \
	-X $(VERSPKG).Version=$(VERSION) \
	-X $(VERSPKG).Commit=$(COMMIT) \
	-X $(VERSPKG).Date=$(DATE)

.PHONY: help build clean test test-race fmt fmt-check vet tidy tidy-check lint docs-check ci

help:
	@echo "credops targets:"
	@echo "  build        Build $(BIN) with version metadata"
	@echo "  clean        Remove build output"
	@echo "  test         Run the test suite (go test ./...)"
	@echo "  test-race    Run the test suite with the race detector"
	@echo "  fmt          Format Go code (rewrites files)"
	@echo "  fmt-check    Fail if any file is unformatted (prints the diff)"
	@echo "  vet          Vet Go code"
	@echo "  tidy         Sync go.sum with go.mod"
	@echo "  tidy-check   Fail if go.mod/go.sum are not tidy (prints the diff)"
	@echo "  lint         Run golangci-lint"
	@echo "  docs-check   Check ADR numbering, the ADR index and relative links"
	@echo "  ci           Run every check the CI workflows run"

build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) $(MAIN)

clean:
	rm -rf bin dist

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	golangci-lint fmt

fmt-check:
	golangci-lint fmt --diff

vet:
	go vet ./...

tidy:
	go mod tidy

tidy-check:
	go mod tidy -diff

lint:
	golangci-lint run ./...

docs-check:
	python3 scripts/check_docs.py

ci: tidy-check fmt-check vet lint build test-race docs-check
