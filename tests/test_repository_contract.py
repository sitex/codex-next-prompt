import json
import re
import subprocess
import unittest
from pathlib import Path
from typing import Final
from urllib.parse import unquote


ROOT: Final = Path(__file__).resolve().parents[1]
MARKDOWN_LINK: Final = re.compile(r"!?\[[^]]*\]\(([^)]+)\)")
FORBIDDEN_METADATA_PATH: Final = "agents/" + "openai.yaml"
FORBIDDEN_POLICY: Final = "allow_" + "implicit_invocation"
REMOVED_PATHS: Final = (
    "cmd/",
    "hooks/",
    "internal/",
    "testdata/",
)
REMOVED_FILES: Final = {
    ".agents/plugins/marketplace.json",
    ".codex-plugin/plugin.json",
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
CURRENT_DOCUMENTS: Final = (
    "README.md",
    "EXAMPLES.md",
    "INSTALL_WITH_LLM.md",
    "RELEASE_NOTES.md",
    "CONTRIBUTING.md",
    "AGENTS.md",
    "SECURITY.md",
    ".github/pull_request_template.md",
    ".github/ISSUE_TEMPLATE/bug_report.md",
    ".github/ISSUE_TEMPLATE/feature_request.md",
)
INSTALL_POLICY_DOCUMENTS: Final = (
    "README.md",
    "INSTALL_WITH_LLM.md",
    "RELEASE_NOTES.md",
    "CONTRIBUTING.md",
    ".github/pull_request_template.md",
)
REMOVED_DOCUMENTATION_TERMS: Final = (
    "/hooks",
    "invoke `/next`",
    "`/next` command",
    "`/next` invocation",
    "codex plugin add",
    "codex plugin list",
    "codex plugin marketplace add",
    "codex plugin marketplace list",
    "plugin is enabled",
    "plugin metadata",
    "skill-only plugin",
    "Suggested next prompt:",
    "SessionStart",
    "Stop hook",
    "hook trust",
    "lifecycle hook",
    ".tar.gz",
    "linux-amd64",
    "darwin-amd64",
    "windows-amd64",
    "GOOS",
    "GOARCH",
    "Go 1.",
    "go test",
    "scripts/package-release.sh",
    FORBIDDEN_METADATA_PATH,
    FORBIDDEN_POLICY,
    "tree/main/skills/next",
    "tree/v0.2.0",
    "/tree/",
    "$skill-installer",
    "custom prompt",
    "plugin marketplace",
)


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

    def test_version_file_is_authoritative_for_v02_release(self) -> None:
        version = (ROOT / "VERSION").read_text(encoding="ascii")

        self.assertEqual(version, "0.2.0\n")

    def test_current_docs_route_only_to_explicit_v02_next_skill(self) -> None:
        violations: list[str] = []

        for path in CURRENT_DOCUMENTS:
            document = (ROOT / path).read_text(encoding="utf-8")
            if "0.2.0" not in document:
                violations.append(f"{path}: missing current version 0.2.0")
            if "$next" not in document:
                violations.append(f"{path}: missing canonical $next invocation")
            if "standalone skill" not in document.lower():
                violations.append(f"{path}: missing standalone skill contract")
            if "0.1.0" in document:
                violations.append(f"{path}: v0.1 history belongs in CHANGELOG.md")
            for term in REMOVED_DOCUMENTATION_TERMS:
                if term in document:
                    violations.append(f"{path}: contains removed term {term!r}")

        changelog = (ROOT / "CHANGELOG.md").read_text(encoding="utf-8")
        for required_term in ("[Unreleased]", "[0.2.0]", "[0.1.0]", "$next"):
            if required_term not in changelog:
                violations.append(f"CHANGELOG.md: missing {required_term!r}")

        self.assertEqual(violations, [])

    def test_current_install_docs_require_verified_release_zip_only(self) -> None:
        violations: list[str] = []

        for path in INSTALL_POLICY_DOCUMENTS:
            document = (ROOT / path).read_text(encoding="utf-8")
            for required_term in (
                "codex-next-prompt-0.2.0.zip",
                "codex-next-prompt-0.2.0.zip.sha256",
                "${CODEX_HOME:-$HOME/.codex}/skills/next",
            ):
                if required_term not in document:
                    violations.append(f"{path}: missing {required_term!r}")

        self.assertEqual(violations, [])

    def test_repository_omits_incompatible_openai_metadata(self) -> None:
        violations: list[str] = []

        for path in tracked_files():
            if path.startswith("tests/"):
                continue
            document = (ROOT / path).read_text(encoding="utf-8", errors="replace")
            if path.endswith(FORBIDDEN_METADATA_PATH):
                violations.append(f"{path}: incompatible metadata file is present")
            if FORBIDDEN_POLICY in document:
                violations.append(f"{path}: incompatible policy is present")

        self.assertEqual(violations, [])

    def test_runtime_matrix_records_codex_0149_explicit_activation_contract(self) -> None:
        matrix = (ROOT / "tests" / "skill_scenarios.md").read_text(encoding="utf-8")

        for contract_term in (
            "Codex 0.149.1",
            "frontmatter description",
            "explicit catalog",
            "$next",
        ):
            with self.subTest(contract_term=contract_term):
                self.assertIn(contract_term, matrix)

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
        self.assertIn("VERSION := $(shell cat VERSION)", makefile)
        self.assertNotIn(".codex-plugin", makefile)
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
