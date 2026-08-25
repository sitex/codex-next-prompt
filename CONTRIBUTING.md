# Contributing to Codex Next Prompt 0.2.0

Codex Next Prompt is a skill-only plugin. Its public action is explicit `$next`
invocation, and its output is ready-to-send prompt text based on the current
conversation.

## Scope

- Keep the repository instruction-only and local-only.
- Preserve explicit `$next` invocation. Don't add automatic invocation or an
  automatic response footer.
- Use only current conversation context. Don't add network access, persistence,
  conversation-record reads, model API calls, executable runtime code, or prompt
  execution.
- Return one default prompt. Use two or three separate prompts only for a real
  fork with materially different paths.
- Ask for missing context instead of inventing values or writing placeholders.
- Keep generated packages, credentials, private signing keys, user data, and
  local planning artifacts out of the repository.

## TDD and Validation

Start contract changes with a focused failing test. Assert routing fields,
package structure, or minimal contract terms, not prose snapshots. Make the
smallest change that passes, then refactor without weakening the test.

Run the narrowest relevant test first, followed by:

```sh
make test
make clean
make check
```

Review the generated ZIP contents and checksum. Don't hide failures with skips,
ignored errors, or a weaker assertion.

## Documentation and Safety

Update current documentation and `CHANGELOG.md` when the public contract changes.
Examples must show explicit `$next` use and must not imply that the skill performs
the recommended action. Installation changes need checksum verification and safe
ZIP path inspection.

## Release

The manifest version is the release source of truth. Produce one portable ZIP
and its `.sha256` file with `make package` or `make check`.

Tags must use `vMAJOR.MINOR.PATCH`, be annotated and signed, and match the
manifest version. The approved public key is
[`RELEASE_SIGNING_KEY.asc`](RELEASE_SIGNING_KEY.asc), fingerprint
`BFA7D43C126EE54A5FC8DD0EBE645A3EFA752D77`. The key bundle contains one primary
key. Never commit private key material.

## Pull Requests

Describe the user-visible contract, scope boundaries, tests run, generated
package inspection, and related issue. Keep the change focused and don't include
generated ZIP files.
