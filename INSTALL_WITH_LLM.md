# Install Codex Next Prompt 0.2.0 with an LLM

This procedure installs the published, skill-only Codex Next Prompt package.
After installation, the canonical invocation is `$next`.

## Safety Contract

- Install version 0.2.0 or newer only from a published GitHub release.
- Download exactly one ZIP and its matching `.sha256` file from the same release.
- Verify the checksum before inspecting or extracting the ZIP.
- Reject absolute paths, parent traversal, symbolic links, and files outside the
  expected top-level release directory.
- Don't build from source, request credentials, read conversation records, or
  change unrelated Codex settings.
- Remove any pre-0.2 installation before registering the new marketplace.
- Report failures and stop. Don't weaken a failed safety check.

## Procedure

### 1. Confirm the release

Resolve a published release at
`https://github.com/sitex/codex-next-prompt/releases`. Require a semantic version
of 0.2.0 or newer. Let `VERSION` be that version.

Download both files into a stable user-selected directory:

```text
codex-next-prompt-VERSION.zip
codex-next-prompt-VERSION.zip.sha256
```

Don't use a temporary directory that will disappear while the marketplace is
registered.

### 2. Verify the checksum

From the download directory, run:

```sh
sha256sum -c "codex-next-prompt-${VERSION}.zip.sha256"
```

If `sha256sum` isn't available, compute SHA-256 with a trusted local tool and
compare the full hexadecimal digest with the checksum file. Stop on any mismatch.

### 3. Inspect the ZIP safely

Set `ARCHIVE` and run this standard-library inspection before extraction:

```sh
ARCHIVE="codex-next-prompt-${VERSION}.zip" python3 - <<'PY'
import os
import stat
import zipfile
from pathlib import PurePosixPath

archive_path = os.environ["ARCHIVE"]
version = archive_path.removeprefix("codex-next-prompt-").removesuffix(".zip")
root = f"codex-next-prompt-{version}"

with zipfile.ZipFile(archive_path) as archive:
    for member in archive.infolist():
        path = PurePosixPath(member.filename)
        mode = member.external_attr >> 16
        if path.is_absolute() or ".." in path.parts:
            raise SystemExit(f"unsafe ZIP path: {member.filename}")
        if not path.parts or path.parts[0] != root:
            raise SystemExit(f"unexpected ZIP root: {member.filename}")
        if stat.S_ISLNK(mode):
            raise SystemExit(f"symbolic link not allowed: {member.filename}")
        print(member.filename)
PY
```

The listing should contain only plugin metadata and the `next` skill beneath the
single expected root. Stop if any entry is unexpected.

### 4. Remove an older installation

Check whether `codex-next-prompt@codex-next-prompt` or its marketplace is already
installed. If so, remove them before extraction and registration:

```sh
codex plugin remove codex-next-prompt@codex-next-prompt
codex plugin marketplace remove codex-next-prompt
```

If an older extracted directory must be deleted, show the exact resolved path to
the user first. Delete only that release directory, never a parent or wildcard
path.

### 5. Extract and install

Extract only after checksum and path inspection succeed:

```sh
python3 -m zipfile -e "codex-next-prompt-${VERSION}.zip" .
codex plugin marketplace add "./codex-next-prompt-${VERSION}"
codex plugin marketplace list
codex plugin list --available
codex plugin add codex-next-prompt@codex-next-prompt
codex plugin list
```

Verify that the plugin is installed and enabled. Keep the extracted directory in
place while the local marketplace is registered.

### 6. Verify discovery

Restart Codex. Confirm that the `next` skill is discoverable in the enabled
skills. Then ask the user to invoke:

```text
$next
```

The skill should return ready-to-send prompt text based on the current
conversation. It must not run that prompt.

## Uninstall

```sh
codex plugin remove codex-next-prompt@codex-next-prompt
codex plugin marketplace remove codex-next-prompt
```

After both commands succeed, remove the exact extracted release directory and
downloaded ZIP files if the user wants them deleted.

## Completion Report

Report the selected version, checksum result, inspected top-level directory,
installation path, plugin discovery result, restart status, and `$next`
invocation result. Don't include credentials, conversation records, private
paths unrelated to this install, or unrelated configuration.
