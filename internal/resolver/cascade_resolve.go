package resolver

import (
	"fmt"

	"github.com/user/envcascade/internal/merger"
)

// ResolveFiles loads and merges the given env files in order (later files
// override earlier ones) and then performs variable interpolation on the
// merged result. It is the primary entry point for cascaded resolution.
func ResolveFiles(paths []string, opts Options) (map[string]string, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("resolver: at least one file path is required")
	}

	merged, err := merger.LoadAndMerge(paths)
	if err != nil {
		return nil, fmt.Errorf("resolver: merge failed: %w", err)
	}

	resolved, err := Resolve(merged, opts)
	if err != nil {
		return nil, fmt.Errorf("resolver: interpolation failed: %w", err)
	}

	return resolved, nil
}
