package diff

import (
	"fmt"
	"io"
	"strings"
)

// Format controls the output style of a diff result.
type Format int

const (
	FormatText Format = iota
	FormatUnified
)

// Fprint writes a human-readable diff result to w.
// FormatText uses +/- prefix lines; FormatUnified uses a compact unified-style.
func Fprint(w io.Writer, r Result, format Format) error {
	switch format {
	case FormatUnified:
		return fprintUnified(w, r)
	default:
		return fprintText(w, r)
	}
}

func fprintText(w io.Writer, r Result) error {
	for _, e := range r.Entries {
		var line string
		switch e.Status {
		case Added:
			line = fmt.Sprintf("+ %s=%s\n", e.Key, e.NewValue)
		case Removed:
			line = fmt.Sprintf("- %s=%s\n", e.Key, e.OldValue)
		case Changed:
			line = fmt.Sprintf("~ %s: %q -> %q\n", e.Key, e.OldValue, e.NewValue)
		case Unchanged:
			line = fmt.Sprintf("  %s=%s\n", e.Key, e.NewValue)
		}
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

func fprintUnified(w io.Writer, r Result) error {
	var sb strings.Builder
	for _, e := range r.Entries {
		switch e.Status {
		case Added:
			sb.WriteString(fmt.Sprintf("+%s=%s\n", e.Key, e.NewValue))
		case Removed:
			sb.WriteString(fmt.Sprintf("-%s=%s\n", e.Key, e.OldValue))
		case Changed:
			sb.WriteString(fmt.Sprintf("-%s=%s\n", e.Key, e.OldValue))
			sb.WriteString(fmt.Sprintf("+%s=%s\n", e.Key, e.NewValue))
		case Unchanged:
			sb.WriteString(fmt.Sprintf(" %s=%s\n", e.Key, e.NewValue))
		}
	}
	_, err := io.WriteString(w, sb.String())
	return err
}
