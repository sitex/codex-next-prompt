# Codex Next Prompt 0.2.0

Version 0.2.0 distributes `$next` as a standalone skill for the current user. It generates
ready-to-send prompt text from the current conversation and never executes its
result.

## Install

Install only the release artifacts published from the signed v0.2.0 tag:

```text
codex-next-prompt-0.2.0.zip
codex-next-prompt-0.2.0.zip.sha256
```

Verify the checksum, inspect and extract the ZIP, and copy the extracted `next`
directory to `${CODEX_HOME:-$HOME/.codex}/skills/next`. Do not install from a
GitHub source tree or mutable branch. The archive contains exactly:

```text
codex-next-prompt-0.2.0/next/SKILL.md
```

Remove the experimental 0.1 plugin and local marketplace registration before the
standalone install. Start a new session and verify with
`codex debug prompt-input '$next'` or by sending `$next`.

## Contract

- Explicit `$next` invocation only.
- Description-driven frontmatter trigger when the user types `$next`.
- One default prompt; two or three only for a real fork.
- Current conversation context only.
- No hooks, automatic invocation, executable runtime, network access,
  persistence, transcript access, model API call, or prompt execution.

Codex 0.149.1 must receive only `SKILL.md`: agent policy metadata that disables
implicit invocation hides the skill from the explicit catalog in this version.
The skill never adds an automatic footer.

The release remains one deterministic ZIP plus one SHA-256 file. Its annotated
tag is signed and verified against the approved public key before publication.
