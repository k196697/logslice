package pipeline

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/transform"
)

// applyTransform applies transformation rules to filter entries.
// It parses the raw rule strings from Config and runs them against
// each entry in the slice, returning the transformed result.
func applyTransform(entries []filter.Entry, cfg *Config) ([]filter.Entry, error) {
	if len(cfg.TransformRules) == 0 {
		return entries, nil
	}

	transformCfg, err := transform.ParseConfig(cfg.TransformRules)
	if err != nil {
		return nil, fmt.Errorf("transform: %w", err)
	}

	if transformCfg == nil {
		return entries, nil
	}

	return transform.RunFromConfig(entries, transformCfg), nil
}
