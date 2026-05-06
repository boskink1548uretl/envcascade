package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcascade/internal/resolver"
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

func TestResolveFiles_BasicInterpolation(t *testing.T) {
	base := writeTempEnv(t, "HOST=localhost\nPORT=5432\n")
	override := writeTempEnv(t, "DSN=postgres://$HOST:$PORT/mydb\n")

	result, err := resolver.ResolveFiles(true, base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "postgres://localhost:5432/mydb"
	if result["DSN"] != want {
		t.Errorf("DSN = %q; want %q", result["DSN"], want)
	}
}

func TestResolveFiles_OverrideBeforeResolve(t *testing.T) {
	base := writeTempEnv(t, "ENV=dev\nAPP_URL=http://localhost\n")
	override := writeTempEnv(t, "ENV=prod\nFULL_URL=${APP_URL}/api\n")

	result, err := resolver.ResolveFiles(false, base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["ENV"] != "prod" {
		t.Errorf("ENV = %q; want %q", result["ENV"], "prod")
	}
	if result["FULL_URL"] != "http://localhost/api" {
		t.Errorf("FULL_URL = %q; want %q", result["FULL_URL"], "http://localhost/api")
	}
}

func TestResolveFiles_MissingVarError(t *testing.T) {
	base := writeTempEnv(t, "URL=http://${UNDEFINED_HOST}/path\n")

	_, err := resolver.ResolveFiles(true, base)
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestResolveFiles_MissingVarLeavesPlaceholder(t *testing.T) {
	base := writeTempEnv(t, "URL=http://${UNDEFINED_HOST}/path\n")

	result, err := resolver.ResolveFiles(false, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["URL"] != "http://${UNDEFINED_HOST}/path" {
		t.Errorf("URL = %q; want placeholder preserved", result["URL"])
	}
}

func TestResolveFiles_NoFiles(t *testing.T) {
	_, err := resolver.ResolveFiles(false)
	if err == nil {
		t.Fatal("expected error for no files, got nil")
	}
}

func TestResolveFiles_MissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nonexistent.env")
	_, err := resolver.ResolveFiles(false, missing)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
