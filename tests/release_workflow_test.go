package tests_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_ReleaseWorkflow_publishes_verified_native_smoked_tag_artifacts(t *testing.T) {
	// Given
	workflowPath := filepath.Join("..", ".github", "workflows", "release.yml")
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("Given read release workflow: %v", err)
	}
	workflow := string(data)

	// When
	requiredTokens := []string{
		"git verify-tag \"$GITHUB_REF_NAME\"",
		"git rev-list -n 1 \"$GITHUB_REF_NAME\"",
		"RELEASE_VERSION=%s",
		"macos-15-intel",
		"macos-latest",
		"ubuntu-24.04-arm",
		"scripts/package-release.sh \"$RELEASE_VERSION\"",
		"smoke-windows.ps1",
		"smoke-windows:",
		"needs: [package, smoke-windows]",
		"gh release view \"$GITHUB_REF_NAME\"",
		"--notes-file RELEASE_NOTES.md",
	}
	forbiddenTokens := []string{"workflow_dispatch", "--clobber", "--generate-notes", "${{ steps.version.outputs.version }}"}

	// Then
	for _, token := range requiredTokens {
		if !strings.Contains(workflow, token) {
			t.Errorf("Then release workflow must contain structural token %q", token)
		}
	}
	for _, token := range forbiddenTokens {
		if strings.Contains(workflow, token) {
			t.Errorf("Then release workflow must omit structural token %q", token)
		}
	}
}
