# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- Bootstrap the public repository with governance documents and a Go 1.25
  standard-library-only toolchain.
- Document response-level `Suggested next prompt:` behavior and the explicit
  absence of composer insertion, network access, telemetry, transcript reads,
  and persistence.
- Track the implementation work in [issue #1](https://github.com/sitex/codex-next-prompt/issues/1).
- Implement fail-open `SessionStart` and `Stop` hooks that request one suitable
  response-level suggestion, preserve exact-output tasks, and never alter the
  composer or run a follow-up automatically.
- Add a marketplace catalog for packaged release roots; source checkouts remain
  non-installable because generated binaries are intentionally untracked.
- Add target-specific release archives, SHA-256 checksums, platform CI smoke
  coverage, release automation, and a self-contained local marketplace for
  Linux, macOS, and Windows.
- Verify the full Codex `0.149.1` lifecycle on Linux: untrusted hooks skip,
  trusted hooks run, suitable replies get a suggestion, exact `OK` output stays
  unchanged, and the stop hook completes. Windows remains CI smoke coverage;
  native macOS execution hasn't been tested locally.
