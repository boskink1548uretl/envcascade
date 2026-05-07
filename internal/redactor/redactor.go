// Package redactor provides utilities for masking sensitive values
// in merged environment maps before display or export.
package redactor

import "strings"

// DefaultSensitivePatterns is the list of key substrings that trigger redaction.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIAL",
	"AUTH",
}

const redactedPlaceholder = "***REDACTED***"

// Options controls redaction behaviour.
type Options struct {
	// SensitivePatterns overrides DefaultSensitivePatterns when non-nil.
	SensitivePatterns []string
	// Placeholder replaces the redacted value; defaults to "***REDACTED***".
	Placeholder string
}

// Redact returns a shallow copy of env with sensitive values masked.
// Keys are matched case-insensitively against the configured patterns.
func Redact(env map[string]string, opts Options) map[string]string {
	patterns := opts.SensitivePatterns
	if patterns == nil {
		patterns = DefaultSensitivePatterns
	}
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = redactedPlaceholder
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, patterns) {
			out[k] = placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitive reports whether key matches any of the given patterns.
func IsSensitive(key string, patterns []string) bool {
	return isSensitive(key, patterns)
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
