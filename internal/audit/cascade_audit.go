package audit

import (
	"fmt"

	"github.com/yourorg/envcascade/internal/merger"
)

// AuditCascade loads and merges the provided .env file paths in order (later
// files override earlier ones) and then runs Audit against the merged result.
// It returns both the merged environment and any findings so callers can
// display or act on them independently.
func AuditCascade(files []string, opts Options) (map[string]string, []Finding, error) {
	if len(files) == 0 {
		return nil, nil, fmt.Errorf("audit: no files provided")
	}

	merged, err := merger.LoadAndMerge(files)
	if err != nil {
		return nil, nil, fmt.Errorf("audit: loading files: %w", err)
	}

	findings, err := Audit(merged, opts)
	if err != nil {
		return merged, nil, fmt.Errorf("audit: running audit: %w", err)
	}

	return merged, findings, nil
}

// HasErrors returns true if any of the findings carry SeverityError.
func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == SeverityError {
			return true
		}
	}
	return false
}
