package exporter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents an output format for exported env vars.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Export writes the merged env map to w in the specified format.
func Export(w io.Writer, env map[string]string, format Format) error {
	switch format {
	case FormatDotEnv:
		return exportDotEnv(w, env)
	case FormatExport:
		return exportShell(w, env)
	case FormatJSON:
		return exportJSON(w, env)
	default:
		return fmt.Errorf("unsupported format: %q", format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func exportDotEnv(w io.Writer, env map[string]string) error {
	for _, k := range sortedKeys(env) {
		v := env[k]
		if strings.ContainsAny(v, " \t\n#") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func exportShell(w io.Writer, env map[string]string) error {
	for _, k := range sortedKeys(env) {
		v := env[k]
		if _, err := fmt.Fprintf(w, "export %s=%q\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(w io.Writer, env map[string]string) error {
	keys := sortedKeys(env)
	if _, err := fmt.Fprintln(w, "{"); err != nil {
		return err
	}
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		if _, err := fmt.Fprintf(w, "  %q: %q%s\n", k, env[k], comma); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w, "}")
	return err
}
