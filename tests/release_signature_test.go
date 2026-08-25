package tests_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_ReleaseSignature_rejects_bundle_with_additional_primary_key(t *testing.T) {
	// Given
	repoRoot := repositoryRoot(t)
	fixture := newReleaseSignatureFixture(t)
	expectedFingerprint := fixture.generateKey(t, "Expected Release <expected@example.invalid>")
	attackerFingerprint := fixture.generateKey(t, "Attacker <attacker@example.invalid>")
	fixture.exportPublicBundle(t, expectedFingerprint, attackerFingerprint)
	fixture.createSignedTag(t, attackerFingerprint)

	// When
	command := exec.Command(
		filepath.Join(repoRoot, "scripts", "verify-release-tag.sh"),
		fixture.publicKeyPath,
		fixture.tagName,
		expectedFingerprint,
	)
	command.Dir = fixture.repository
	command.Env = append(os.Environ(), "RUNNER_TEMP="+fixture.runnerTemp)
	output, err := command.CombinedOutput()

	// Then
	if err == nil {
		t.Fatalf("Then appended attacker primary key must be rejected; output: %s", output)
	}
	if !strings.Contains(string(output), "exactly one primary key") {
		t.Fatalf("Then rejection must identify the primary-key bundle violation; output: %s", output)
	}
}

func Test_ReleaseSignature_accepts_expected_primary_when_tag_uses_signing_subkey(t *testing.T) {
	// Given
	repoRoot := repositoryRoot(t)
	fixture := newReleaseSignatureFixture(t)
	expectedFingerprint, signingFingerprint := fixture.generatePrimaryWithSigningSubkey(t, "Expected Release <expected@example.invalid>")
	fixture.exportPublicBundle(t, expectedFingerprint)
	fixture.createSignedTag(t, signingFingerprint+"!")

	// When
	command := exec.Command(
		filepath.Join(repoRoot, "scripts", "verify-release-tag.sh"),
		fixture.publicKeyPath,
		fixture.tagName,
		expectedFingerprint,
	)
	command.Dir = fixture.repository
	command.Env = append(os.Environ(), "RUNNER_TEMP="+fixture.runnerTemp)
	output, err := command.CombinedOutput()

	// Then
	if err != nil {
		t.Fatalf("Then expected signing subkey must verify through its primary fingerprint: %v\n%s", err, output)
	}
}

type releaseSignatureFixture struct {
	gnupgHome     string
	publicKeyPath string
	repository    string
	runnerTemp    string
	tagName       string
}

func newReleaseSignatureFixture(t *testing.T) releaseSignatureFixture {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("release signature integration runs on the Linux release-verification platform")
	}
	for _, executable := range []string{"bash", "git", "gpg"} {
		if _, err := exec.LookPath(executable); err != nil {
			t.Skipf("release signature integration requires %s: %v", executable, err)
		}
	}
	root := t.TempDir()
	gnupgHome := filepath.Join(root, "signing-gnupg")
	if err := os.Mkdir(gnupgHome, 0o700); err != nil {
		t.Fatalf("Given create signing GNUPGHOME: %v", err)
	}
	repository := filepath.Join(root, "repository")
	if err := os.Mkdir(repository, 0o755); err != nil {
		t.Fatalf("Given create disposable repository: %v", err)
	}
	runnerTemp := filepath.Join(root, "runner-temp")
	if err := os.Mkdir(runnerTemp, 0o755); err != nil {
		t.Fatalf("Given create runner temp: %v", err)
	}
	runCommand(t, repository, os.Environ(), "git", "init", "--quiet")
	runCommand(t, repository, os.Environ(), "git", "config", "user.name", "Release Test")
	runCommand(t, repository, os.Environ(), "git", "config", "user.email", "release@example.invalid")
	if err := os.WriteFile(filepath.Join(repository, "payload.txt"), []byte("release payload\n"), 0o644); err != nil {
		t.Fatalf("Given write release payload: %v", err)
	}
	runCommand(t, repository, os.Environ(), "git", "add", "payload.txt")
	runCommand(t, repository, os.Environ(), "git", "commit", "--quiet", "-m", "release payload")
	return releaseSignatureFixture{
		gnupgHome:     gnupgHome,
		publicKeyPath: filepath.Join(root, "release-keys.asc"),
		repository:    repository,
		runnerTemp:    runnerTemp,
		tagName:       "v1.0.0",
	}
}

func (fixture releaseSignatureFixture) generateKey(t *testing.T, identity string) string {
	t.Helper()
	environment := append(os.Environ(), "GNUPGHOME="+fixture.gnupgHome)
	runCommand(t, fixture.repository, environment, "gpg", "--batch", "--passphrase", "", "--quick-gen-key", identity, "ed25519", "sign", "0")
	output := runCommand(t, fixture.repository, environment, "gpg", "--batch", "--with-colons", "--fingerprint", identity)
	return firstFingerprint(t, output, identity)
}

func (fixture releaseSignatureFixture) generatePrimaryWithSigningSubkey(t *testing.T, identity string) (string, string) {
	t.Helper()
	environment := append(os.Environ(), "GNUPGHOME="+fixture.gnupgHome)
	runCommand(t, fixture.repository, environment, "gpg", "--batch", "--passphrase", "", "--quick-gen-key", identity, "ed25519", "cert", "0")
	primaryFingerprint := firstFingerprint(t, runCommand(t, fixture.repository, environment, "gpg", "--batch", "--with-colons", "--fingerprint", identity), identity)
	runCommand(t, fixture.repository, environment, "gpg", "--batch", "--passphrase", "", "--quick-add-key", primaryFingerprint, "ed25519", "sign", "0")
	output := runCommand(t, fixture.repository, environment, "gpg", "--batch", "--with-colons", "--fingerprint", "--fingerprint", identity)
	fingerprints := make([]string, 0, 2)
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Split(line, ":")
		if len(fields) > 9 && fields[0] == "fpr" {
			fingerprints = append(fingerprints, fields[9])
		}
	}
	if len(fingerprints) != 2 {
		t.Fatalf("Given generated key %q must expose primary and signing-subkey fingerprints; got %v", identity, fingerprints)
	}
	return fingerprints[0], fingerprints[1]
}

func firstFingerprint(t *testing.T, output, identity string) string {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Split(line, ":")
		if len(fields) > 9 && fields[0] == "fpr" {
			return fields[9]
		}
	}
	t.Fatalf("Given generated key %q must expose a fingerprint", identity)
	return ""
}

func (fixture releaseSignatureFixture) exportPublicBundle(t *testing.T, fingerprints ...string) {
	t.Helper()
	environment := append(os.Environ(), "GNUPGHOME="+fixture.gnupgHome)
	arguments := append([]string{"--batch", "--armor", "--export"}, fingerprints...)
	command := exec.Command("gpg", arguments...)
	command.Dir = fixture.repository
	command.Env = environment
	output, err := command.Output()
	if err != nil {
		t.Fatalf("Given export public key bundle: %v", err)
	}
	if err := os.WriteFile(fixture.publicKeyPath, output, 0o644); err != nil {
		t.Fatalf("Given write public key bundle: %v", err)
	}
}

func (fixture releaseSignatureFixture) createSignedTag(t *testing.T, fingerprint string) {
	t.Helper()
	environment := append(os.Environ(), "GNUPGHOME="+fixture.gnupgHome)
	runCommand(t, fixture.repository, environment, "git", "-c", "user.signingkey="+fingerprint, "tag", "-s", fixture.tagName, "-m", "release")
}

func runCommand(t *testing.T, directory string, environment []string, name string, arguments ...string) string {
	t.Helper()
	command := exec.Command(name, arguments...)
	command.Dir = directory
	command.Env = environment
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("Given run %s: %v\n%s", fmt.Sprintf("%s %s", name, strings.Join(arguments, " ")), err, output)
	}
	return string(output)
}
