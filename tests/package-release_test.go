package tests_test

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type archiveEntry struct {
	mode os.FileMode
	size int64
}

func Test_PackageRelease_creates_target_specific_unix_archive_and_checksum(t *testing.T) {
	requirePOSIXPackagingHost(t)

	// Given
	repoRoot := repositoryRoot(t)
	distDir := t.TempDir()
	version := "1.2.3-test"

	// When
	runPackageRelease(t, repoRoot, distDir, version, "linux", "amd64")

	// Then
	archivePath := filepath.Join(distDir, "codex-next-prompt-"+version+"-linux-amd64.tar.gz")
	entries := readTarGz(t, archivePath)
	assertReleaseEntries(t, entries, "linux", "amd64")
	assertExecutable(t, entries, "codex-next-prompt-"+version+"/hooks/run")
	assertExecutable(t, entries, "codex-next-prompt-"+version+"/bin/linux-amd64/codex-next-prompt")
	assertChecksum(t, archivePath)
}

func Test_PackageRelease_creates_target_specific_windows_archive_and_checksum(t *testing.T) {
	requirePOSIXPackagingHost(t)

	// Given
	repoRoot := repositoryRoot(t)
	distDir := t.TempDir()
	version := "1.2.3-test"

	// When
	runPackageRelease(t, repoRoot, distDir, version, "windows", "arm64")

	// Then
	archivePath := filepath.Join(distDir, "codex-next-prompt-"+version+"-windows-arm64.zip")
	entries := readZip(t, archivePath)
	assertReleaseEntries(t, entries, "windows", "arm64")
	assertChecksum(t, archivePath)
}

func Test_PackageRelease_rejects_invalid_version_and_target(t *testing.T) {
	requirePOSIXPackagingHost(t)

	tests := []struct {
		name    string
		version string
		goos    string
		goarch  string
	}{
		{name: "path separator in version", version: "1.2/3", goos: "linux", goarch: "amd64"},
		{name: "unsupported operating system", version: "1.2.3", goos: "plan9", goarch: "amd64"},
		{name: "unsupported architecture", version: "1.2.3", goos: "linux", goarch: "386"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			repoRoot := repositoryRoot(t)
			distDir := t.TempDir()
			command := exec.Command(filepath.Join(repoRoot, "scripts", "package-release.sh"), test.version, test.goos, test.goarch)
			command.Dir = repoRoot
			command.Env = append(os.Environ(), "DIST_DIR="+distDir)

			// When
			err := command.Run()

			// Then
			if err == nil {
				t.Fatal("Then invalid packaging input must fail")
			}
		})
	}
}

func runPackageRelease(t *testing.T, repoRoot, distDir, version, goos, goarch string) {
	t.Helper()

	command := exec.Command(filepath.Join(repoRoot, "scripts", "package-release.sh"), version, goos, goarch)
	command.Dir = repoRoot
	command.Env = append(os.Environ(), "DIST_DIR="+distDir)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("When package %s/%s: %v\n%s", goos, goarch, err, output)
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("Given resolve repository root: %v", err)
	}
	return repoRoot
}

func requirePOSIXPackagingHost(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("release packaging script is exercised on POSIX CI hosts")
	}
}

func readTarGz(t *testing.T, archivePath string) map[string]archiveEntry {
	t.Helper()

	file, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("Then open archive %q: %v", archivePath, err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("Then open gzip stream %q: %v", archivePath, err)
	}
	defer gzipReader.Close()

	entries := make(map[string]archiveEntry)
	tarReader := tar.NewReader(gzipReader)
	for {
		header, nextErr := tarReader.Next()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			t.Fatalf("Then read tar archive %q: %v", archivePath, nextErr)
		}
		if header.Typeflag == tar.TypeReg {
			entries[strings.TrimPrefix(header.Name, "./")] = archiveEntry{mode: os.FileMode(header.Mode), size: header.Size}
		}
	}
	return entries
}

func readZip(t *testing.T, archivePath string) map[string]archiveEntry {
	t.Helper()

	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatalf("Then open zip archive %q: %v", archivePath, err)
	}
	defer reader.Close()

	entries := make(map[string]archiveEntry)
	for _, file := range reader.File {
		if !file.FileInfo().IsDir() {
			entries[file.Name] = archiveEntry{mode: file.Mode(), size: int64(file.UncompressedSize64)}
		}
	}
	return entries
}

func assertReleaseEntries(t *testing.T, entries map[string]archiveEntry, goos, goarch string) {
	t.Helper()

	root := "codex-next-prompt-1.2.3-test/"
	binaryName := "codex-next-prompt"
	if goos == "windows" {
		binaryName += ".exe"
	}
	want := map[string]bool{
		root + ".agents/plugins/marketplace.json":              true,
		root + ".codex-plugin/plugin.json":                     true,
		root + "hooks/hooks.json":                              true,
		root + "hooks/run":                                     true,
		root + "hooks/run.cmd":                                 true,
		root + "bin/" + goos + "-" + goarch + "/" + binaryName: true,
	}
	for path, entry := range entries {
		if !want[path] {
			t.Errorf("Then archive contains unexpected file %q", path)
		}
		if entry.size == 0 {
			t.Errorf("Then archive file %q must not be empty", path)
		}
		for _, forbidden := range []string{".git/", "internal/", "cmd/", "tests/", "testdata/", "thoughts/", "README.md", "go.mod", "go.sum"} {
			if strings.Contains(path, forbidden) {
				t.Errorf("Then archive must not include source/private path %q", path)
			}
		}
	}
	for path := range want {
		if _, exists := entries[path]; !exists {
			t.Errorf("Then archive is missing %q", path)
		}
	}
}

func assertExecutable(t *testing.T, entries map[string]archiveEntry, path string) {
	t.Helper()

	entry, exists := entries[path]
	if !exists {
		t.Fatalf("Then executable entry %q must exist", path)
	}
	if entry.mode.Perm()&0o111 == 0 {
		t.Errorf("Then archive entry %q mode = %o, want executable", path, entry.mode.Perm())
	}
}

func assertChecksum(t *testing.T, archivePath string) {
	t.Helper()

	archive, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("Then read archive for checksum: %v", err)
	}
	checksum, err := os.ReadFile(archivePath + ".sha256")
	if err != nil {
		t.Fatalf("Then read checksum file: %v", err)
	}
	actualHash := sha256.Sum256(archive)
	want := fmt.Sprintf("%x  %s\n", actualHash, filepath.Base(archivePath))
	if string(checksum) != want {
		t.Errorf("Then checksum = %q, want %q", checksum, want)
	}
}
