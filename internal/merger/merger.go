package merger

import (
	"fmt"
	"maps"
)

// Layer represents a named environment layer (e.g., base, dev, staging, prod).
type Layer struct {
	Name   string
	Values map[string]string
}

// MergeResult holds the final merged environment and metadata about overrides.
type MergeResult struct {
	Env       map[string]string
	Overrides map[string][]OverrideRecord
}

// OverrideRecord tracks which layer set or overrode a key.
type OverrideRecord struct {
	Layer string
	Value string
}

// Merge applies layers in order, with later layers overriding earlier ones.
// The first layer is treated as the base. Each subsequent layer overrides keys.
func Merge(layers []Layer) (*MergeResult, error) {
	if len(layers) == 0 {
		return nil, fmt.Errorf("merger: at least one layer is required")
	}

	result := &MergeResult{
		Env:       make(map[string]string),
		Overrides: make(map[string][]OverrideRecord),
	}

	for _, layer := range layers {
		for k, v := range layer.Values {
			result.Overrides[k] = append(result.Overrides[k], OverrideRecord{
				Layer: layer.Name,
				Value: v,
			})
			result.Env[k] = v
		}
	}

	return result, nil
}

// MergeInto merges src into dst without modifying src.
// Keys in src override keys in dst.
func MergeInto(dst, src map[string]string) map[string]string {
	out := make(map[string]string, len(dst))
	maps.Copy(out, dst)
	maps.Copy(out, src)
	return out
}
