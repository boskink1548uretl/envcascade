package merger_test

import (
	"testing"

	"github.com/yourorg/envcascade/internal/merger"
)

func TestMerge_SingleLayer(t *testing.T) {
	layers := []merger.Layer{
		{Name: "base", Values: map[string]string{"APP_ENV": "base", "PORT": "8080"}},
	}
	res, err := merger.Merge(layers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["APP_ENV"] != "base" {
		t.Errorf("expected APP_ENV=base, got %q", res.Env["APP_ENV"])
	}
	if res.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", res.Env["PORT"])
	}
}

func TestMerge_OverrideOrder(t *testing.T) {
	layers := []merger.Layer{
		{Name: "base", Values: map[string]string{"APP_ENV": "base", "DB_HOST": "localhost"}},
		{Name: "dev", Values: map[string]string{"APP_ENV": "dev"}},
		{Name: "prod", Values: map[string]string{"APP_ENV": "prod", "DB_HOST": "prod-db"}},
	}
	res, err := merger.Merge(layers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["APP_ENV"] != "prod" {
		t.Errorf("expected APP_ENV=prod, got %q", res.Env["APP_ENV"])
	}
	if res.Env["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST=prod-db, got %q", res.Env["DB_HOST"])
	}
}

func TestMerge_OverrideRecords(t *testing.T) {
	layers := []merger.Layer{
		{Name: "base", Values: map[string]string{"KEY": "v1"}},
		{Name: "dev", Values: map[string]string{"KEY": "v2"}},
	}
	res, err := merger.Merge(layers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	records := res.Overrides["KEY"]
	if len(records) != 2 {
		t.Fatalf("expected 2 override records, got %d", len(records))
	}
	if records[0].Layer != "base" || records[0].Value != "v1" {
		t.Errorf("unexpected first record: %+v", records[0])
	}
	if records[1].Layer != "dev" || records[1].Value != "v2" {
		t.Errorf("unexpected second record: %+v", records[1])
	}
}

func TestMerge_EmptyLayersError(t *testing.T) {
	_, err := merger.Merge(nil)
	if err == nil {
		t.Error("expected error for empty layers, got nil")
	}
}

func TestMergeInto_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}
	out := merger.MergeInto(dst, src)
	if dst["B"] != "2" {
		t.Error("MergeInto mutated dst")
	}
	if out["B"] != "99" {
		t.Errorf("expected out[B]=99, got %q", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected out[C]=3, got %q", out["C"])
	}
}
