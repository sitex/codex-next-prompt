#!/usr/bin/env python3
# /// script
# requires-python = ">=3.11"
# dependencies = []
# ///
# How to run: python3 scripts/package_release.py VERSION

import hashlib
import os
import re
import stat
import sys
import tempfile
import zipfile
from dataclasses import dataclass
from pathlib import Path
from typing import Final, NewType


Version = NewType("Version", str)
SEMVER: Final = re.compile(
    r"^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)"
    r"(?:-(?:(?:0|[1-9][0-9]*)|[0-9]*[A-Za-z-][0-9A-Za-z-]*)"
    r"(?:\.(?:(?:0|[1-9][0-9]*)|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*)?"
    r"(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$"
)
PAYLOAD: Final = (
    (Path("skills/next/SKILL.md"), Path("next/SKILL.md")),
    (Path("skills/next/agents/openai.yaml"), Path("next/agents/openai.yaml")),
)
ZIP_TIMESTAMP: Final = (1980, 1, 1, 0, 0, 0)
FILE_MODE: Final = stat.S_IFREG | 0o644


@dataclass(frozen=True, slots=True)
class PackageError(Exception):
    message: str

    def __str__(self) -> str:
        return self.message


def parse_version(raw_version: str) -> Version:
    if SEMVER.fullmatch(raw_version) is None:
        raise PackageError(f"invalid semantic version: {raw_version}")
    return Version(raw_version)


def read_source_version(root: Path) -> Version:
    version_path = root / "VERSION"
    ensure_payload_file(root, version_path)
    try:
        raw_version = version_path.read_text(encoding="ascii")
    except UnicodeDecodeError as error:
        raise PackageError("VERSION must contain an ASCII semantic version") from error
    if raw_version != raw_version.strip() + "\n":
        raise PackageError("VERSION must contain one semantic version followed by a newline")
    return parse_version(raw_version.strip())


def ensure_payload_file(root: Path, path: Path) -> None:
    relative_path = path.relative_to(root)
    current = root
    for component in relative_path.parts:
        current /= component
        if current.is_symlink():
            raise PackageError(f"payload path must not be a symlink: {relative_path}")
    if not path.is_file():
        raise PackageError(f"missing payload file: {relative_path}")
    if path.stat().st_size == 0:
        raise PackageError(f"payload file must not be empty: {relative_path}")


def package_release(root: Path, output_dir: Path, version: Version) -> tuple[Path, Path]:
    source_version = read_source_version(root)
    if source_version != version:
        raise PackageError(
            f"version argument {version} does not match VERSION {source_version}"
        )

    files = tuple((source_path, archive_path, root / source_path) for source_path, archive_path in PAYLOAD)
    for _source_path, _archive_path, source_path in files:
        ensure_payload_file(root, source_path)

    output_dir.mkdir(parents=True, exist_ok=True)
    release_name = f"codex-next-prompt-{version}"
    archive_path = output_dir / f"{release_name}.zip"
    checksum_path = output_dir / f"{archive_path.name}.sha256"

    with tempfile.NamedTemporaryFile(dir=output_dir, delete=False) as temporary:
        temporary_archive = Path(temporary.name)
    try:
        with zipfile.ZipFile(temporary_archive, "w") as archive:
            for _source_relative, archive_relative, source_path in files:
                info = zipfile.ZipInfo(f"{release_name}/{archive_relative.as_posix()}", ZIP_TIMESTAMP)
                info.compress_type = zipfile.ZIP_STORED
                info.create_system = 3
                info.external_attr = FILE_MODE << 16
                archive.writestr(info, source_path.read_bytes())
        os.replace(temporary_archive, archive_path)
    finally:
        temporary_archive.unlink(missing_ok=True)

    digest = hashlib.sha256(archive_path.read_bytes()).hexdigest()
    checksum_path.write_text(f"{digest}  {archive_path.name}\n", encoding="ascii")
    return archive_path, checksum_path


def main(arguments: list[str]) -> int:
    if len(arguments) != 1:
        print("usage: scripts/package_release.py VERSION", file=sys.stderr)
        return 2
    try:
        version = parse_version(arguments[0])
        output_dir = Path(os.environ.get("DIST_DIR", "dist"))
        archive_path, checksum_path = package_release(Path.cwd(), output_dir, version)
    except (PackageError, OSError) as error:
        print(f"package_release: {error}", file=sys.stderr)
        return 1
    print(archive_path)
    print(checksum_path)
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
