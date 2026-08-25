package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const maxInputBytes = 4 * 1024 * 1024

type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func decodeSingleObject[T any](reader io.Reader, destination *T) error {
	input, err := io.ReadAll(io.LimitReader(reader, maxInputBytes+1))
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	if len(input) > maxInputBytes {
		return fmt.Errorf("decode input: exceeds %d-byte limit", maxInputBytes)
	}

	decoder := json.NewDecoder(bytes.NewReader(input))
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
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("encode output: %w", err)
	}
	return nil
}
