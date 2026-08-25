package hook

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type runResult struct {
	exitCode int
	stdout   string
	stderr   string
}

func runFixture(t *testing.T, command string, fixturePath string) runResult {
	t.Helper()

	input, err := os.Open(filepath.Join("..", "..", "testdata", fixturePath))
	if err != nil {
		t.Fatalf("Given fixture %q: open: %v", fixturePath, err)
	}
	t.Cleanup(func() {
		if closeErr := input.Close(); closeErr != nil {
			t.Errorf("close fixture %q: %v", fixturePath, closeErr)
		}
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(command, Streams{
		Stdin:  input,
		Stdout: &stdout,
		Stderr: &stderr,
	})

	return runResult{
		exitCode: exitCode,
		stdout:   stdout.String(),
		stderr:   stderr.String(),
	}
}

func decodeSingleJSONObject(t *testing.T, output string) map[string]json.RawMessage {
	t.Helper()

	decoder := json.NewDecoder(bytes.NewBufferString(output))
	var object map[string]json.RawMessage
	if err := decoder.Decode(&object); err != nil {
		t.Fatalf("Then stdout must contain one JSON object: %v\nstdout: %q", err, output)
	}

	var trailing json.RawMessage
	err := decoder.Decode(&trailing)
	if err != io.EOF {
		t.Fatalf("Then stdout must contain exactly one JSON object: trailing decode error = %v", err)
	}

	return object
}

func assertSingleTrailingNewline(t *testing.T, output string) {
	t.Helper()

	if len(output) == 0 || output[len(output)-1] != '\n' {
		t.Fatalf("Then stdout must end with one newline: %q", output)
	}
	if len(output) > 1 && output[len(output)-2] == '\n' {
		t.Fatalf("Then stdout must end with exactly one newline: %q", output)
	}
}

func Test_Run_accepts_trailing_whitespace_after_single_input_object(t *testing.T) {
	// Given
	fixturePath := filepath.Join("session-start", "startup-trailing-whitespace.json")

	// When
	result := runFixture(t, "session-start", fixturePath)

	// Then
	if result.exitCode != 0 {
		t.Fatalf("Then exit code = %d, want 0", result.exitCode)
	}
	if result.stderr != "" {
		t.Fatalf("Then stderr = %q, want empty", result.stderr)
	}
	assertSingleTrailingNewline(t, result.stdout)
	decodeSingleJSONObject(t, result.stdout)
}

func Test_Run_rejects_second_input_object_without_corrupting_stdout(t *testing.T) {
	// Given
	fixturePath := filepath.Join("session-start", "second-object.json")

	// When
	result := runFixture(t, "session-start", fixturePath)

	// Then
	if result.exitCode != 0 {
		t.Fatalf("Then exit code = %d, want fail-open 0", result.exitCode)
	}
	if result.stdout != "" {
		t.Fatalf("Then stdout = %q, want empty protocol stream", result.stdout)
	}
	if result.stderr == "" {
		t.Fatal("Then stderr must contain a concise diagnostic")
	}
}

func Test_Run_reports_malformed_input_only_on_stderr(t *testing.T) {
	// Given
	fixturePath := filepath.Join("stop", "malformed.json")

	// When
	result := runFixture(t, "stop", fixturePath)

	// Then
	if result.exitCode != 0 {
		t.Fatalf("Then exit code = %d, want fail-open 0", result.exitCode)
	}
	if result.stdout != "" {
		t.Fatalf("Then stdout = %q, want empty protocol stream", result.stdout)
	}
	if result.stderr == "" {
		t.Fatal("Then stderr must contain a concise diagnostic")
	}
}

func Test_Run_rejects_oversized_input_without_corrupting_stdout(t *testing.T) {
	// Given
	input := `{"hook_event_name":"Stop","last_assistant_message":"` + strings.Repeat("a", 4*1024*1024) + `"}`
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// When
	exitCode := Run("stop", Streams{Stdin: strings.NewReader(input), Stdout: &stdout, Stderr: &stderr})

	// Then
	if exitCode != 0 {
		t.Fatalf("Then exit code = %d, want fail-open 0", exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("Then stdout = %q, want empty protocol stream", stdout.String())
	}
	if stderr.Len() == 0 {
		t.Fatal("Then stderr must contain a concise diagnostic")
	}
}
