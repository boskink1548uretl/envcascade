package resolver

import (
	"fmt"

	"github.com/yourorg/envcascade/internal/merger"
)

// ResolveFiles loads and merges the given env files in order, then resolves
// all variable interpolations in the merged result. If errorOnMissing is true,
// any unresolvable reference returns an error.
func ResolveFiles(errorOnMissing bool, files ...string) (map[string]string, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("resolver: no files provided")
	}

	merged, err := merger.LoadAndMerge(files...)
	if err != nil {
		return nil, fmt.Errorf("resolver: merge failed: %w", err)
	}

	resolved, err := Resolve(merged, errorOnMissing)
	if err != nil {
		return nil, fmt.Errorf("resolver: interpolation failed: %w", err)
	}

	return resolved, nil
}
