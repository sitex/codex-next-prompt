# Contributing

Thank you for helping improve Codex Next Prompt. Product scope and the public
behavior contract are tracked in [issue #1](https://github.com/sitex/codex-next-prompt/issues/1).

## Before Opening a Change

- Keep the plugin local-only and standard-library-only.
- Preserve the response-level `Suggested next prompt:` behavior.
- Do not add composer insertion, network access, telemetry, transcript reads,
  persistence, or automatic continuation.
- Add or update tests before implementation changes. Use a red, green, refactor
  cycle and table-driven tests for protocol behavior.
- Keep hook stdout machine-readable and reserve stderr for concise diagnostics.
- Keep public documentation free of credentials, user data, private paths, and
  local planning artifacts.

## Validation

Use Go 1.25 and run the narrowest relevant command first, followed by:

```sh
make check
```

Do not hide failures with broad skips or ignored errors. Changes that affect
the public contract should explain the behavior and update `CHANGELOG.md`.
Release tags must use strict `vMAJOR.MINOR.PATCH`-style versions, be annotated,
and be signed by a key available to the GitHub Actions runner for `git verify-tag`.

## Pull Requests

Describe the user-visible behavior, scope boundaries, tests run, and any
known limitations. Keep each pull request focused and link the relevant issue.
Do not include generated binaries, release archives, credentials, transcripts,
or unrelated repository files.
