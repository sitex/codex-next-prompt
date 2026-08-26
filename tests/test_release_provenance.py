import os
import re
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from typing import Final


ROOT: Final = Path(__file__).resolve().parents[1]
VERIFIER: Final = ROOT / "scripts" / "verify-release-tag.sh"
EXPECTED_FINGERPRINT: Final = "BFA7D43C126EE54A5FC8DD0EBE645A3EFA752D77"
PINNED_ACTION: Final = re.compile(r"^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+@[0-9a-f]{40}(?:\s+#.*)?$")


def run(command: list[str], cwd: Path, environment: dict[str, str]) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        command,
        cwd=cwd,
        env=environment,
        capture_output=True,
        text=True,
        check=False,
    )


@unittest.skipUnless(sys.platform.startswith("linux"), "release provenance runs on Linux")
class TestReleaseSignature(unittest.TestCase):
    def setUp(self) -> None:
        for executable in ("bash", "git", "gpg"):
            if shutil.which(executable) is None:
                self.skipTest(f"release provenance requires {executable}")

    def test_verifier_rejects_attacker_primary_in_public_bundle(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            root, repository, environment = self._repository(Path(temporary))
            expected = self._generate_key(repository, environment, "Expected <expected@example.invalid>", "sign")
            attacker = self._generate_key(repository, environment, "Attacker <attacker@example.invalid>", "sign")
            public_key = root / "release-keys.asc"
            self._export(repository, environment, public_key, expected, attacker)
            self._tag(repository, environment, attacker)

            result = run(
                [str(VERIFIER), str(public_key), "v1.0.0", expected],
                repository,
                {**os.environ, "RUNNER_TEMP": str(root)},
            )

            self.assertNotEqual(result.returncode, 0)
            self.assertIn("exactly one primary key", result.stderr)

    def test_verifier_accepts_tag_signed_by_expected_signing_subkey(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            root, repository, environment = self._repository(Path(temporary))
            identity = "Expected <expected@example.invalid>"
            primary = self._generate_key(repository, environment, identity, "cert")
            added = run(
                ["gpg", "--batch", "--passphrase", "", "--quick-add-key", primary, "ed25519", "sign", "0"],
                repository,
                environment,
            )
            self.assertEqual(added.returncode, 0, added.stderr)
            fingerprints = self._fingerprints(repository, environment, identity)
            self.assertEqual(len(fingerprints), 2)
            public_key = root / "release-key.asc"
            self._export(repository, environment, public_key, primary)
            self._tag(repository, environment, f"{fingerprints[1]}!")

            result = run(
                [str(VERIFIER), str(public_key), "v1.0.0", primary],
                repository,
                {**os.environ, "RUNNER_TEMP": str(root)},
            )

            self.assertEqual(result.returncode, 0, result.stderr)

    def _repository(self, root: Path) -> tuple[Path, Path, dict[str, str]]:
        gnupg_home = root / "gnupg"
        repository = root / "repository"
        gnupg_home.mkdir(mode=0o700)
        repository.mkdir()
        environment = {**os.environ, "GNUPGHOME": str(gnupg_home)}
        commands = (
            ["git", "init", "--quiet"],
            ["git", "config", "user.name", "Release Test"],
            ["git", "config", "user.email", "release@example.invalid"],
        )
        for command in commands:
            result = run(command, repository, environment)
            self.assertEqual(result.returncode, 0, result.stderr)
        (repository / "payload.txt").write_text("release payload\n", encoding="utf-8")
        for command in (["git", "add", "payload.txt"], ["git", "commit", "--quiet", "-m", "payload"]):
            result = run(command, repository, environment)
            self.assertEqual(result.returncode, 0, result.stderr)
        return root, repository, environment

    def _generate_key(self, repository: Path, environment: dict[str, str], identity: str, usage: str) -> str:
        result = run(
            ["gpg", "--batch", "--passphrase", "", "--quick-gen-key", identity, "ed25519", usage, "0"],
            repository,
            environment,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        return self._fingerprints(repository, environment, identity)[0]

    def _fingerprints(self, repository: Path, environment: dict[str, str], identity: str) -> tuple[str, ...]:
        result = run(
            ["gpg", "--batch", "--with-colons", "--fingerprint", "--fingerprint", identity],
            repository,
            environment,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        return tuple(line.split(":")[9] for line in result.stdout.splitlines() if line.startswith("fpr:"))

    def _export(self, repository: Path, environment: dict[str, str], path: Path, *fingerprints: str) -> None:
        result = run(["gpg", "--batch", "--armor", "--export", *fingerprints], repository, environment)
        self.assertEqual(result.returncode, 0, result.stderr)
        path.write_text(result.stdout, encoding="ascii")

    def _tag(self, repository: Path, environment: dict[str, str], fingerprint: str) -> None:
        result = run(
            ["git", "-c", f"user.signingkey={fingerprint}", "tag", "-s", "v1.0.0", "-m", "release"],
            repository,
            environment,
        )
        self.assertEqual(result.returncode, 0, result.stderr)


class TestReleaseWorkflow(unittest.TestCase):
    def test_release_workflow_enforces_provenance_and_exact_assets(self) -> None:
        workflow = (ROOT / ".github" / "workflows" / "release.yml").read_text(encoding="utf-8")

        required = (
            "permissions:\n  contents: read",
            "runs-on: ubuntu-latest",
            "persist-credentials: false",
            'git fetch --force --no-tags origin "refs/tags/$GITHUB_REF_NAME:refs/tags/$GITHUB_REF_NAME"',
            "scripts/verify-release-tag.sh RELEASE_SIGNING_KEY.asc \"$GITHUB_REF_NAME\"",
            EXPECTED_FINGERPRINT,
            'test "$tag_commit" = "$GITHUB_SHA"',
            'git merge-base --is-ancestor "$tag_commit" origin/main',
            'test "$(cat VERSION)" = "$version"',
            "python3 scripts/package_release.py \"$version\"",
            "gh api -i \"repos/$GITHUB_REPOSITORY/releases/tags/$GITHUB_REF_NAME\"",
            "dist/codex-next-prompt-${{ env.RELEASE_VERSION }}.zip",
            "dist/codex-next-prompt-${{ env.RELEASE_VERSION }}.zip.sha256",
        )
        for token in required:
            with self.subTest(token=token):
                self.assertIn(token, workflow)
        for token in ("matrix:", "setup-go", "windows", "macos", "dist/*"):
            with self.subTest(forbidden=token):
                self.assertNotIn(token, workflow.lower())
        self.assertIn("publish:\n    permissions:\n      contents: write", workflow)
        self.assertEqual(workflow.count("contents: write"), 1)

    def test_publish_checks_out_event_sha_before_creating_release(self) -> None:
        workflow = (ROOT / ".github" / "workflows" / "release.yml").read_text(encoding="utf-8")
        publish = workflow.split("  publish:\n", 1)[1]

        checkout = publish.find("uses: actions/checkout@")
        publish_command = publish.find("gh release create")
        self.assertGreaterEqual(checkout, 0)
        self.assertGreater(publish_command, checkout)
        self.assertIn("ref: ${{ github.sha }}", publish[:publish_command])
        self.assertIn("persist-credentials: false", publish[:publish_command])
        self.assertIn("--notes-file RELEASE_NOTES.md", publish[publish_command:])

    def test_package_checks_out_verified_event_sha(self) -> None:
        workflow = (ROOT / ".github" / "workflows" / "release.yml").read_text(encoding="utf-8")
        package = workflow.split("  package:\n", 1)[1].split("  publish:\n", 1)[0]

        self.assertIn("ref: ${{ github.sha }}", package)
        self.assertIn("persist-credentials: false", package)

    def test_workflows_pin_actions_and_disable_checkout_credentials(self) -> None:
        for relative_path in (".github/workflows/ci.yml", ".github/workflows/release.yml"):
            workflow = (ROOT / relative_path).read_text(encoding="utf-8")
            action_lines = tuple(line.strip()[len("uses: ") :] for line in workflow.splitlines() if line.strip().startswith("uses: "))
            with self.subTest(path=relative_path):
                self.assertTrue(action_lines)
                self.assertTrue(all(PINNED_ACTION.fullmatch(line) for line in action_lines))
                self.assertEqual(workflow.count("persist-credentials: false"), workflow.count("actions/checkout@"))


if __name__ == "__main__":
    unittest.main()
