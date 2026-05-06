// Package resolver provides variable interpolation for merged env maps.
// It expands references like ${VAR} or $VAR within values using the
// same merged map, supporting a single-pass topological resolution.
package resolver

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var interpolationRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls resolution behaviour.
type Options struct {
	// FallbackToOS allows falling back to os.Getenv when a key is not
	// present in the merged map.
	FallbackToOS bool
	// ErrorOnMissing returns an error when a referenced variable cannot
	// be resolved instead of leaving the placeholder intact.
	ErrorOnMissing bool
}

// Resolve performs variable interpolation on all values in env.
// It returns a new map; the original is not mutated.
func Resolve(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := expandValue(v, env, opts)
		if err != nil {
			return nil, fmt.Errorf("resolver: key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func expandValue(value string, env map[string]string, opts Options) (string, error) {
	var expandErr error
	result := interpolationRe.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := env[name]; ok {
			return v
		}
		if opts.FallbackToOS {
			if v := os.Getenv(name); v != "" {
				return v
			}
		}
		if opts.ErrorOnMissing {
			expandErr = fmt.Errorf("unresolved variable %q", name)
			return match
		}
		return match
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "${") 
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
