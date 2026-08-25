# Changelog

All notable changes to this project are documented here.

## [Unreleased]

- Pivot the 0.2.0 distribution from a non-discoverable plugin package to a
  standalone user skill.
- Remove the plugin manifest and local marketplace metadata.
- Package exactly `next/SKILL.md` in one deterministic ZIP.
- Add the post-release immutable `$skill-installer` source and
  `${CODEX_HOME:-$HOME/.codex}/skills/next` installation path.
- Remove agent policy metadata because Codex 0.149.1 hides the skill from the
  explicit catalog when implicit invocation is disabled.
- Document description-driven activation when the user types `$next`.

## [0.2.0] - 2026-08-26

### Breaking

- Replace automatic response suggestions with the explicit `$next` skill.
- Remove the executable runtime and all automatic invocation behavior.
- Replace platform-specific packages with one portable standalone skill ZIP and
  matching SHA-256 file.
- Require removal of the 0.1 plugin and marketplace registration before a clean
  0.2.0 standalone install.

### Added

- Generate one ready-to-send next prompt from the current conversation.
- Return two or three separate prompts only when the user faces a real fork.
- Ask for missing context instead of inventing values or emitting placeholders.
- Match the conversation's language and never execute the generated prompt.
- Add deterministic standard-library packaging, package contract tests, safe ZIP
  installation guidance, and signed release verification.

### Security

- Keep the package instruction-only, with no network access, persistence,
  conversation-record reads, model API calls, secret access, or command execution.

## [0.1.0] - 2026-08-20

- Initial experimental release with automatic response-level suggestions, an
  executable lifecycle integration, and separate operating-system packages.
- Added checksums, signed-tag verification, local marketplace installation, and
  public project governance.
