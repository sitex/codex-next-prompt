# Codex Next Prompt

Codex Next Prompt is a local-only Codex plugin that asks the active model to
append one concise follow-up suggestion to suitable final responses.

The supported response-level UX is:

```text
Suggested next prompt: <a specific prompt the user could submit>
```

The suggestion is advisory. It may be omitted when there is no meaningful next
action, when the response is an exact-output task, for a terse acknowledgement
or safety refusal, or while a question is waiting for user-provided data.

Exact-output suppression has been verified against a response required to be
exactly `OK`: the plugin left the output unchanged.

## Usage Example

Start a new Codex session and send:

```text
Explain why focused tests are useful after changing code and recommend what to do next.
```

The response should finish with one relevant line similar to:

```text
Suggested next prompt: Run the focused tests for my changed module and explain any failures.
```

The wording is model-generated and can vary. See [EXAMPLES.md](EXAMPLES.md) for
exact-output suppression, waiting-for-input behavior, and troubleshooting.

## Install

Codex `0.149.1` or newer is required. The source repository does not contain
generated binaries and is not a directly installable marketplace. Download the
release archive for your platform and its checksum from the matching
[GitHub release](https://github.com/sitex/codex-next-prompt/releases), then verify
and extract both files from the same directory. For example, on Linux `amd64`:

```sh
VERSION=0.1.0
curl -fLO "https://github.com/sitex/codex-next-prompt/releases/download/v${VERSION}/codex-next-prompt-${VERSION}-linux-amd64.tar.gz"
curl -fLO "https://github.com/sitex/codex-next-prompt/releases/download/v${VERSION}/codex-next-prompt-${VERSION}-linux-amd64.tar.gz.sha256"
sha256sum -c "codex-next-prompt-${VERSION}-linux-amd64.tar.gz.sha256"
tar -xzf "codex-next-prompt-${VERSION}-linux-amd64.tar.gz"
```

Use the extracted release root as the local marketplace. Its packaged catalog
resolves `./` to that same root, which contains the target binary:

```sh
codex plugin marketplace add "./codex-next-prompt-${VERSION}"
codex plugin marketplace list
codex plugin list --available
codex plugin add codex-next-prompt@codex-next-prompt
```

### Install with an LLM

Give your coding agent this prompt:

```text
Install Codex Next Prompt by following the instructions at:
https://raw.githubusercontent.com/sitex/codex-next-prompt/main/INSTALL_WITH_LLM.md

Follow the checksum and hook-trust safety requirements exactly. Do not install
from source, do not bypass hook trust, and stop if no published release exists.
```

The machine-oriented procedure is also available as
[INSTALL_WITH_LLM.md](INSTALL_WITH_LLM.md). The LLM can download, verify, and
register the plugin, but the user must review and trust the hooks interactively
through `/hooks`.

Keep the extracted directory in place while the marketplace and plugin are
installed. Start a new Codex session, open `/hooks`, review the `SessionStart`
and `Stop` commands and their local paths, and trust them only after confirming
they point into the verified release directory. The plugin can then add
`Suggested next prompt:` to suitable final responses.

Codex `0.149.1` has no separate disable command. Remove the plugin to stop its
hooks, and remove the static local marketplace when it is no longer needed:

```text
codex plugin remove codex-next-prompt@codex-next-prompt
codex plugin marketplace remove codex-next-prompt
rm -rf "./codex-next-prompt-${VERSION}"
```

## Scope

The plugin uses supported Codex lifecycle hooks. Session-start guidance asks the
model for the response-level line, and a stop hook validates the completed
response without blocking it. The minimum supported Codex version is `0.149.1`.

This project does **not** provide composer insertion, ghost text, draft
mutation, prefill, Tab/Right acceptance, or any other native composer API.
There is no network access, telemetry, transcript reading, session-history
persistence, analytics, external model call, nested Codex invocation, or
automatic follow-up execution.

The plugin has no configuration and doesn't change the composer. It only guides
response generation and checks the completed response through lifecycle hooks.

## Platforms

Release packaging covers Linux and macOS on `amd64` and `arm64`, plus Windows
on `amd64` and `arm64`. The full install, trust, response, and stop lifecycle has
been tested with Codex `0.149.1` on Linux. Windows smoke tests run in CI. Native
macOS execution hasn't been tested locally.

## Privacy

The plugin runs locally and makes no network requests. It has no authentication
flow, telemetry, analytics, persistence, or external model call. It doesn't read
transcripts or session history. Codex still uses its normal model connection to
produce responses.

## Development

Contributors need Go 1.25. The runtime uses only the Go standard library.

```sh
make fmt
make vet
make test
make build
make smoke
make package
make check
```

To test packaging directly or build a target release archive:

```sh
go test ./tests -run 'Test_(MarketplaceCatalog|PluginManifest|HookManifest|HookLaunchers)' -count=1
make package VERSION=0.1.0 GOOS=linux GOARCH=amd64
```

Release archives and checksums are produced by `scripts/package-release.sh`.
Each archive is a self-contained local marketplace root; the repository copy of
`.agents/plugins/marketplace.json` defines that packaged distribution catalog,
not a source-checkout installation path.
GitHub release automation builds the supported platform matrix from a version
tag. Release tags must be annotated and signed because CI runs `git verify-tag`
and confirms the tag resolves to the workflow commit before packaging. CI
imports the repository's public [`RELEASE_SIGNING_KEY.asc`](RELEASE_SIGNING_KEY.asc)
into an isolated keyring, requires it to contain exactly one primary key, and
requires the tag signature's primary fingerprint to equal
`BFA7D43C126EE54A5FC8DD0EBE645A3EFA752D77`.
Rotating this public key or fingerprint requires a reviewed code change; private
signing keys must never be committed. The public key is source-only and is not
included in plugin runtime release archives.
Maintainers should run `make check` before creating a release tag.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development and review rules.

## Security

See [SECURITY.md](SECURITY.md) for responsible disclosure instructions.

## License

Copyright (c) 2026 Rocky. Released under the [MIT License](LICENSE).
