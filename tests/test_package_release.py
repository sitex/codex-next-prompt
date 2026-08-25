import hashlib
import os
import subprocess
import tempfile
import unittest
import zipfile
from pathlib import Path
from typing import Final


ROOT: Final = Path(__file__).resolve().parents[1]
PACKAGER: Final = ROOT / "scripts" / "package_release.py"
VERSION: Final = "1.2.3"
ARCHIVE_NAME: Final = f"codex-next-prompt-{VERSION}.zip"
PAYLOAD: Final = (
    "next/SKILL.md",
    "next/agents/openai.yaml",
)


def create_source(root: Path, version: str = VERSION) -> None:
    files = {
        "VERSION": f"{version}\n",
        "skills/next/SKILL.md": "---\nname: next\n---\n\n# Next\n",
        "skills/next/agents/openai.yaml": "policy:\n  allow_implicit_invocation: false\n",
    }
    for relative_path, content in files.items():
        path = root / relative_path
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")


def run_packager(source: Path, output: Path, version: str = VERSION) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        ["python3", str(PACKAGER), version],
        cwd=source,
        env={**os.environ, "DIST_DIR": str(output)},
        capture_output=True,
        text=True,
        check=False,
    )


class TestPackageRelease(unittest.TestCase):
    def test_package_contains_exact_deterministic_portable_payload(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            workspace = Path(temporary)
            source = workspace / "source"
            first_output = workspace / "first"
            second_output = workspace / "second"
            create_source(source)

            first = run_packager(source, first_output)
            second = run_packager(source, second_output)

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            first_archive = first_output / ARCHIVE_NAME
            second_archive = second_output / ARCHIVE_NAME
            self.assertEqual(first_archive.read_bytes(), second_archive.read_bytes())

            archive_root = f"codex-next-prompt-{VERSION}"
            expected_names = tuple(f"{archive_root}/{path}" for path in PAYLOAD)
            with zipfile.ZipFile(first_archive) as archive:
                self.assertEqual(tuple(archive.namelist()), expected_names)
                self.assertIsNone(archive.testzip())
                for info in archive.infolist():
                    with self.subTest(path=info.filename):
                        self.assertEqual(info.date_time, (1980, 1, 1, 0, 0, 0))
                        self.assertEqual(info.create_system, 3)
                        self.assertEqual(info.compress_type, zipfile.ZIP_STORED)
                        self.assertEqual(info.external_attr >> 16, 0o100644)

            digest = hashlib.sha256(first_archive.read_bytes()).hexdigest()
            checksum = (first_output / f"{ARCHIVE_NAME}.sha256").read_text(encoding="ascii")
            self.assertEqual(checksum, f"{digest}  {ARCHIVE_NAME}\n")

    def test_package_rejects_invalid_version_and_source_version_mismatch(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            workspace = Path(temporary)
            source = workspace / "source"
            create_source(source)

            for version in ("v1.2.3", "01.2.3", "1.2.3-01", "1.2", "1.2.4"):
                with self.subTest(version=version):
                    result = run_packager(source, workspace / version.replace("/", "_"), version)
                    self.assertNotEqual(result.returncode, 0)

    def test_package_rejects_missing_empty_and_symlink_source_version(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            workspace = Path(temporary)

            for condition in ("missing", "empty", "symlink"):
                with self.subTest(condition=condition):
                    source = workspace / condition
                    create_source(source)
                    target = source / "VERSION"
                    if condition == "missing":
                        target.unlink()
                    elif condition == "empty":
                        target.write_bytes(b"")
                    else:
                        replacement = source / "replacement-version"
                        replacement.write_text(f"{VERSION}\n", encoding="ascii")
                        target.unlink()
                        target.symlink_to(replacement)

                    result = run_packager(source, workspace / f"version-{condition}")
                    self.assertNotEqual(result.returncode, 0)
                    self.assertFalse((workspace / f"version-{condition}" / ARCHIVE_NAME).exists())

    def test_package_rejects_missing_empty_and_symlink_payload(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            workspace = Path(temporary)

            for condition in ("missing", "empty", "symlink"):
                with self.subTest(condition=condition):
                    source = workspace / condition
                    create_source(source)
                    target = source / "skills" / "next" / "SKILL.md"
                    if condition == "missing":
                        target.unlink()
                    elif condition == "empty":
                        target.write_bytes(b"")
                    else:
                        replacement = source / "replacement.md"
                        replacement.write_text("replacement\n", encoding="utf-8")
                        target.unlink()
                        target.symlink_to(replacement)

                    result = run_packager(source, workspace / f"dist-{condition}")
                    self.assertNotEqual(result.returncode, 0)
                    self.assertFalse((workspace / f"dist-{condition}" / ARCHIVE_NAME).exists())

    def test_package_rejects_symlinked_payload_parent(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            workspace = Path(temporary)
            source = workspace / "source"
            create_source(source)
            real_skills = workspace / "real-skills"
            (source / "skills").rename(real_skills)
            (source / "skills").symlink_to(real_skills, target_is_directory=True)

            result = run_packager(source, workspace / "dist")

            self.assertNotEqual(result.returncode, 0)
            self.assertFalse((workspace / "dist" / ARCHIVE_NAME).exists())


if __name__ == "__main__":
    unittest.main()
