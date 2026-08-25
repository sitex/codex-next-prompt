package hook

import (
	"encoding/json"
	"fmt"
	"io"
)

type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type hookInput struct {
	HookEventName string `json:"hook_event_name"`
}

func decodeSingleObject[T any](reader io.Reader, destination *T) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode input: %w", err)
	}

	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode input: multiple JSON values")
		}
		return fmt.Errorf("decode input trailing data: %w", err)
	}

	return nil
}

func writeJSON[T any](writer io.Writer, value T) error {
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("encode output: %w", err)
	}
	return nil
}
