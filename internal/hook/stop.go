package hook

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	suggestionPrefix = "Suggested next prompt:"
	stopWarning      = "codex-next-prompt: invalid Suggested next prompt line"
)

type stopInput struct {
	HookEventName        string  `json:"hook_event_name"`
	LastAssistantMessage *string `json:"last_assistant_message"`
}

type stopOutput struct {
	SystemMessage string `json:"systemMessage"`
}

func runStop(streams Streams) error {
	var input stopInput
	if err := decodeSingleObject(streams.Stdin, &input); err != nil {
		return err
	}
	if input.HookEventName != "Stop" {
		return fmt.Errorf("stop: unexpected hook event %q", input.HookEventName)
	}
	if input.LastAssistantMessage == nil || validSuggestion(*input.LastAssistantMessage) {
		return nil
	}

	return writeJSON(streams.Stdout, stopOutput{SystemMessage: stopWarning})
}

func validSuggestion(message string) bool {
	normalized := strings.ReplaceAll(strings.ReplaceAll(message, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(normalized, "\n")
	markerIndex := -1
	markerCount := 0
	finalNonEmptyIndex := -1

	for index, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			finalNonEmptyIndex = index
		}
		if strings.HasPrefix(trimmedLine, suggestionPrefix) {
			markerIndex = index
			markerCount++
		}
	}

	if markerCount == 0 {
		return true
	}
	if markerCount != 1 || markerIndex != finalNonEmptyIndex {
		return false
	}

	markerLine := strings.TrimSpace(lines[markerIndex])
	suggestion := strings.TrimSpace(strings.TrimPrefix(markerLine, suggestionPrefix))
	return suggestion != "" && utf8.RuneCountInString(markerLine) <= 240
}
