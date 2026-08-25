package tests_test

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_POSIXArchiveSmoke_executes_extracted_native_launcher(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX archive smoke runs on Linux and macOS")
	}

	// Given
	repoRoot := repositoryRoot(t)
	distDir := t.TempDir()
	version := "1.2.3-test"
	runPackageRelease(t, repoRoot, distDir, version, runtime.GOOS, runtime.GOARCH)
	archivePath := filepath.Join(distDir, "codex-next-prompt-"+version+"-"+runtime.GOOS+"-"+runtime.GOARCH+".tar.gz")
	command := exec.Command(filepath.Join(repoRoot, "tests", "smoke-archive-posix.sh"), archivePath)
	command.Dir = repoRoot
	command.Env = os.Environ()

	// When
	output, err := command.CombinedOutput()

	// Then
	if err != nil {
		t.Fatalf("Then extracted native archive smoke must pass: %v\n%s", err, output)
	}
}

func Test_POSIXArchiveSmoke_rejects_unsafe_entries_before_extraction(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX archive smoke runs on Linux and macOS")
	}

	tests := []struct {
		name     string
		header   tar.Header
		contents string
	}{
		{name: "parent traversal", header: tar.Header{Name: "../escape", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg}, contents: "x"},
		{name: "absolute path", header: tar.Header{Name: "/tmp/escape", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg}, contents: "x"},
		{name: "symbolic link", header: tar.Header{Name: "plugin/link", Linkname: "target", Mode: 0o777, Typeflag: tar.TypeSymlink}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Given
			archivePath := filepath.Join(t.TempDir(), "unsafe.tar.gz")
			writeTarGzEntry(t, archivePath, test.header, test.contents)
			command := exec.Command(filepath.Join(repositoryRoot(t), "tests", "smoke-archive-posix.sh"), archivePath)

			// When
			output, err := command.CombinedOutput()

			// Then
			if err == nil {
				t.Fatal("Then unsafe archive entry must be rejected")
			}
			if !strings.Contains(string(output), "unsafe archive") {
				t.Fatalf("Then rejection must occur during archive preflight: %s", output)
			}
		})
	}
}

func writeTarGzEntry(t *testing.T, archivePath string, header tar.Header, contents string) {
	t.Helper()

	archive, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("Given create archive: %v", err)
	}
	gzipWriter := gzip.NewWriter(archive)
	tarWriter := tar.NewWriter(gzipWriter)
	if err := tarWriter.WriteHeader(&header); err != nil {
		t.Fatalf("Given write tar header: %v", err)
	}
	if contents != "" {
		if _, err := tarWriter.Write([]byte(contents)); err != nil {
			t.Fatalf("Given write tar contents: %v", err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("Given close tar writer: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatalf("Given close gzip writer: %v", err)
	}
	if err := archive.Close(); err != nil {
		t.Fatalf("Given close archive: %v", err)
	}
}
