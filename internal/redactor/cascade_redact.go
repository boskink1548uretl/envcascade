package redactor

import (
	"github.com/yourorg/envcascade/internal/merger"
)

// RedactFiles loads and merges the given .env files in order (later files
// override earlier ones) and returns a redacted copy of the merged map.
// It is a convenience wrapper around merger.LoadAndMerge + Redact.
func RedactFiles(opts Options, files ...string) (map[string]string, error) {
	if len(files) == 0 {
		return nil, merger.ErrNoLayers
	}

	merged, err := merger.LoadAndMerge(files...)
	if err != nil {
		return nil, err
	}

	return Redact(merged, opts), nil
}
