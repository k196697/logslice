package transform

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
)

// Config holds the transform configuration passed from the CLI or pipeline.
type Config struct {
	// Exprs is a list of raw transform expressions, e.g. ["rename:level=severity", "drop:debug"].
	Exprs []string
}

// ParseConfig parses all expressions in cfg and returns the resulting rules.
// Returns an error if any expression is invalid.
func ParseConfig(cfg Config) ([]Rule, error) {
	rules := make([]Rule, 0, len(cfg.Exprs))
	for _, expr := range cfg.Exprs {
		r, err := ParseRule(expr)
		if err != nil {
			return nil, fmt.Errorf("transform: %w", err)
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// RunFromConfig is a convenience wrapper that parses cfg and applies all rules
// to entries. Returns an error if any expression cannot be parsed.
func RunFromConfig(entries []filter.Entry, cfg Config) ([]filter.Entry, error) {
	if len(cfg.Exprs) == 0 {
		return entries, nil
	}
	rules, err := ParseConfig(cfg)
	if err != nil {
		return nil, err
	}
	return Run(entries, rules), nil
}
