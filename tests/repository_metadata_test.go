package tests_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var markdownLink = regexp.MustCompile(`\[[^]]+\]\(([^)]+)\)`)

func Test_RepositoryMetadata_all_tracked_JSON_is_valid(t *testing.T) {
	// Given
	repoRoot := repositoryRoot(t)
	paths := trackedFiles(t, repoRoot, "*.json")

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			if strings.Contains(path, "testdata/") {
				t.Skip("protocol fixtures intentionally include malformed and multi-object JSON")
			}
			data, err := os.ReadFile(filepath.Join(repoRoot, path))
			if err != nil {
				t.Fatalf("Given read JSON: %v", err)
			}

			// When
			valid := json.Valid(data)

			// Then
			if !valid {
				t.Fatal("Then tracked JSON must be valid")
			}
		})
	}
}

func Test_RepositoryMetadata_local_Markdown_links_resolve(t *testing.T) {
	// Given
	repoRoot := repositoryRoot(t)
	paths := trackedFiles(t, repoRoot, "*.md")

	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(repoRoot, path))
		if err != nil {
			t.Fatalf("Given read Markdown %q: %v", path, err)
		}
		for _, match := range markdownLink.FindAllStringSubmatch(string(data), -1) {
			target := strings.SplitN(match[1], "#", 2)[0]
			if target == "" || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
				continue
			}

			// When
			_, statErr := os.Stat(filepath.Join(repoRoot, filepath.Dir(path), filepath.FromSlash(target)))

			// Then
			if statErr != nil {
				t.Errorf("Then local link %q in %s must resolve: %v", match[1], path, statErr)
			}
		}
	}
}

func trackedFiles(t *testing.T, repoRoot, pattern string) []string {
	t.Helper()
	command := exec.Command("git", "ls-files", pattern)
	command.Dir = repoRoot
	output, err := command.Output()
	if err != nil {
		t.Fatalf("Given list tracked %s files: %v", pattern, err)
	}
	return strings.Fields(string(output))
}
