SHELL := /bin/sh

.DEFAULT_GOAL := check

VERSION ?= $(shell sed -n 's/^[[:space:]]*"version": "\([^"]*\)",*$$/\1/p' .codex-plugin/plugin.json)

.PHONY: test package check clean

test:
	python3 -m unittest discover -s tests -p 'test_*.py' -v

package:
	python3 scripts/package-release.py "$(VERSION)"

check: test package

clean:
	rm -rf dist
	find . -type d -name __pycache__ -prune -exec rm -rf {} +
