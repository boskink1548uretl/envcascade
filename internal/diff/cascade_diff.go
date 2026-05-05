package diff

import (
	"fmt"

	"github.com/yourorg/envcascade/internal/loader"
)

// LayerDiff represents the diff between two adjacent environment layers.
type LayerDiff struct {
	FromFile string
	ToFile   string
	Result   Result
}

// CompareFiles loads two .env files and returns their diff.
func CompareFiles(fromPath, toPath string) (LayerDiff, error) {
	base, err := loader.LoadFile(fromPath)
	if err != nil {
		return LayerDiff{}, fmt.Errorf("loading base file %q: %w", fromPath, err)
	}
	override, err := loader.LoadFile(toPath)
	if err != nil {
		return LayerDiff{}, fmt.Errorf("loading override file %q: %w", toPath, err)
	}
	return LayerDiff{
		FromFile: fromPath,
		ToFile:   toPath,
		Result:   Compare(base, override),
	}, nil
}

// CompareChain diffs each consecutive pair in a list of .env file paths.
// e.g. [dev, staging, prod] produces diffs for dev->staging and staging->prod.
func CompareChain(paths []string) ([]LayerDiff, error) {
	if len(paths) < 2 {
		return nil, fmt.Errorf("at least two file paths are required for a chain diff")
	}
	diffs := make([]LayerDiff, 0, len(paths)-1)
	for i := 0; i < len(paths)-1; i++ {
		d, err := CompareFiles(paths[i], paths[i+1])
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, d)
	}
	return diffs, nil
}
