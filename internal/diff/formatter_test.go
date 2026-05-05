package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envcascade/internal/diff"
)

func TestFprintText_Added(t *testing.T) {
	r := diff.Result{
		Entries: []diff.Entry{
			{Key: "NEW", Status: diff.Added, NewValue: "yes"},
		},
	}
	var buf bytes.Buffer
	if err := diff.Fprint(&buf, r, diff.FormatText); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "+ NEW=yes") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestFprintText_Removed(t *testing.T) {
	r := diff.Result{
		Entries: []diff.Entry{
			{Key: "OLD", Status: diff.Removed, OldValue: "gone"},
		},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, r, diff.FormatText)
	if !strings.Contains(buf.String(), "- OLD=gone") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestFprintText_Changed(t *testing.T) {
	r := diff.Result{
		Entries: []diff.Entry{
			{Key: "HOST", Status: diff.Changed, OldValue: "local", NewValue: "remote"},
		},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, r, diff.FormatText)
	out := buf.String()
	if !strings.Contains(out, "~ HOST") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFprintUnified_ChangedShowsBothLines(t *testing.T) {
	r := diff.Result{
		Entries: []diff.Entry{
			{Key: "PORT", Status: diff.Changed, OldValue: "3000", NewValue: "8080"},
		},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, r, diff.FormatUnified)
	out := buf.String()
	if !strings.Contains(out, "-PORT=3000") {
		t.Errorf("missing removal line: %s", out)
	}
	if !strings.Contains(out, "+PORT=8080") {
		t.Errorf("missing addition line: %s", out)
	}
}

func TestFprintUnified_Unchanged(t *testing.T) {
	r := diff.Result{
		Entries: []diff.Entry{
			{Key: "STABLE", Status: diff.Unchanged, OldValue: "v1", NewValue: "v1"},
		},
	}
	var buf bytes.Buffer
	diff.Fprint(&buf, r, diff.FormatUnified)
	out := buf.String()
	if !strings.HasPrefix(out, " STABLE=v1") {
		t.Errorf("unexpected output: %q", out)
	}
}
