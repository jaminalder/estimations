# Simple Makefile for Estimation Poker
GO ?= go
GOFLAGS ?=
ENVVARS := GOCACHE=$(CURDIR)/tmp/gocache GOTMPDIR=$(CURDIR)/tmp
TESTFLAGS ?= -v

# --- capture extra words after the target as the commit message ---
# e.g. make pushall "my new commit"
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
MSG  ?= $(if $(strip $(ARGS)),$(ARGS),update)

.PHONY: run dev build test fmt lint tidy clean tools pushall

run:
	$(ENVVARS) $(GO) run $(GOFLAGS) ./cmd/server

# Requires: go install github.com/air-verse/air@latest
dev:
	air

build:
	$(ENVVARS) $(GO) build $(GOFLAGS) -o bin/server ./cmd/server

test:
	$(ENVVARS) $(GO) test $(GOFLAGS) ./... -count=1 $(TESTFLAGS)

fmt:
	$(ENVVARS) $(GO) fmt $(GOFLAGS) ./...
	$(ENVVARS) $(GO) vet $(GOFLAGS) ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo 'golangci-lint not found. Run: make tools'; exit 1; }
	$(ENVVARS) GOLANGCI_LINT_CACHE=$(CURDIR)/tmp/golangci-cache XDG_CACHE_HOME=$(CURDIR)/tmp golangci-lint run

tidy:
	$(ENVVARS) $(GO) mod tidy $(GOFLAGS)

clean:
	rm -rf bin tmp

tools:
	GOTOOLCHAIN=go1.25.0 $(GO) install github.com/air-verse/air@latest
	GOTOOLCHAIN=go1.25.0 $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOTOOLCHAIN=go1.25.0 $(GO) install mvdan.cc/gofumpt@latest
	GOTOOLCHAIN=go1.25.0 $(GO) install golang.org/x/tools/cmd/goimports@latest

# Add all changes, commit with optional message, then push
pushall:
	@echo "Running fmt and lint before push..."
	@$(MAKE) fmt
	@$(MAKE) lint
	@echo "Commit message: '$(MSG)'"
	@git add -A
	@if git diff --cached --quiet; then \
		echo "No changes to commit."; \
	else \
		git commit -m "$(MSG)"; \
	fi
	@git push

# Prevent make from treating extra message words as targets
%:
	@:
