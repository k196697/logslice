package pipeline

import (
	"fmt"

	"github.com/logslice/logslice/internal/aggregate"
	"github.com/logslice/logslice/internal/filter"
)

// AggregateConfig holds the pipeline-level aggregate configuration parsed
// from CLI flags.
type AggregateConfig struct {
	// GroupBy is the field to group log entries by.
	GroupBy string
	// CountField is the output field name for the per-group count.
	CountField string
}

// ParseAggregateConfig builds an AggregateConfig from raw flag values.
// Returns nil when groupBy is empty (feature disabled).
func ParseAggregateConfig(groupBy, countField string) (*AggregateConfig, error) {
	if groupBy == "" {
		return nil, nil
	}
	cf := countField
	if cf == "" {
		cf = aggregate.DefaultOptions().CountField
	}
	return &AggregateConfig{GroupBy: groupBy, CountField: cf}, nil
}

// applyAggregate runs the aggregate step when cfg is non-nil.
func applyAggregate(entries []filter.Entry, cfg *AggregateConfig) ([]filter.Entry, error) {
	if cfg == nil {
		return entries, nil
	}
	opts := aggregate.Options{
		GroupBy:    cfg.GroupBy,
		CountField: cfg.CountField,
	}
	result, err := aggregate.Run(entries, opts)
	if err != nil {
		return nil, fmt.Errorf("aggregate step: %w", err)
	}
	return result, nil
}
