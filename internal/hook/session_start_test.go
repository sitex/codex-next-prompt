package hook

import (
	"encoding/json"
	"encoding/xml"
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

type nextPromptRules struct {
	XMLName    xml.Name `xml:"next_prompt_rules"`
	NoTools    string   `xml:"no_tools,attr"`
	NoComposer string   `xml:"no_composer,attr"`
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
			var rules nextPromptRules
			if err := xml.Unmarshal([]byte(envelope.HookSpecificOutput.AdditionalContext), &rules); err != nil {
				t.Fatalf("Then additionalContext must be structured XML: %v", err)
			}
			if rules.NoTools != "true" || rules.NoComposer != "true" {
				t.Fatalf("Then rule attributes = no_tools:%q no_composer:%q, want true", rules.NoTools, rules.NoComposer)
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
