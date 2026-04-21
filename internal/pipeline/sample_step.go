package pipeline

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/sampling"
)

// SampleConfig carries raw sampling parameters forwarded from CLI flags.
type SampleConfig struct {
	Rate  string
	Every string
	Seed  int64
}

// applySampling runs the sampling step when sampling is configured.
// It returns the original slice unchanged when no sampling is requested.
func applySampling(entries []filter.Entry, cfg SampleConfig) ([]filter.Entry, error) {
	if cfg.Rate == "" && cfg.Every == "" {
		return entries, nil
	}

	opts, err := sampling.ParseConfig(sampling.Config{
		RateStr:  cfg.Rate,
		EveryStr: cfg.Every,
		Seed:     cfg.Seed,
	})
	if err != nil {
		return nil, fmt.Errorf("pipeline sample step: %w", err)
	}

	return sampling.Run(entries, opts), nil
}
