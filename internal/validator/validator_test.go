package validator_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/yourorg/envcascade/internal/validator"
)

func TestValidate_RequiredKeyPresent(t *testing.T) {
	env := map[string]string{"APP_ENV": "production"}
	schema := validator.Schema{
		"APP_ENV": {Required: true},
	}
	if err := validator.Validate(env, schema); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_RequiredKeyMissing(t *testing.T) {
	env := map[string]string{}
	schema := validator.Schema{
		"DATABASE_URL": {Required: true},
	}
	err := validator.Validate(env, schema)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("expected error to mention DATABASE_URL, got: %v", err)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	schema := validator.Schema{
		"PORT": {Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
	}
	if err := validator.Validate(env, schema); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	env := map[string]string{"PORT": "not-a-port"}
	schema := validator.Schema{
		"PORT": {Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
	}
	err := validator.Validate(env, schema)
	if err == nil {
		t.Fatal("expected pattern mismatch error")
	}
	if !strings.Contains(err.Error(), "PORT") {
		t.Errorf("expected error to mention PORT, got: %v", err)
	}
}

func TestValidate_AllowedValues_Valid(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "info"}
	schema := validator.Schema{
		"LOG_LEVEL": {AllowedValues: []string{"debug", "info", "warn", "error"}},
	}
	if err := validator.Validate(env, schema); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_AllowedValues_Invalid(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "verbose"}
	schema := validator.Schema{
		"LOG_LEVEL": {AllowedValues: []string{"debug", "info", "warn", "error"}},
	}
	err := validator.Validate(env, schema)
	if err == nil {
		t.Fatal("expected allowed-values violation error")
	}
}

func TestValidate_OptionalKeyAbsent_NoError(t *testing.T) {
	env := map[string]string{}
	schema := validator.Schema{
		"OPTIONAL_KEY": {Pattern: regexp.MustCompile(`^[a-z]+$`)},
	}
	if err := validator.Validate(env, schema); err != nil {
		t.Fatalf("expected no error for absent optional key, got: %v", err)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "bad"}
	schema := validator.Schema{
		"DATABASE_URL": {Required: true},
		"LOG_LEVEL":    {AllowedValues: []string{"info", "debug"}},
	}
	err := validator.Validate(env, schema)
	if err == nil {
		t.Fatal("expected multiple violations")
	}
	ve, ok := err.(*validator.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d: %v", len(ve.Violations), ve.Violations)
	}
}
