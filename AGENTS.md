# Project Instructions

- Keep the implementation local-only, offline, and standard-library-only.
- Preserve the public response-level `Suggested next prompt:` contract.
- Do not implement composer insertion, network access, telemetry, transcript
  reads, persistence, or automatic follow-up execution.
- Follow TDD: write a failing behavioral test first, make it pass minimally,
  then refactor without weakening the contract.
- Hook stdout is a protocol stream. Emit only documented machine-readable
  output; send concise diagnostics to stderr and never mix the two.
- Run `make check` before reporting completion. Do not hide failures with
  ignored errors or silent skips.
- Keep credentials, user data, generated binaries, release archives, and local
  planning artifacts out of the repository.
- Stage explicit paths only when preparing a commit; do not use broad staging.
