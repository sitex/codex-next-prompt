package tests_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
