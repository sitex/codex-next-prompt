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
and be signed by the public key in
[`RELEASE_SIGNING_KEY.asc`](RELEASE_SIGNING_KEY.asc), fingerprint
`BFA7D43C126EE54A5FC8DD0EBE645A3EFA752D77`. Public key rotation requires a
reviewed code change that updates the key and expected fingerprint together.
The committed key bundle must contain exactly one primary key. Never commit a
private signing key.

## Pull Requests

Describe the user-visible behavior, scope boundaries, tests run, and any
known limitations. Keep each pull request focused and link the relevant issue.
Do not include generated binaries, release archives, credentials, transcripts,
or unrelated repository files.
