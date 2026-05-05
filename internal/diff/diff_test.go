package diff_test

import (
	"testing"

	"github.com/yourorg/envcascade/internal/diff"
)

func TestCompare_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	override := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := diff.Compare(base, override)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	found := findEntry(result, "NEW_KEY")
	if found == nil || found.Status != diff.Added {
		t.Errorf("expected NEW_KEY to be Added, got %+v", found)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "gone"}
	override := map[string]string{"FOO": "bar"}

	result := diff.Compare(base, override)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	found := findEntry(result, "OLD")
	if found == nil || found.Status != diff.Removed {
		t.Errorf("expected OLD to be Removed, got %+v", found)
	}
}

func TestCompare_Changed(t *testing.T) {
	base := map[string]string{"DB_HOST": "localhost"}
	override := map[string]string{"DB_HOST": "prod.db.internal"}

	result := diff.Compare(base, override)

	found := findEntry(result, "DB_HOST")
	if found == nil || found.Status != diff.Changed {
		t.Errorf("expected DB_HOST to be Changed")
	}
	if found.OldValue != "localhost" || found.NewValue != "prod.db.internal" {
		t.Errorf("unexpected values: %+v", found)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	base := map[string]string{"PORT": "8080"}
	override := map[string]string{"PORT": "8080"}

	result := diff.Compare(base, override)

	if result.HasChanges() {
		t.Fatal("expected no changes")
	}
	found := findEntry(result, "PORT")
	if found == nil || found.Status != diff.Unchanged {
		t.Errorf("expected PORT to be Unchanged")
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := diff.Compare(map[string]string{}, map[string]string{})
	if result.HasChanges() {
		t.Fatal("expected no changes for empty maps")
	}
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2", "M": "3"}
	override := map[string]string{"Z": "1", "A": "99", "M": "3"}

	result := diff.Compare(base, override)

	for i := 1; i < len(result.Entries); i++ {
		if result.Entries[i-1].Key > result.Entries[i].Key {
			t.Errorf("entries not sorted: %s > %s", result.Entries[i-1].Key, result.Entries[i].Key)
		}
	}
}

func findEntry(r diff.Result, key string) *diff.Entry {
	for i := range r.Entries {
		if r.Entries[i].Key == key {
			return &r.Entries[i]
		}
	}
	return nil
}
