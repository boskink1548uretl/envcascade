package audit_test

import (
	"testing"

	"github.com/yourorg/envcascade/internal/audit"
)

func TestAudit_RequiredKeyPresent(t *testing.T) {
	env := map[string]string{"APP_ENV": "production"}
	findings, err := audit.Audit(env, audit.Options{RequiredKeys: []string{"APP_ENV"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, f := range findings {
		if f.Severity == audit.SeverityError {
			t.Errorf("unexpected error finding: %+v", f)
		}
	}
}

func TestAudit_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{"OTHER": "val"}
	findings, err := audit.Audit(env, audit.Options{RequiredKeys: []string{"APP_ENV"}, AllowExtra: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Key == "APP_ENV" && f.Severity == audit.SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected error finding for missing required key APP_ENV")
	}
}

func TestAudit_SensitiveKeyWarning(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret"}
	findings, err := audit.Audit(env, audit.Options{
		SensitiveKeys: []string{"DB_PASSWORD"},
		AllowExtra:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Key == "DB_PASSWORD" && f.Severity == audit.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for sensitive key with non-empty value")
	}
}

func TestAudit_ExtraKeyInfo(t *testing.T) {
	env := map[string]string{"KNOWN": "v", "EXTRA": "x"}
	findings, err := audit.Audit(env, audit.Options{
		RequiredKeys: []string{"KNOWN"},
		AllowExtra:   false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Key == "EXTRA" && f.Severity == audit.SeverityInfo {
			found = true
		}
	}
	if !found {
		t.Error("expected info finding for extra key")
	}
}

func TestAudit_AllowExtraSuppressesInfo(t *testing.T) {
	env := map[string]string{"KNOWN": "v", "EXTRA": "x"}
	findings, err := audit.Audit(env, audit.Options{
		RequiredKeys: []string{"KNOWN"},
		AllowExtra:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, f := range findings {
		if f.Key == "EXTRA" {
			t.Errorf("did not expect finding for EXTRA when AllowExtra=true, got %+v", f)
		}
	}
}

func TestAudit_NilEnv(t *testing.T) {
	findings, err := audit.Audit(nil, audit.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings for nil env, got %d", len(findings))
	}
}
