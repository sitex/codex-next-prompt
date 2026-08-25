# Codex Next Prompt 0.2.0

Version 0.2.0 replaces the automatic response behavior with one explicit,
skill-only action. Invoke `$next` to generate ready-to-send prompt text from the
current conversation. The skill never executes its result.

## Breaking Migration

Remove the previous plugin and marketplace registration before installing
0.2.0. Delete the old extracted release directory only after checking its exact
path, then install the new package from a clean directory. Existing automatic
response behavior doesn't carry forward.

## Package

Every release has one portable asset and one checksum:

```text
codex-next-prompt-0.2.0.zip
codex-next-prompt-0.2.0.zip.sha256
```

Verify SHA-256, inspect the ZIP paths, extract the single release root, register
that root as a local marketplace, install the plugin, and restart Codex. See
[INSTALL_WITH_LLM.md](INSTALL_WITH_LLM.md) for the safe procedure.

## Contract

- `$next` is the canonical and explicit invocation.
- Input comes only from the current conversation.
- Output is one default prompt, or two to three prompts for a real fork.
- Missing context produces a request for context, not invented values.
- Output follows the conversation's language.
- The package contains instructions and metadata only.
- There is no automatic invocation, executable runtime, network access,
  persistence, conversation-record access, model API call, or prompt execution.

Release tags remain annotated, signed, and checked against the approved public
key before the single ZIP is published.
