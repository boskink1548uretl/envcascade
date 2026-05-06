package audit

import (
	"fmt"
	"io"
)

// Fprint writes a human-readable audit report to w.
// Each finding is printed with its severity, key, and message.
// Returns the number of bytes written and any write error.
func Fprint(w io.Writer, findings []Finding) (int, error) {
	if len(findings) == 0 {
		return fmt.Fprintln(w, "audit: no findings — all checks passed")
	}

	total := 0
	for _, f := range findings {
		icon := severityIcon(f.Severity)
		n, err := fmt.Fprintf(w, "%s [%s] %s: %s\n", icon, f.Severity, f.Key, f.Message)
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// Summary prints a one-line summary (counts per severity) to w.
func Summary(w io.Writer, findings []Finding) (int, error) {
	counts := map[Severity]int{}
	for _, f := range findings {
		counts[f.Severity]++
	}
	return fmt.Fprintf(w, "audit summary: %d error(s), %d warning(s), %d info(s)\n",
		counts[SeverityError],
		counts[SeverityWarning],
		counts[SeverityInfo],
	)
}

func severityIcon(s Severity) string {
	switch s {
	case SeverityError:
		return "✗"
	case SeverityWarning:
		return "!"
	default:
		return "i"
	}
}
