package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcascade/internal/audit"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestAuditCascade_BasicMergeAndAudit(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=dev\nDB_HOST=localhost\n")
	override := writeTempEnv(t, "APP_ENV=staging\n")

	env, findings, err := audit.AuditCascade([]string{base, override}, audit.Options{
		RequiredKeys: []string{"APP_ENV", "DB_HOST"},
		AllowExtra:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", env["APP_ENV"])
	}
	for _, f := range findings {
		if f.Severity == audit.SeverityError {
			t.Errorf("unexpected error finding: %+v", f)
		}
	}
}

func TestAuditCascade_MissingRequiredKey(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=dev\n")

	_, findings, err := audit.AuditCascade([]string{base}, audit.Options{
		RequiredKeys: []string{"APP_ENV", "DB_PASSWORD"},
		AllowExtra:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !audit.HasErrors(findings) {
		t.Error("expected HasErrors=true for missing DB_PASSWORD")
	}
}

func TestAuditCascade_NoFiles(t *testing.T) {
	_, _, err := audit.AuditCascade(nil, audit.Options{})
	if err == nil {
		t.Error("expected error for no files")
	}
}

func TestAuditCascade_MissingFile(t *testing.T) {
	_, _, err := audit.AuditCascade([]string{"/nonexistent/.env"}, audit.Options{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestHasErrors_EmptyFindings(t *testing.T) {
	if audit.HasErrors(nil) {
		t.Error("HasErrors should be false for nil findings")
	}
}
