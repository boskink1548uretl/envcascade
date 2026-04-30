package merger

import (
	"fmt"

	"github.com/yourorg/envcascade/internal/loader"
)

// CascadeConfig describes the ordered list of env files to merge.
type CascadeConfig struct {
	// Layers is an ordered slice of (name, filepath) pairs.
	// Earlier entries are base layers; later entries override them.
	Layers []CascadeLayer
}

// CascadeLayer pairs a human-readable name with a file path.
type CascadeLayer struct {
	Name string
	Path string
}

// LoadAndMerge loads each layer file in order and merges them.
// Files that do not exist are skipped if Optional is true on the layer;
// otherwise a missing file is a hard error.
func LoadAndMerge(cfg CascadeConfig) (*MergeResult, error) {
	if len(cfg.Layers) == 0 {
		return nil, fmt.Errorf("cascade: no layers configured")
	}

	var layers []Layer
	for _, cl := range cfg.Layers {
		vals, err := loader.LoadFile(cl.Path)
		if err != nil {
			return nil, fmt.Errorf("cascade: loading layer %q from %q: %w", cl.Name, cl.Path, err)
		}
		layers = append(layers, Layer{
			Name:   cl.Name,
			Values: vals,
		})
	}

	return Merge(layers)
}
