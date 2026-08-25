package tests_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_ProductBuilds_disable_VCS_stamping_for_linked_worktrees(t *testing.T) {
	// Given
	repoRoot := repositoryRoot(t)
	buildFiles := []string{
		"Makefile",
		filepath.Join("scripts", "package-release.sh"),
		filepath.Join("tests", "smoke-posix.sh"),
		filepath.Join("tests", "smoke-windows.ps1"),
		filepath.Join(".github", "workflows", "ci.yml"),
	}

	for _, relativePath := range buildFiles {
		t.Run(relativePath, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(repoRoot, relativePath))
			if err != nil {
				t.Fatalf("Given read build surface: %v", err)
			}

			// When
			buildCount := strings.Count(string(data), "go build")
			disabledCount := strings.Count(string(data), "-buildvcs=false")

			// Then
			if buildCount == 0 || disabledCount < buildCount {
				t.Fatalf("Then %s must disable VCS stamping on all %d Go builds, got %d flags", relativePath, buildCount, disabledCount)
			}
		})
	}
}
