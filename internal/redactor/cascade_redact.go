package redactor

import (
	"fmt"

	"github.com/yourorg/envcascade/internal/merger"
)

// RedactFiles loads and merges the given .env files in order (later files
// override earlier ones), then redacts any sensitive keys from the merged
// result. At least one file path must be provided.
func RedactFiles(opts Options, files ...string) (map[string]string, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("redactor: at least one file path is required")
	}

	merged, err := merger.LoadAndMerge(files...)
	if err != nil {
		return nil, fmt.Errorf("redactor: failed to load and merge files: %w", err)
	}

	return Redact(merged, opts), nil
}
