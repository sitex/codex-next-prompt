## Summary

Describe the user-visible or contributor-visible change for version 0.2.0.
The product is a standalone skill invoked as `$next`.

## `$next` Contract

- [ ] Keeps `$next` explicit and returns ready-to-send prompt text only.
- [ ] Uses only current conversation context and never executes the result.
- [ ] Returns one default prompt, with separate prompts only for a real fork.
- [ ] Requests missing context instead of inventing values or placeholders.
- [ ] Adds no executable runtime, automatic invocation, network access,
      persistence, conversation-record reads, model API calls, or secret access.

## TDD and Safety

- [ ] Added or updated a focused failing test before the implementation change.
- [ ] Tests assert routing or repository contract terms, not prose snapshots.
- [ ] Package changes keep one portable ZIP containing exactly `next/SKILL.md`,
      with checksum verification and safe paths.
- [ ] Codex 0.149.1 discovery remains description-driven when the user types
      `$next`, with no agent policy metadata or automatic footer.
- [ ] Installation remains available after release through `$skill-installer`
      from `tree/v0.2.0/skills/next` or manually at
      `${CODEX_HOME:-$HOME/.codex}/skills/next`.
- [ ] Contains no credentials, private keys, user data, generated ZIP files, or
      unrelated changes.

## Validation

Commands run:

```text

```

## Related issue

Closes #12
