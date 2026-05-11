package redactor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcascade/internal/redactor"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestRedactFiles_BasicRedaction(t *testing.T) {
	base := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret\nAPI_KEY=abc123\n")
	override := writeTempEnv(t, "DB_PASSWORD=newsecret\nDEBUG=true\n")

	result, err := redactor.RedactFiles(redactor.Options{}, base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", result["APP_NAME"])
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY to be redacted, got %q", result["API_KEY"])
	}
	if result["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", result["DEBUG"])
	}
}

func TestRedactFiles_CustomPlaceholder(t *testing.T) {
	base := writeTempEnv(t, "SECRET_TOKEN=topsecret\nHOST=localhost\n")

	opts := redactor.Options{Placeholder: "***"}
	result, err := redactor.RedactFiles(opts, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["SECRET_TOKEN"] != "***" {
		t.Errorf("expected SECRET_TOKEN=***, got %q", result["SECRET_TOKEN"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", result["HOST"])
	}
}

func TestRedactFiles_NoFiles(t *testing.T) {
	_, err := redactor.RedactFiles(redactor.Options{})
	if err == nil {
		t.Error("expected error for no files, got nil")
	}
}

func TestRedactFiles_MissingFile(t *testing.T) {
	_, err := redactor.RedactFiles(redactor.Options{}, "/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRedactFiles_OverrideOrderRespected(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=development\nDB_HOST=localhost\n")
	override := writeTempEnv(t, "APP_ENV=production\n")

	result, err := redactor.RedactFiles(redactor.Options{}, base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production after override, got %q", result["APP_ENV"])
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
}
