package exporter

import (
	"fmt"
	"io"

	"github.com/yourorg/envcascade/internal/merger"
	"github.com/yourorg/envcascade/internal/validator"
)

// ExportOptions controls the behaviour of ExportCascade.
type ExportOptions struct {
	// Schema is optional; when non-nil the merged env is validated before export.
	Schema *validator.Schema
	Format Format
}

// ExportCascade loads and merges the given env files in order (lowest to
// highest precedence), optionally validates the result, then writes it to w.
func ExportCascade(w io.Writer, files []string, opts ExportOptions) error {
	if len(files) == 0 {
		return fmt.Errorf("exporter: at least one file is required")
	}

	env, err := merger.LoadAndMerge(files)
	if err != nil {
		return fmt.Errorf("exporter: merge failed: %w", err)
	}

	if opts.Schema != nil {
		if verrs := validator.Validate(env, opts.Schema); len(verrs) > 0 {
			return fmt.Errorf("exporter: validation failed (%d error(s)): %v", len(verrs), verrs)
		}
	}

	if err := Export(w, env, opts.Format); err != nil {
		return fmt.Errorf("exporter: write failed: %w", err)
	}
	return nil
}
