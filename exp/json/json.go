package json

import (
	"encoding/json"
	"fmt"
	"io"
)

// From reads from a reader and unmarshals the response into the provided type.
func From[T any](r io.Reader, t T) (T, error) {
	bts, err := io.ReadAll(r)
	if err != nil {
		return t, fmt.Errorf("failed to read response: %w", err)
	}
	if err := json.Unmarshal(bts, &t); err != nil {
		return t, fmt.Errorf("failed to parse body: %w: %s", err, bts)
	}
	return t, nil
}

// To marshals the provide type to a writer.
func To[T any](w io.Writer, t T) error {
	bts, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	if _, err := w.Write(bts); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
