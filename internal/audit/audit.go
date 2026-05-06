// Package audit provides utilities for auditing merged environment variable
// sets against a reference environment, reporting missing, extra, and
// redacted-sensitive keys.
package audit

import "sort"

// Severity indicates how serious an audit finding is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single audit observation.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Options controls audit behaviour.
type Options struct {
	// RequiredKeys lists keys that must be present in the merged env.
	RequiredKeys []string
	// SensitiveKeys lists keys whose values should not be logged / exported.
	SensitiveKeys []string
	// AllowExtra permits keys not in RequiredKeys without raising a warning.
	AllowExtra bool
}

// Audit inspects env against opts and returns a (possibly empty) slice of
// Findings. A non-nil error is only returned for programming mistakes (e.g.
// nil env map).
func Audit(env map[string]string, opts Options) ([]Finding, error) {
	if env == nil {
		return nil, nil
	}

	var findings []Finding

	requiredSet := toSet(opts.RequiredKeys)
	sensitiveSet := toSet(opts.SensitiveKeys)

	// Check all required keys are present.
	for _, key := range sorted(opts.RequiredKeys) {
		if _, ok := env[key]; !ok {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "required key is missing",
				Severity: SeverityError,
			})
		}
	}

	// Walk the actual env.
	for _, key := range sorted(keys(env)) {
		if sensitiveSet[key] && env[key] != "" {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "sensitive key has a non-empty value (ensure it is not committed)",
				Severity: SeverityWarning,
			})
		}
		if !opts.AllowExtra && len(requiredSet) > 0 && !requiredSet[key] {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key is not listed in required keys",
				Severity: SeverityInfo,
			})
		}
	}

	return findings, nil
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func sorted(ss []string) []string {
	cp := append([]string(nil), ss...)
	sort.Strings(cp)
	return cp
}
