package exporter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envcascade/internal/exporter"
	"github.com/yourorg/envcascade/internal/validator"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestExportCascade_BasicMergeAndExport(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=dev\nDB_HOST=localhost\n")
	override := writeTempEnv(t, "APP_ENV=prod\n")

	var buf strings.Builder
	err := exporter.ExportCascade(&buf, []string{base, override}, exporter.ExportOptions{
		Format: exporter.FormatDotEnv,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=prod") {
		t.Errorf("expected overridden APP_ENV=prod, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost from base, got:\n%s", out)
	}
}

func TestExportCascade_ValidationError(t *testing.T) {
	f := writeTempEnv(t, "APP_ENV=dev\n")
	schema := &validator.Schema{
		Rules: []validator.Rule{
			{Key: "REQUIRED_KEY", Required: true},
		},
	}
	var buf strings.Builder
	err := exporter.ExportCascade(&buf, []string{f}, exporter.ExportOptions{
		Format: exporter.FormatDotEnv,
		Schema: schema,
	})
	if err == nil {
		t.Error("expected validation error, got nil")
	}
}

func TestExportCascade_NoFiles(t *testing.T) {
	var buf strings.Builder
	err := exporter.ExportCascade(&buf, nil, exporter.ExportOptions{Format: exporter.FormatDotEnv})
	if err == nil {
		t.Error("expected error for empty file list")
	}
}

func TestExportCascade_MissingFile(t *testing.T) {
	var buf strings.Builder
	err := exporter.ExportCascade(&buf, []string{filepath.Join(t.TempDir(), "nonexistent.env")}, exporter.ExportOptions{
		Format: exporter.FormatDotEnv,
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
