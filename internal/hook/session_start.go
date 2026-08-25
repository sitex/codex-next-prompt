package hook

import "fmt"

const sessionStartContext = "For suitable final responses, append at most one line with the exact prefix `Suggested next prompt:`. Make it the final non-empty line, write the suggestion in the user's language, and omit it when no meaningful next action exists or for an exact-output response."

type sessionStartInput struct {
	HookEventName string `json:"hook_event_name"`
	Source        string `json:"source"`
}

type sessionStartOutput struct {
	HookSpecificOutput sessionStartSpecificOutput `json:"hookSpecificOutput"`
}

type sessionStartSpecificOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext"`
}

func runSessionStart(streams Streams) error {
	var input sessionStartInput
	if err := decodeSingleObject(streams.Stdin, &input); err != nil {
		return err
	}
	if input.HookEventName != "SessionStart" {
		return fmt.Errorf("session-start: unexpected hook event %q", input.HookEventName)
	}

	switch input.Source {
	case "startup", "resume", "clear", "compact":
	default:
		return fmt.Errorf("session-start: unsupported source %q", input.Source)
	}

	output := sessionStartOutput{
		HookSpecificOutput: sessionStartSpecificOutput{
			HookEventName:     "SessionStart",
			AdditionalContext: sessionStartContext,
		},
	}
	return writeJSON(streams.Stdout, output)
}
