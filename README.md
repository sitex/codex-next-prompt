# Codex Next Prompt

Codex Next Prompt is a local-only Codex plugin that asks the active model to
append one concise follow-up suggestion to suitable final responses.

The supported response-level UX is:

```text
Suggested next prompt: <a specific prompt the user could submit>
```

The suggestion is advisory. It may be omitted when there is no meaningful next
action, when the response is an exact-output task, for a terse acknowledgement
or safety refusal, or while a question is waiting for user-provided data.

## Scope

The plugin uses supported Codex lifecycle hooks. Session-start guidance asks the
model for the response-level line, and a stop hook validates the completed
response without blocking it. The minimum supported Codex version is `0.149.1`.

This project does **not** provide composer insertion, ghost text, draft
mutation, prefill, Tab/Right acceptance, or any other native composer API.
There is no network access, telemetry, transcript reading, session-history
persistence, analytics, external model call, nested Codex invocation, or
automatic follow-up execution.

Phase 1 establishes the public contract and contributor toolchain. Hook source,
tests, packaging, and release automation are added in later phases.

## Development

Contributors need Go 1.25. The runtime uses only the Go standard library and
does not require Go for an eventual precompiled release archive.

```sh
make fmt
make vet
make test
make build
make smoke
make package
make check
```

The implementation targets and smoke fixtures are intentionally absent during
the repository bootstrap phase, so the corresponding targets fail fast until
their planned phases land.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development and review rules.

## Security

See [SECURITY.md](SECURITY.md) for responsible disclosure instructions.

## Project Status

Product work is tracked in [issue #1](https://github.com/sitex/codex-next-prompt/issues/1).
This repository is public and currently contains the Phase 1 bootstrap only.

## License

Copyright (c) 2026 Rocky. Released under the [MIT License](LICENSE).
