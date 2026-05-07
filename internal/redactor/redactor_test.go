package redactor_test

import (
	"testing"

	"github.com/yourorg/envcascade/internal/redactor"
)

func TestRedact_SensitiveKeyIsRedacted(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_HOST":    "localhost",
	}
	out := redactor.Redact(env, redactor.Options{})
	if out["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected redacted, got %q", out["DB_PASSWORD"])
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", out["APP_HOST"])
	}
}

func TestRedact_CaseInsensitiveMatch(t *testing.T) {
	env := map[string]string{"db_password": "secret"}
	out := redactor.Redact(env, redactor.Options{})
	if out["db_password"] != "***REDACTED***" {
		t.Errorf("expected redacted, got %q", out["db_password"])
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"API_KEY": "abc123"}
	out := redactor.Redact(env, redactor.Options{Placeholder: "<hidden>"})
	if out["API_KEY"] != "<hidden>" {
		t.Errorf("expected <hidden>, got %q", out["API_KEY"])
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	env := map[string]string{
		"STRIPE_KEY": "sk_live_123",
		"API_KEY":    "abc",
	}
	out := redactor.Redact(env, redactor.Options{
		SensitivePatterns: []string{"STRIPE"},
	})
	if out["STRIPE_KEY"] != "***REDACTED***" {
		t.Errorf("expected redacted, got %q", out["STRIPE_KEY"])
	}
	// API_KEY should NOT be redacted with custom patterns that don't include it
	if out["API_KEY"] != "abc" {
		t.Errorf("expected abc, got %q", out["API_KEY"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "original"}
	_ = redactor.Redact(env, redactor.Options{})
	if env["SECRET_KEY"] != "original" {
		t.Error("Redact mutated the input map")
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	out := redactor.Redact(map[string]string{}, redactor.Options{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"AUTH_TOKEN", true},
		{"APP_NAME", false},
		{"PRIVATE_KEY", true},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := redactor.IsSensitive(tc.key, redactor.DefaultSensitivePatterns)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
