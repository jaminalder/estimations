# Simple Makefile for Estimation Poker
GO ?= go

# --- capture extra words after the target as the commit message ---
# e.g. make pushall "my new commit"
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
MSG  ?= $(if $(strip $(ARGS)),$(ARGS),update)

.PHONY: run dev build test fmt tidy clean tools pushall

run:
	$(GO) run ./cmd/server

# Requires: go install github.com/air-verse/air@latest
dev:
	air

build:
	$(GO) build -o bin/server ./cmd/server

test:
	$(GO) test ./... -count=1

fmt:
	$(GO) fmt ./...
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin tmp

tools:
	$(GO) install github.com/air-verse/air@latest

# Add all changes, commit with optional message, then push
pushall:
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
