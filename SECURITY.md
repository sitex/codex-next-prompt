# Security Policy for 0.2.0

## Supported Versions

Only the latest published version receives security fixes. The current public
action is explicit `$next` invocation.

## Threat Surface

Codex Next Prompt is instruction-only. The release contains plugin metadata and
skill text, with no executable runtime. `$next` receives only the current
conversation context and returns prompt text. It doesn't execute the result.

Security reports are especially relevant if the package gains unexpected
automatic invocation, network access, persistence, conversation-record access,
model API calls, file access outside supplied context, command execution, or
secret collection.

Installation integrity also matters. Releases provide one ZIP and a matching
SHA-256 file. A report should identify unsafe archive paths, checksum bypasses,
unexpected package files, or release-signature problems.

## Reporting a Vulnerability

Don't open a public issue for an undisclosed vulnerability. Use the repository's
GitHub security advisory mechanism when available. Otherwise, contact the
maintainer through the private contact method on the GitHub profile.

Include the affected version or commit, reproduction steps, impact, and a minimal
proof of concept when safe. Don't include credentials, personal data,
conversation records, or other sensitive material.
