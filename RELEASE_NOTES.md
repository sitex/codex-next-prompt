# Codex Next Prompt 0.2.0

Version 0.2.0 distributes `$next` as a standalone skill for the current user. It generates
ready-to-send prompt text from the current conversation and never executes its
result.

## Install

Recommended: use `$skill-installer` with:

```text
https://github.com/sitex/codex-next-prompt/tree/main/skills/next
```

For a manual release install, verify the ZIP checksum, extract it, and copy the
extracted `next` directory to `$CODEX_HOME/skills/next` (default
`$HOME/.codex/skills/next`). The archive contains exactly:

```text
codex-next-prompt-0.2.0/next/SKILL.md
codex-next-prompt-0.2.0/next/agents/openai.yaml
```

Remove the experimental 0.1 plugin and local marketplace registration before the
standalone install. Start a new session and verify with
`codex debug prompt-input '$next'` or by sending `$next`.

## Contract

- Explicit `$next` invocation only.
- One default prompt; two or three only for a real fork.
- Current conversation context only.
- No hooks, automatic invocation, executable runtime, network access,
  persistence, transcript access, model API call, or prompt execution.

The release remains one deterministic ZIP plus one SHA-256 file. Its annotated
tag is signed and verified against the approved public key before publication.
