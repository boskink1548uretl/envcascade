package diff

import "sort"

// EntryStatus describes how a key differs between two env maps.
type EntryStatus string

const (
	Added    EntryStatus = "added"
	Removed  EntryStatus = "removed"
	Changed  EntryStatus = "changed"
	Unchanged EntryStatus = "unchanged"
)

// Entry represents a single key comparison result.
type Entry struct {
	Key      string
	Status   EntryStatus
	OldValue string
	NewValue string
}

// Result holds the full diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if any entry is not Unchanged.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Status != Unchanged {
			return true
		}
	}
	return false
}

// Compare produces a diff between a base env map and an override env map.
// Keys present only in base are Removed; only in override are Added;
// in both with different values are Changed; same value are Unchanged.
func Compare(base, override map[string]string) Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, oldVal := range base {
		seen[k] = true
		if newVal, ok := override[k]; ok {
			if newVal == oldVal {
				entries = append(entries, Entry{Key: k, Status: Unchanged, OldValue: oldVal, NewValue: newVal})
			} else {
				entries = append(entries, Entry{Key: k, Status: Changed, OldValue: oldVal, NewValue: newVal})
			}
		} else {
			entries = append(entries, Entry{Key: k, Status: Removed, OldValue: oldVal})
		}
	}

	for k, newVal := range override {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Status: Added, NewValue: newVal})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Result{Entries: entries}
}
