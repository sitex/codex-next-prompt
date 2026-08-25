# Install Codex Next Prompt with an LLM

This document is written for an LLM or coding agent installing
`sitex/codex-next-prompt` on behalf of a user.

## Safety Contract

- Install only from a published GitHub Release. Do not build or install from the
  source checkout.
- Download the archive and its matching `.sha256` file from the same release.
- Stop immediately if the checksum does not match.
- Do not use `--dangerously-bypass-hook-trust`.
- Do not approve or simulate `/hooks` trust for the user. The user must review
  and trust the hooks interactively.
- Keep the extracted release directory in place while its marketplace and
  plugin are installed.
- Do not request credentials, API keys, or network access beyond downloading
  the public release and running the normal Codex model connection.

## Installation Procedure

### 1. Check prerequisites

Run:

```text
codex --version
```

Require Codex `0.149.1` or newer. If `codex` is unavailable or too old, explain
the requirement and stop.

### 2. Detect the platform

Select exactly one release target:

| Operating system | Architecture | Target |
|---|---|---|
| Linux | x86_64 / amd64 | `linux-amd64` |
| Linux | arm64 / aarch64 | `linux-arm64` |
| macOS | x86_64 / amd64 | `darwin-amd64` |
| macOS | arm64 / Apple Silicon | `darwin-arm64` |
| Windows | AMD64 / x86_64 | `windows-amd64` |
| Windows | ARM64 | `windows-arm64` |

If the operating system or architecture is not listed, report that it is
unsupported and stop.

### 3. Resolve a published release

Query the latest GitHub Release for
`https://github.com/sitex/codex-next-prompt` and obtain its tag/version. If no
release exists, tell the user that installation cannot continue until a release
is published. Do not fall back to the source repository.

For version `VERSION` and target `TARGET`, the expected filenames are:

- Linux/macOS: `codex-next-prompt-VERSION-TARGET.tar.gz`
- Windows: `codex-next-prompt-VERSION-TARGET.zip`
- Checksum: the archive filename plus `.sha256`

### 4. Download and verify

Create a stable user-selected installation directory, not a temporary directory
that will be deleted automatically. Download the archive and checksum there.

On Linux:

```sh
sha256sum -c "ARCHIVE.sha256"
```

On macOS:

```sh
shasum -a 256 -c "ARCHIVE.sha256"
```

On Windows PowerShell, parse the expected hash from the checksum file and compare
it with:

```powershell
(Get-FileHash -Algorithm SHA256 ARCHIVE).Hash
```

Comparison on Windows must be case-insensitive. If verification fails, delete
the downloaded files, report the failure, and stop.

### 5. Extract the archive

Extract the archive into the stable installation directory. The extracted root
is named `codex-next-prompt-VERSION` and contains a self-contained local
marketplace.

Do not move individual files out of this directory.

### 6. Register and install

Run these commands with the extracted root path:

```text
codex plugin marketplace add "PATH_TO_EXTRACTED_ROOT"
codex plugin marketplace list
codex plugin list --available
codex plugin add codex-next-prompt@codex-next-prompt
codex plugin list
```

Verify that `codex-next-prompt@codex-next-prompt` is installed and enabled. If
any command fails, report the exact command and error. Do not bypass Codex
security controls.

### 7. Hand control back to the user

Tell the user to:

1. Start a new interactive Codex session.
2. Open `/hooks`.
3. Review the `SessionStart` and `Stop` commands and confirm they point inside
   the checksum-verified release directory.
4. Trust the hooks only if the paths and commands are correct.

Do not claim installation is fully active until the user completes this review.

## Uninstall

Run:

```text
codex plugin remove codex-next-prompt@codex-next-prompt
codex plugin marketplace remove codex-next-prompt
```

After both commands succeed, the extracted release directory and downloaded
archive may be removed.

## Completion Report

Report:

- detected Codex version;
- selected release version and target;
- checksum result;
- extracted marketplace path;
- plugin list result;
- that interactive `/hooks` review is still required.

Never include authentication files, tokens, session transcripts, or unrelated
Codex configuration in the report.
