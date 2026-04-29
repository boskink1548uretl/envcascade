package loader

import (
	"os"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envcascade-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFile_BasicKeyValue(t *testing.T) {
	path := writeTempEnvFile(t, "APP_ENV=production\nDEBUG=false\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", env["APP_ENV"])
	}
	if env["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", env["DEBUG"])
	}
}

func TestLoadFile_CommentsAndBlanksIgnored(t *testing.T) {
	content := "# this is a comment\n\nKEY=value\n"
	path := writeTempEnvFile(t, content)
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
}

func TestLoadFile_QuotedValues(t *testing.T) {
	path := writeTempEnvFile(t, `DB_URL="postgres://localhost/mydb"\nSECRET='mysecret'\n`)
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("expected unquoted DB_URL, got %q", env["DB_URL"])
	}
}

func TestLoadFile_MissingEquals(t *testing.T) {
	path := writeTempEnvFile(t, "BADLINE\n")
	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for line missing '=', got nil")
	}
}

func TestLoadFile_FileNotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestStripQuotes(t *testing.T) {
	cases := []struct{ input, want string }{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`noquotes`, "noquotes"},
		{`"mismatch'`, `"mismatch'`},
	}
	for _, c := range cases {
		got := stripQuotes(c.input)
		if got != c.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
