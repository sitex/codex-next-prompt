package hook

import "fmt"

func Run(command string, streams Streams) int {
	var err error
	switch command {
	case "session-start":
		err = runSessionStart(streams)
	case "stop":
		err = runStop(streams)
	default:
		err = fmt.Errorf("unknown command %q", command)
	}

	if err != nil {
		if _, writeErr := fmt.Fprintf(streams.Stderr, "codex-next-prompt: %v\n", err); writeErr != nil {
			return 1
		}
	}
	return 0
}
