package merger_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcascade/internal/merger"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadAndMerge_BasicCascade(t *testing.T) {
	baseFile := writeTempEnv(t, "APP_ENV=base\nPORT=8080\nDB_HOST=localhost\n")
	devFile := writeTempEnv(t, "APP_ENV=dev\nDEBUG=true\n")

	cfg := merger.CascadeConfig{
		Layers: []merger.CascadeLayer{
			{Name: "base", Path: baseFile},
			{Name: "dev", Path: devFile},
		},
	}

	res, err := merger.LoadAndMerge(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["APP_ENV"] != "dev" {
		t.Errorf("expected APP_ENV=dev, got %q", res.Env["APP_ENV"])
	}
	if res.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", res.Env["PORT"])
	}
	if res.Env["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", res.Env["DEBUG"])
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", res.Env["DB_HOST"])
	}
}

func TestLoadAndMerge_MissingFileError(t *testing.T) {
	cfg := merger.CascadeConfig{
		Layers: []merger.CascadeLayer{
			{Name: "base", Path: filepath.Join(t.TempDir(), "nonexistent.env")},
		},
	}
	_, err := merger.LoadAndMerge(cfg)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadAndMerge_NoLayersError(t *testing.T) {
	_, err := merger.LoadAndMerge(merger.CascadeConfig{})
	if err == nil {
		t.Error("expected error for empty config, got nil")
	}
}

func TestLoadAndMerge_EmptyFileIsValid(t *testing.T) {
	baseFile := writeTempEnv(t, "APP_ENV=production\n")
	emptyFile := writeTempEnv(t, "")

	cfg := merger.CascadeConfig{
		Layers: []merger.CascadeLayer{
			{Name: "base", Path: baseFile},
			{Name: "empty", Path: emptyFile},
		},
	}

	res, err := merger.LoadAndMerge(cfg)
	if err != nil {
		t.Fatalf("unexpected error for empty layer file: %v", err)
	}
	if res.Env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", res.Env["APP_ENV"])
	}
}
