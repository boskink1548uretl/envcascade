package resolver_test

import (
	"testing"

	"github.com/user/envcascade/internal/resolver"
)

func TestResolve_NoInterpolation(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out, err := resolver.Resolve(env, resolver.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "5432" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestResolve_BraceInterpolation(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "8080",
	}
	out, err := resolver.Resolve(env, resolver.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BASE_URL"] != "http://example.com:8080" {
		t.Errorf("got %q", out["BASE_URL"])
	}
}

func TestResolve_DollarInterpolation(t *testing.T) {
	env := map[string]string{
		"GREETING": "Hello $NAME",
		"NAME":     "World",
	}
	out, err := resolver.Resolve(env, resolver.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "Hello World" {
		t.Errorf("got %q", out["GREETING"])
	}
}

func TestResolve_MissingVar_LeavesPlaceholder(t *testing.T) {
	env := map[string]string{"URL": "http://${MISSING_HOST}"}
	out, err := resolver.Resolve(env, resolver.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://${MISSING_HOST}" {
		t.Errorf("expected placeholder intact, got %q", out["URL"])
	}
}

func TestResolve_ErrorOnMissing(t *testing.T) {
	env := map[string]string{"URL": "http://${MISSING_HOST}"}
	_, err := resolver.Resolve(env, resolver.Options{ErrorOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{
		"DSN": "postgres://${USER}@${HOST}/db",
		"USER": "admin",
		"HOST": "db.local",
	}
	copy := map[string]string{}
	for k, v := range original {
		copy[k] = v
	}
	_, err := resolver.Resolve(original, resolver.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range copy {
		if original[k] != v {
			t.Errorf("input mutated at key %q", k)
		}
	}
}

func TestResolve_FallbackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	env := map[string]string{"COMBINED": "prefix-${OS_VAR}"}
	out, err := resolver.Resolve(env, resolver.Options{FallbackToOS: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["COMBINED"] != "prefix-from-os" {
		t.Errorf("got %q", out["COMBINED"])
	}
}
