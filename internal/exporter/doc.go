// Package exporter provides utilities for serialising a merged env map into
// multiple output formats.
//
// Supported formats:
//
//	- dotenv  — KEY=VALUE pairs, values containing whitespace or special
//	            characters are double-quoted.
//	- export  — POSIX shell "export KEY=VALUE" statements suitable for
//	            sourcing in bash/zsh.
//	- json    — A JSON object mapping string keys to string values.
//
// The high-level ExportCascade helper combines the merger and validator
// packages so callers can load, merge, validate, and export in a single call:
//
//	err := exporter.ExportCascade(os.Stdout, []string{".env", ".env.prod"}, exporter.ExportOptions{
//	    Format: exporter.FormatExport,
//	    Schema: mySchema,
//	})
package exporter
