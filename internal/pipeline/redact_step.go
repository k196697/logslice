package pipeline

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/redact"
)

// applyRedact applies redaction rules to the pipeline entries.
// If no redaction config is provided, entries are returned unchanged.
func applyRedact(entries []filter.Entry, fields []string, patterns []string, mask string) ([]filter.Entry, error) {
	if len(fields) == 0 && len(patterns) == 0 {
		return entries, nil
	}

	cfg, err := redact.ParseConfig(fields, patterns, mask)
	if err != nil {
		return nil, fmt.Errorf("redact: invalid config: %w", err)
	}
	if cfg == nil {
		return entries, nil
	}

	return redact.RunFromConfig(entries, cfg)
}
