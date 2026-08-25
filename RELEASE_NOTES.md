# Codex Next Prompt

This release packages Codex Next Prompt as a self-contained local marketplace
for Linux, macOS, and Windows on `amd64` and `arm64`. Download the archive and
matching `.sha256` file for one platform from the same release.

## Compatibility

- Codex `0.149.1` or newer is required.
- The extracted release directory must remain in place while its marketplace
  and plugin are installed.
- Release tags are annotated and signed; the release workflow verifies the tag
  signature and exact source commit before packaging.

## Limitations

- Suggestions are response-level text only, not composer insertion or prefill.
- The plugin does not execute the suggested prompt or run tools on its behalf.
- The plugin is local-only and has no network access, telemetry, persistence,
  transcript reads, or external model call.
- The static local marketplace does not provide remote update discovery.
