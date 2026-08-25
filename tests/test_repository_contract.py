import json
import re
import subprocess
import unittest
from pathlib import Path
from typing import Final
from urllib.parse import unquote


ROOT: Final = Path(__file__).resolve().parents[1]
MARKDOWN_LINK: Final = re.compile(r"!?\[[^]]*\]\(([^)]+)\)")
REMOVED_PATHS: Final = (
    "cmd/",
    "hooks/",
    "internal/",
    "testdata/",
)
REMOVED_FILES: Final = {
    "go.mod",
    "scripts/package-release.sh",
    "tests/archive_smoke_test.go",
    "tests/build_contract_test.go",
    "tests/marketplace_test.go",
    "tests/package-release_test.go",
    "tests/packaging_test.go",
    "tests/release_signature_test.go",
    "tests/release_workflow_test.go",
    "tests/repository_metadata_test.go",
    "tests/smoke-archive-posix.sh",
    "tests/smoke-posix.sh",
    "tests/smoke-windows.ps1",
}


def tracked_files() -> tuple[str, ...]:
    result = subprocess.run(
        ["git", "ls-files", "--cached", "--others", "--exclude-standard"],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )
    return tuple(
        path
        for path in result.stdout.splitlines()
        if path and (ROOT / path).is_file()
    )


class TestRepositoryContract(unittest.TestCase):
    def test_repository_contains_no_legacy_runtime_when_skill_only(self) -> None:
        paths = tracked_files()

        forbidden = sorted(
            path
            for path in paths
            if path.endswith((".go", ".cmd", ".ps1"))
            or path in REMOVED_FILES
            or path.startswith(REMOVED_PATHS)
        )

        self.assertEqual(forbidden, [])

    def test_all_json_is_valid_when_protocol_fixtures_are_removed(self) -> None:
        json_paths = tuple(path for path in tracked_files() if path.endswith(".json"))

        invalid: list[str] = []
        for path in json_paths:
            try:
                json.loads((ROOT / path).read_text(encoding="utf-8"))
            except (json.JSONDecodeError, UnicodeDecodeError):
                invalid.append(path)

        self.assertEqual(invalid, [])

    def test_relative_markdown_links_resolve_when_repository_is_skill_only(self) -> None:
        broken: list[str] = []
        markdown_paths = tuple(path for path in tracked_files() if path.endswith(".md"))

        for path in markdown_paths:
            document = (ROOT / path).read_text(encoding="utf-8")
            for match in MARKDOWN_LINK.finditer(document):
                raw_target = match.group(1).strip().strip("<>")
                target = unquote(raw_target.split("#", 1)[0])
                if not target or "://" in target or target.startswith("mailto:"):
                    continue
                if not ((ROOT / path).parent / target).exists():
                    broken.append(f"{path}: {raw_target}")

        self.assertEqual(broken, [])

    def test_release_signing_key_is_public_only_when_retained(self) -> None:
        key = (ROOT / "RELEASE_SIGNING_KEY.asc").read_text(encoding="utf-8")

        self.assertIn("-----BEGIN PGP PUBLIC KEY BLOCK-----", key)
        self.assertNotIn("-----BEGIN PGP PRIVATE KEY BLOCK-----", key)
        self.assertNotIn("-----BEGIN PGP SECRET KEY BLOCK-----", key)

    def test_manifest_exposes_only_v02_skills_when_runtime_is_removed(self) -> None:
        manifest = json.loads(
            (ROOT / ".codex-plugin" / "plugin.json").read_text(encoding="utf-8")
        )

        self.assertEqual(manifest.get("version"), "0.2.0")
        self.assertEqual(manifest.get("skills"), "./skills/")
        for forbidden_field in ("hooks", "mcpServers", "apps"):
            self.assertNotIn(forbidden_field, manifest)

    def test_text_attributes_match_python_skill_repository(self) -> None:
        attributes = (ROOT / ".gitattributes").read_text(encoding="utf-8")

        self.assertIn("*.py text eol=lf", attributes)
        for removed_rule in ("*.go", "*.cmd", "*.ps1", "hooks/"):
            self.assertNotIn(removed_rule, attributes)

    def test_ignore_rules_match_python_test_and_package_outputs(self) -> None:
        ignore = (ROOT / ".gitignore").read_text(encoding="utf-8")

        self.assertIn("__pycache__/", ignore)
        self.assertIn("*.py[cod]", ignore)
        for go_artifact in ("/bin/", "*.test", "*.out", "coverage.out", "coverprofile.out"):
            self.assertNotIn(go_artifact, ignore)

    def test_makefile_uses_python_for_test_package_check_and_clean(self) -> None:
        makefile = (ROOT / "Makefile").read_text(encoding="utf-8")

        for target in ("test:", "package:", "check:", "clean:"):
            self.assertIn(target, makefile)
        self.assertIn("python3 -m unittest discover", makefile)
        self.assertIn("python3 scripts/package_release.py", makefile)
        self.assertIn("__pycache__", makefile)
        self.assertNotIn("go ", makefile.lower())

    def test_ci_runs_python_tests_in_one_ubuntu_job(self) -> None:
        workflow = (ROOT / ".github" / "workflows" / "ci.yml").read_text(
            encoding="utf-8"
        )
        jobs = workflow.split("jobs:\n", 1)[1]

        self.assertEqual(len(re.findall(r"^  [a-zA-Z0-9_-]+:\s*$", jobs, re.MULTILINE)), 1)
        self.assertIn("runs-on: ubuntu-latest", workflow)
        self.assertRegex(workflow, r"actions/checkout@[0-9a-f]{40}")
        self.assertIn(
            "actions/setup-python@ece7cb06caefa5fff74198d8649806c4678c61a1",
            workflow,
        )
        self.assertIn("run: make check", workflow)
        for old_surface in ("setup-go", "matrix:", "macos", "windows", "go test"):
            self.assertNotIn(old_surface, workflow.lower())


if __name__ == "__main__":
    unittest.main()
