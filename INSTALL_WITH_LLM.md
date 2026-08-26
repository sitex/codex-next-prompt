# Install the 0.2.0 Standalone Skill with an LLM

This procedure installs Codex Next Prompt as a standalone user skill. Its only
invocation is explicit `$next`.

## Release Installation

1. Select the release published from the signed v0.2.0 tag and download exactly
   `codex-next-prompt-0.2.0.zip` and its matching `.sha256` file.
2. Run `sha256sum -c codex-next-prompt-0.2.0.zip.sha256`.
3. Inspect the ZIP and reject absolute paths, parent traversal, symbolic links,
   empty files, or any payload other than:

```text
codex-next-prompt-0.2.0/next/SKILL.md
```

4. Extract the ZIP only after the checks pass.
5. Resolve the destination as `${CODEX_HOME:-$HOME/.codex}/skills/next`.
6. Refuse to overwrite an existing destination without first identifying whether
   it is an old installation and obtaining approval to replace that exact path.
7. Copy the extracted `next` directory to the destination. The final file must
   be `${CODEX_HOME:-$HOME/.codex}/skills/next/SKILL.md`.
8. Start a new Codex session.

Do not install from a GitHub source tree or mutable branch, build from source,
request credentials, read conversation records, weaken a failed safety check,
or change unrelated Codex settings.

## Migration from 0.1

The 0.1 release used an experimental plugin and local marketplace. Remove both
old registrations before installing the 0.2.0 standalone skill. Delete an old
extracted release or skill directory only after showing and validating its exact
resolved path. Never delete a parent or wildcard path.

## Codex 0.149.1 Compatibility

Install only `SKILL.md`. Runtime testing showed that agent policy metadata which
disables implicit invocation hides the skill from the explicit catalog in Codex
0.149.1. The supported trigger is description-driven: the user types `$next`, and
the frontmatter description declares that explicit intent. There is no automatic
footer.

## Verification

Use either:

```sh
codex debug prompt-input '$next'
```

or start a new session and send:

```text
$next
```

The result must be ready-to-send prompt text based only on the current
conversation. The skill must not execute it. There is no custom slash command,
automatic invocation, or hook.

## Uninstall

Remove only `${CODEX_HOME:-$HOME/.codex}/skills/next`. Report the selected
version, checksum result, inspected payload,
destination path, restart status, and `$next` verification result.
