package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Required bool
	Pattern  *regexp.Regexp
	AllowedValues []string
}

// Schema maps variable names to their validation rules.
type Schema map[string]Rule

// ValidationError holds all violations found during validation.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d violation(s):\n  - %s",
		len(e.Violations), strings.Join(e.Violations, "\n  - "))
}

// Validate checks the provided env map against the given schema.
// It returns a *ValidationError if any violations are found, or nil on success.
func Validate(env map[string]string, schema Schema) error {
	var violations []string

	for key, rule := range schema {
		val, exists := env[key]

		if rule.Required && !exists {
			violations = append(violations, fmt.Sprintf("required key %q is missing", key))
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != nil && !rule.Pattern.MatchString(val) {
			violations = append(violations,
				fmt.Sprintf("key %q value %q does not match pattern %q", key, val, rule.Pattern.String()))
		}

		if len(rule.AllowedValues) > 0 {
			allowed := false
			for _, av := range rule.AllowedValues {
				if val == av {
					allowed = true
					break
				}
			}
			if !allowed {
				violations = append(violations,
					fmt.Sprintf("key %q value %q is not in allowed values %v", key, val, rule.AllowedValues))
			}
		}
	}

	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}
