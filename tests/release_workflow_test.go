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
		"permissions:\n  contents: read",
		"publish:\n    name: Publish immutable GitHub release\n    permissions:\n      contents: write",
		`git fetch --force --no-tags origin "refs/tags/$GITHUB_REF_NAME:refs/tags/$GITHUB_REF_NAME"`,
		"scripts/verify-release-tag.sh RELEASE_SIGNING_KEY.asc \"$GITHUB_REF_NAME\"",
		"BFA7D43C126EE54A5FC8DD0EBE645A3EFA752D77",
		"git fetch --no-tags origin main",
		"git merge-base --is-ancestor \"$tag_commit\" origin/main",
		"git rev-list -n 1 \"$GITHUB_REF_NAME\"",
		"RELEASE_VERSION=%s",
		"macos-15-intel",
		"macos-latest",
		"ubuntu-24.04-arm",
		"scripts/package-release.sh \"$RELEASE_VERSION\"",
		"smoke-windows.ps1",
		`$archives = @(Get-ChildItem dist/*.zip -File)`,
		`if ($archives.Count -ne 1)`,
		`tests/smoke-windows.ps1 $archives[0].FullName`,
		"smoke-windows:",
		"needs: [package, smoke-windows]",
		"gh api -i \"repos/$GITHUB_REPOSITORY/releases/tags/$GITHUB_REF_NAME\"",
		"http_status=",
		"case \"$http_status\" in",
		"404)",
		"200)",
		"--notes-file RELEASE_NOTES.md",
	}
	forbiddenTokens := []string{"workflow_dispatch", "--clobber", "--generate-notes", "${{ steps.version.outputs.version }}", "-Single"}

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

func Test_ReleaseTagVerifier_binds_single_VALID_SIGNATURE_to_expected_primary_key(t *testing.T) {
	// Given
	scriptPath := filepath.Join("..", "scripts", "verify-release-tag.sh")
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Given read release tag verifier: %v", err)
	}
	script := string(data)

	// When
	requiredTokens := []string{
		"mktemp -d",
		"gpg --batch --with-colons --show-keys --fingerprint",
		"${#primary_fingerprints[@]} -ne 1",
		"git verify-tag --raw",
		"$2 == \"VALIDSIG\" {print $NF}",
		"${#valid_primary_fingerprints[@]} -ne 1",
		"${valid_primary_fingerprints[0]} != \"$expected_primary_fingerprint\"",
	}

	// Then
	for _, token := range requiredTokens {
		if !strings.Contains(script, token) {
			t.Errorf("Then release verifier must contain structural token %q", token)
		}
	}
}

func Test_Workflows_pin_actions_and_disable_checkout_credentials(t *testing.T) {
	// Given
	workflowPaths := []string{
		filepath.Join("..", ".github", "workflows", "ci.yml"),
		filepath.Join("..", ".github", "workflows", "release.yml"),
	}
	pins := map[string]string{
		"actions/checkout@":          "actions/checkout@11d5960a326750d5838078e36cf38b85af677262",
		"actions/setup-go@":          "actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16",
		"actions/upload-artifact@":   "actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02",
		"actions/download-artifact@": "actions/download-artifact@634f93cb2916e3fdff6788551b99b062d0335ce0",
	}

	for _, workflowPath := range workflowPaths {
		workflowData, err := os.ReadFile(workflowPath)
		if err != nil {
			t.Fatalf("Given read workflow %q: %v", workflowPath, err)
		}
		workflow := string(workflowData)

		// When/Then
		for actionPrefix, pinnedAction := range pins {
			if strings.Count(workflow, actionPrefix) != strings.Count(workflow, pinnedAction) {
				t.Errorf("Then every %s use in %s must use reviewed SHA", actionPrefix, workflowPath)
			}
		}
		checkoutCount := strings.Count(workflow, "actions/checkout@")
		if strings.Count(workflow, "persist-credentials: false") != checkoutCount {
			t.Errorf("Then %s checkout credential guards = %d, want %d", workflowPath, strings.Count(workflow, "persist-credentials: false"), checkoutCount)
		}
	}
}
