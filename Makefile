SHELL := /bin/sh

.DEFAULT_GOAL := check

GO_FILES := $(shell find . -type f -name '*.go' -not -path './.git/*')
VERSION ?= $(shell sed -n 's/^[[:space:]]*"version": "\([^"]*\)",*$$/\1/p' .codex-plugin/plugin.json)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: fmt vet test build smoke package check

fmt:
	@test -n "$(GO_FILES)" || { printf '%s\n' 'fmt: no Go source files found' >&2; exit 1; }
	gofmt -w $(GO_FILES)

vet:
	GOTOOLCHAIN=auto go vet ./...

test:
	GOTOOLCHAIN=auto go test -race -shuffle=on -count=1 ./...

build:
	mkdir -p bin
	GOTOOLCHAIN=auto go build -trimpath -o bin/codex-next-prompt ./cmd/codex-next-prompt

smoke:
	@test -x tests/smoke-posix.sh || { printf '%s\n' 'smoke: tests/smoke-posix.sh is not available' >&2; exit 1; }
	tests/smoke-posix.sh

package:
	@test -x scripts/package-release.sh || { printf '%s\n' 'package: scripts/package-release.sh is not available' >&2; exit 1; }
	scripts/package-release.sh "$(VERSION)" "$(GOOS)" "$(GOARCH)"

check:
	@test -n "$(GO_FILES)" || { printf '%s\n' 'check: no Go source files found' >&2; exit 1; }
	@test -z "$$(gofmt -l $(GO_FILES))" || { printf '%s\n' 'check: Go files are not formatted' >&2; exit 1; }
	GOTOOLCHAIN=auto go vet ./...
	GOTOOLCHAIN=auto go test -race -shuffle=on -count=1 ./...
	mkdir -p bin
	GOTOOLCHAIN=auto go build -trimpath -o bin/codex-next-prompt ./cmd/codex-next-prompt
	test -x tests/smoke-posix.sh || { printf '%s\n' 'check: tests/smoke-posix.sh is not available' >&2; exit 1; }
	tests/smoke-posix.sh
