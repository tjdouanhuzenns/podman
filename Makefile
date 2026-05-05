# Makefile for podman
# See docs/tutorials/podman_tutorial.md for build instructions

EXPORT_GOPATH := $(shell if [ "$(GOPATH)" != "" ]; then echo "GOPATH=$(GOPATH)"; fi)

GO ?= go
GOFMT ?= $(GO)fmt
GO_LDFLAGS ?= -w -s

# Version information
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell cat VERSION 2>/dev/null || echo "0.0.0")

# Build flags
LDFLAGS := -X github.com/containers/podman/v5/libpod/define.gitCommit=$(GIT_COMMIT) \
	-X github.com/containers/podman/v5/libpod/define.buildInfo=$(BUILD_DATE) \
	$(GO_LDFLAGS)

# Directories
BINDIR ?= $(DESTDIR)/usr/local/bin
MANDIR ?= $(DESTDIR)/usr/local/share/man
COMPLETIONSDIR ?= $(DESTDIR)/usr/share/bash-completion/completions

# Binary names
BINARY_NAME := podman
REMOTE_BINARY_NAME := podman-remote

.PHONY: all
all: binaries

.PHONY: binaries
binaries: podman podman-remote  ## Build podman and podman-remote binaries

.PHONY: podman
podman:  ## Build the podman binary
	$(GO) build \
		-ldflags "$(LDFLAGS)" \
		-tags "$(BUILDTAGS)" \
		-o bin/$(BINARY_NAME) \
		./cmd/podman

.PHONY: podman-remote
podman-remote:  ## Build the podman-remote binary
	$(GO) build \
		-ldflags "$(LDFLAGS)" \
		-tags "$(BUILDTAGS) remote" \
		-o bin/$(REMOTE_BINARY_NAME) \
		./cmd/podman

.PHONY: test
test: unit integration  ## Run all tests

.PHONY: unit
unit:  ## Run unit tests
	# Use -count=1 to disable test result caching
	$(GO) test -v -count=1 ./...

.PHONY: integration
integration:  ## Run integration tests
	$(GO) test -v -tags integration ./test/...

.PHONY: lint
lint:  ## Run linters
	golangci-lint run ./...

.PHONY: fmt
fmt:  ## Format Go source files
	$(GOFMT) -w $(shell find . -name '*.go' -not -path './vendor/*')

.PHONY: fmt-check
fmt-check:  ## Check Go source file formatting
	@out=$$($(GOFMT) -l $$(find . -name '*.go' -not -path './vendor/*')); \
	if [ -n "$$out" ]; then \
		echo "Files not formatted:$$out"; \
		exit 1; \
	fi

.PHONY: vendor
vendor:  ## Update vendor directory
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: install
install:  ## Install podman binary
	install -d -m 755 $(BINDIR)
	install -m 755 bin/$(BINARY_NAME) $(BINDIR)/$(BINARY_NAME)

.PHONY: install.remote
install.remote:  ## Install podman-remote binary
	install -d -m 755 $(BINDIR)
	install -m 755 bin/$(REMOTE_BINARY_NAME) $(BINDIR)/$(REMOTE_BINARY_NAME)

.PHONY: clean
clean:  ## Remove build artifacts
	rm -rf bin/
	rm -f *.out

.PHONY: help
help:  ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
