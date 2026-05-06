package exporter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envcascade/internal/exporter"
)

func TestExport_DotEnv(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"DB_HOST": "localhost",
	}
	var buf strings.Builder
	if err := exporter.Export(&buf, env, exporter.FormatDotEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", out)
	}
}

func TestExport_DotEnv_QuotesSpacedValues(t *testing.T) {
	env := map[string]string{"GREETING": "hello world"}
	var buf strings.Builder
	_ = exporter.Export(&buf, env, exporter.FormatDotEnv)
	if !strings.Contains(buf.String(), `"hello world"`) {
		t.Errorf("expected quoted value for spaced string, got: %s", buf.String())
	}
}

func TestExport_Shell(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	var buf strings.Builder
	if err := exporter.Export(&buf, env, exporter.FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(buf.String(), "export ") {
		t.Errorf("expected shell export prefix, got: %s", buf.String())
	}
}

func TestExport_JSON(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	var buf strings.Builder
	if err := exporter.Export(&buf, env, exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"KEY": "value"`) {
		t.Errorf("expected JSON key-value, got:\n%s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	env := map[string]string{"X": "1"}
	var buf strings.Builder
	err := exporter.Export(&buf, env, exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExport_EmptyMap(t *testing.T) {
	var buf strings.Builder
	if err := exporter.Export(&buf, map[string]string{}, exporter.FormatDotEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "" {
		t.Errorf("expected empty output for empty map, got: %q", buf.String())
	}
}
