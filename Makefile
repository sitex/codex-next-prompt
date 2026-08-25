SHELL := /bin/sh

.DEFAULT_GOAL := check

GO_FILES := $(shell find . -type f -name '*.go' -not -path './.git/*')

.PHONY: fmt vet test build smoke package check

fmt:
	@test -n "$(GO_FILES)" || { printf '%s\n' 'fmt: no Go source files found' >&2; exit 1; }
	gofmt -w $(GO_FILES)

vet:
	go vet ./...

test:
	go test -race -shuffle=on -count=1 ./...

build:
	go build ./cmd/...

smoke:
	@test -x tests/smoke-posix.sh || { printf '%s\n' 'smoke: tests/smoke-posix.sh is not available' >&2; exit 1; }
	tests/smoke-posix.sh

package:
	@test -x scripts/package-release.sh || { printf '%s\n' 'package: scripts/package-release.sh is not available' >&2; exit 1; }
	scripts/package-release.sh

check:
	@test -n "$(GO_FILES)" || { printf '%s\n' 'check: no Go source files found' >&2; exit 1; }
	@test -z "$$(gofmt -l $(GO_FILES))" || { printf '%s\n' 'check: Go files are not formatted' >&2; exit 1; }
	go vet ./...
	go test -race -shuffle=on -count=1 ./...
	go build ./cmd/...
	test -x tests/smoke-posix.sh || { printf '%s\n' 'check: tests/smoke-posix.sh is not available' >&2; exit 1; }
	tests/smoke-posix.sh
