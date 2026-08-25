package hook

import "fmt"

const sessionStartContext = `<next_prompt_rules no_tools="true" no_composer="true">For suitable final responses, append at most one final non-empty line with the exact prefix ` + "`Suggested next prompt:`" + `. Write a specific prompt the user could submit verbatim, in the user's language. This suggestion is response text only: do not execute tools or describe composer insertion, prefill, ghost text, or acceptance controls. Omit it when no meaningful next action exists, for exact-output responses, terse acknowledgements, safety refusals, or questions already waiting for user-provided data. Do not imply authorization for unrequested work.</next_prompt_rules>`

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
