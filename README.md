# Codex Next Prompt

Codex Next Prompt 0.2.0 adds one explicit Codex skill: `$next`.

Invoke `$next` when you want a ready-to-send prompt for the best next step in
the current conversation. The skill reads only the context already available in
that conversation. It returns prompt text, not commentary, and never executes
the prompt.

## What `$next` Returns

The default is one clear prompt for the most useful next action.

When the conversation contains a real decision with materially different paths,
the skill may return two or three separate prompts, one per path. Cosmetic
variations stay in one prompt. If required context is missing, `$next` asks for
that context instead of inventing values or leaving placeholders.

The output follows the conversation's language. It can recommend inspection,
planning, implementation, testing, review, or another concrete action, but it
doesn't run commands, edit files, send messages, or continue the task itself.

## Example

After discussing a failing test, invoke:

```text
$next
```

Possible output:

```text
Reproduce the focused test failure, identify the root cause, apply the smallest fix, and rerun the focused test before the full suite.
```

See [EXAMPLES.md](EXAMPLES.md) for forks, missing context, Russian output, and
the no-execution boundary.

## Install 0.2.0

Download these two assets from the
[v0.2.0 release](https://github.com/sitex/codex-next-prompt/releases/tag/v0.2.0):

```text
codex-next-prompt-0.2.0.zip
codex-next-prompt-0.2.0.zip.sha256
```

Verify the checksum from the directory containing both files:

```sh
sha256sum -c codex-next-prompt-0.2.0.zip.sha256
```

Inspect the ZIP before extraction:

```sh
python3 -m zipfile -l codex-next-prompt-0.2.0.zip
```

The archive must contain only the `codex-next-prompt-0.2.0/` directory and its
plugin metadata and skill files. Extract it, then register the extracted root as
a local marketplace:

```sh
python3 -m zipfile -e codex-next-prompt-0.2.0.zip .
codex plugin marketplace add "./codex-next-prompt-0.2.0"
codex plugin marketplace list
codex plugin list --available
codex plugin add codex-next-prompt@codex-next-prompt
codex plugin list
```

Restart Codex after installation. Confirm that the plugin is enabled and the
`next` skill is discoverable. Enabled skills may also appear in Codex's skill or
slash listing, but the canonical invocation is `$next`.

For an agent-run installation with path checks before extraction, use
[INSTALL_WITH_LLM.md](INSTALL_WITH_LLM.md).

### Upgrade from the previous release

The pre-0.2 package used a different runtime design. Remove its installed plugin
and marketplace first, delete its extracted release directory only after
checking the path, then install 0.2.0 from the ZIP above. Don't place 0.2.0 over
an older extracted directory.

```sh
codex plugin remove codex-next-prompt@codex-next-prompt
codex plugin marketplace remove codex-next-prompt
```

### Uninstall

```sh
codex plugin remove codex-next-prompt@codex-next-prompt
codex plugin marketplace remove codex-next-prompt
```

After both commands succeed, remove the extracted
`codex-next-prompt-0.2.0/` directory if you no longer need it.

## Boundaries

Version 0.2.0 is instruction-only and skill-only. It has no executable runtime,
background process, automatic footer, automatic invocation, composer mutation,
network request, telemetry, persistence, conversation-record access, model API
call, or nested Codex session. It doesn't inspect files unless their contents
are already present in the current conversation.

## Development and Release

The repository uses Python's standard library for tests and deterministic ZIP
packaging:

```sh
make test
make package
make check
make clean
```

Use TDD for contract changes: add a focused failing test, make the smallest
change that passes, then run the full repository checks. Release packaging must
produce exactly one portable ZIP and its SHA-256 file.

Release tags use `vMAJOR.MINOR.PATCH`, are annotated and signed, and must match
the version in `.codex-plugin/plugin.json`. The release workflow verifies the
tag against [`RELEASE_SIGNING_KEY.asc`](RELEASE_SIGNING_KEY.asc) and fingerprint
`BFA7D43C126EE54A5FC8DD0EBE645A3EFA752D77`. Never commit a private signing key.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the skill-only development contract.

## Security

See [SECURITY.md](SECURITY.md) for the instruction-only threat surface and
private reporting guidance.

## License

Copyright (c) 2026 Rocky. Released under the [MIT License](LICENSE).
