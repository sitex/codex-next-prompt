# Codex Next Prompt

Codex Next Prompt 0.2.0 is a standalone user skill invoked explicitly as
`$next`. It generates ready-to-send prompt text for the best next step in the
current conversation and never executes that prompt.

## Behavior

The default result is one concrete prompt. The skill returns two or three
separate prompts only when the conversation leaves a real, materially different
fork unresolved. If required context is missing, it asks for that context rather
than inventing values or emitting placeholders. Output follows the conversation's
language.

```text
$next
```

Possible output:

```text
Reproduce the focused test failure, identify the root cause, apply the smallest fix, and rerun the focused test before the full suite.
```

See [EXAMPLES.md](EXAMPLES.md) for more scenarios.

## Install 0.2.0

### Recommended: `$skill-installer`

In Codex, ask `$skill-installer` to install the skill from this repository path:

```text
https://github.com/sitex/codex-next-prompt/tree/main/skills/next
```

Start a new Codex session after installation.

### Manual: release ZIP

Download both assets from the
[v0.2.0 release](https://github.com/sitex/codex-next-prompt/releases/tag/v0.2.0):

```text
codex-next-prompt-0.2.0.zip
codex-next-prompt-0.2.0.zip.sha256
```

Verify and inspect them before extraction:

```sh
sha256sum -c codex-next-prompt-0.2.0.zip.sha256
python3 -m zipfile -l codex-next-prompt-0.2.0.zip
python3 -m zipfile -e codex-next-prompt-0.2.0.zip .
```

The ZIP contains exactly:

```text
codex-next-prompt-0.2.0/next/SKILL.md
codex-next-prompt-0.2.0/next/agents/openai.yaml
```

Copy the extracted `next` directory to the user skill directory:

```sh
mkdir -p "${CODEX_HOME:-$HOME/.codex}/skills"
cp -R codex-next-prompt-0.2.0/next "${CODEX_HOME:-$HOME/.codex}/skills/next"
```

Do not nest it as `skills/next/next`. Start a new session after copying it.

## Migrate from 0.1

Version 0.1 used an experimental plugin and local marketplace. Remove that old
plugin installation and marketplace registration before installing the 0.2.0
standalone skill. Also remove any old `next` directory from
`${CODEX_HOME:-$HOME/.codex}/skills` before copying the replacement.

## Verify and Uninstall

Verify in a new session:

```sh
codex debug prompt-input '$next'
```

You can also send `$next` in that session and confirm that the response is prompt
text only. There is no custom slash command.

Uninstall the standalone skill by removing only:

```text
$CODEX_HOME/skills/next
```

When `CODEX_HOME` is unset, the default is `$HOME/.codex`.

## Boundaries

Version 0.2.0 has no hooks, executable runtime, background process, automatic
invocation, network request, telemetry, persistence, transcript access, model API
call, or nested Codex session. It uses only context already available in the
current conversation.

## Development and Release

```sh
make test
make package
make check
make clean
```

The top-level `VERSION` file is authoritative. Packaging produces one
deterministic ZIP with exactly the two standalone skill files and one SHA-256
file. Release tags use `vMAJOR.MINOR.PATCH`, are annotated and signed, and must
match `VERSION`. The workflow verifies the approved public signing key before
publishing.

See [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md), and
[INSTALL_WITH_LLM.md](INSTALL_WITH_LLM.md).

## License

Copyright (c) 2026 Rocky. Released under the [MIT License](LICENSE).
