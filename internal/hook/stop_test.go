package hook

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"
)

func Test_Run_stop_completes_silently_for_valid_or_absent_suggestion(t *testing.T) {
	// Given
	fixtures := []string{
		"valid-final.json",
		"absent.json",
		"crlf-unicode.json",
		"null-message.json",
		"missing-message.json",
		"stop-hook-active-true.json",
		"stop-hook-active-false.json",
		"unknown-fields.json",
	}

	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixturePath := filepath.Join("stop", fixtureName)

			// When
			result := runFixture(t, "stop", fixturePath)

			// Then
			if result.exitCode != 0 {
				t.Fatalf("Then exit code = %d, want 0", result.exitCode)
			}
			if result.stdout != "" {
				t.Fatalf("Then stdout = %q, want silent completion", result.stdout)
			}
			if result.stderr != "" {
				t.Fatalf("Then stderr = %q, want empty", result.stderr)
			}
		})
	}
}

func Test_StopInput_decodes_stop_hook_active_contract(t *testing.T) {
	// Given
	data := []byte(`{"hook_event_name":"Stop","stop_hook_active":true}`)

	// When
	var input stopInput
	err := json.Unmarshal(data, &input)

	// Then
	if err != nil {
		t.Fatalf("Then decode stop input: %v", err)
	}
	if !input.StopHookActive {
		t.Fatal("Then stop_hook_active must be represented as true")
	}
}

func Test_Run_stop_warns_for_invalid_present_suggestion(t *testing.T) {
	// Given
	fixtures := []string{
		"duplicate.json",
		"empty.json",
		"non-final.json",
		"over-240-unicode.json",
		"whitespace-empty.json",
	}

	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixturePath := filepath.Join("stop", fixtureName)

			// When
			result := runFixture(t, "stop", fixturePath)

			// Then
			if result.exitCode != 0 {
				t.Fatalf("Then exit code = %d, want 0", result.exitCode)
			}
			if result.stderr != "" {
				t.Fatalf("Then stderr = %q, want empty", result.stderr)
			}
			assertSingleTrailingNewline(t, result.stdout)

			object := decodeSingleJSONObject(t, result.stdout)
			var systemMessage string
			if err := json.Unmarshal(object["systemMessage"], &systemMessage); err != nil {
				t.Fatalf("Then systemMessage must be a string: %v", err)
			}
			if systemMessage != "codex-next-prompt: invalid Suggested next prompt line" {
				t.Fatalf("Then systemMessage = %q, want stable protocol warning", systemMessage)
			}
			assertNoStopContinuationFields(t, object)
		})
	}
}

func Test_Run_stop_counts_overlong_marker_in_unicode_code_points(t *testing.T) {
	// Given
	fixturePath := filepath.Join("stop", "over-240-unicode.json")
	fixture := readStopFixture(t, fixturePath)
	if fixture.LastAssistantMessage == nil {
		t.Fatal("Given overlong fixture must contain last_assistant_message")
	}
	if utf8.RuneCountInString(*fixture.LastAssistantMessage) <= 240 {
		t.Fatalf("Given marker rune count = %d, want over 240", utf8.RuneCountInString(*fixture.LastAssistantMessage))
	}

	// When
	result := runFixture(t, "stop", fixturePath)

	// Then
	if result.stdout == "" {
		t.Fatal("Then overlong Unicode marker must produce a warning")
	}
}

func Test_Run_stop_fails_open_for_wrong_event_shape_or_second_object(t *testing.T) {
	// Given
	fixtures := []string{"wrong-event.json", "second-object.json"}

	for _, fixtureName := range fixtures {
		t.Run(fixtureName, func(t *testing.T) {
			fixturePath := filepath.Join("stop", fixtureName)

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
		})
	}
}

type stopFixture struct {
	LastAssistantMessage *string `json:"last_assistant_message"`
}

func readStopFixture(t *testing.T, fixturePath string) stopFixture {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", fixturePath))
	if err != nil {
		t.Fatalf("Given fixture %q: read: %v", fixturePath, err)
	}

	var fixture stopFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatalf("Given fixture %q: decode: %v", fixturePath, err)
	}
	return fixture
}

func assertNoStopContinuationFields(t *testing.T, object map[string]json.RawMessage) {
	t.Helper()

	forbidden := []string{
		"decision",
		"reason",
		"continue",
		"stopReason",
		"hookSpecificOutput",
		"suppressOutput",
		"additionalContext",
	}
	for _, field := range forbidden {
		if _, exists := object[field]; exists {
			t.Errorf("Then Stop output must not contain %q", field)
		}
	}
	if len(object) != 1 {
		t.Errorf("Then Stop warning must contain only systemMessage, got %d fields", len(object))
	}
}
