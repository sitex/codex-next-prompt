package hook

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

type sessionStartEnvelope struct {
	HookSpecificOutput struct {
		HookEventName     string `json:"hookEventName"`
		AdditionalContext string `json:"additionalContext"`
	} `json:"hookSpecificOutput"`
}

func decodeSessionStartOutput(t *testing.T, output string) sessionStartEnvelope {
	t.Helper()

	decodeSingleJSONObject(t, output)

	var envelope sessionStartEnvelope
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("Then decode SessionStart output: %v", err)
	}
	return envelope
}

func Test_Run_session_start_emits_guidance_for_supported_sources(t *testing.T) {
	// Given
	sources := []string{"startup", "resume", "clear", "compact"}

	for _, source := range sources {
		t.Run(source, func(t *testing.T) {
			fixturePath := filepath.Join("session-start", source+".json")

			// When
			result := runFixture(t, "session-start", fixturePath)

			// Then
			if result.exitCode != 0 {
				t.Fatalf("Then exit code = %d, want 0", result.exitCode)
			}
			if result.stderr != "" {
				t.Fatalf("Then stderr = %q, want empty", result.stderr)
			}

			envelope := decodeSessionStartOutput(t, result.stdout)
			if envelope.HookSpecificOutput.HookEventName != "SessionStart" {
				t.Fatalf("Then hookEventName = %q, want SessionStart", envelope.HookSpecificOutput.HookEventName)
			}
			if envelope.HookSpecificOutput.AdditionalContext == "" {
				t.Fatal("Then additionalContext must be non-empty")
			}
			if strings.Count(envelope.HookSpecificOutput.AdditionalContext, "Suggested next prompt:") != 1 {
				t.Fatalf("Then additionalContext must contain the public marker exactly once")
			}
		})
	}
}

func Test_Run_session_start_ignores_unknown_fields(t *testing.T) {
	// Given
	fixturePath := filepath.Join("session-start", "unknown-fields.json")

	// When
	result := runFixture(t, "session-start", fixturePath)

	// Then
	if result.exitCode != 0 {
		t.Fatalf("Then exit code = %d, want 0", result.exitCode)
	}
	if result.stderr != "" {
		t.Fatalf("Then stderr = %q, want empty", result.stderr)
	}
	decodeSessionStartOutput(t, result.stdout)
}

func Test_Run_session_start_fails_open_for_wrong_event_name(t *testing.T) {
	// Given
	fixturePath := filepath.Join("session-start", "wrong-event.json")

	// When
	result := runFixture(t, "session-start", fixturePath)

	// Then
	if result.exitCode != 0 {
		t.Fatalf("Then exit code = %d, want fail-open 0", result.exitCode)
	}
	if result.stdout != "" {
		t.Fatalf("Then stdout = %q, want empty", result.stdout)
	}
	if result.stderr == "" {
		t.Fatal("Then stderr must contain a concise diagnostic")
	}
}

func Test_Run_session_start_fails_open_for_unknown_source(t *testing.T) {
	// Given
	fixturePath := filepath.Join("session-start", "unknown-source.json")

	// When
	result := runFixture(t, "session-start", fixturePath)

	// Then
	if result.exitCode != 0 {
		t.Fatalf("Then exit code = %d, want fail-open 0", result.exitCode)
	}
	if result.stdout != "" {
		t.Fatalf("Then stdout = %q, want empty", result.stdout)
	}
	if result.stderr == "" {
		t.Fatal("Then stderr must contain a concise diagnostic")
	}
}
