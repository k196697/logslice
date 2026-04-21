package sampling

import (
	"fmt"
	"strconv"
	"strings"
)

// Config holds raw string configuration for sampling, typically sourced from
// CLI flags.
type Config struct {
	RateStr  string
	EveryStr string
	Seed     int64
}

// ParseConfig converts a Config into Options, returning an error for invalid
// values.
func ParseConfig(cfg Config) (Options, error) {
	opts := DefaultOptions()
	opts.Seed = cfg.Seed

	if cfg.RateStr != "" {
		v, err := strconv.ParseFloat(strings.TrimSpace(cfg.RateStr), 64)
		if err != nil {
			return opts, fmt.Errorf("sampling: invalid rate %q: %w", cfg.RateStr, err)
		}
		if v < 0 || v > 1 {
			return opts, fmt.Errorf("sampling: rate must be between 0.0 and 1.0, got %f", v)
		}
		opts.Rate = v
	}

	if cfg.EveryStr != "" {
		n, err := strconv.Atoi(strings.TrimSpace(cfg.EveryStr))
		if err != nil {
			return opts, fmt.Errorf("sampling: invalid every value %q: %w", cfg.EveryStr, err)
		}
		if n < 1 {
			return opts, fmt.Errorf("sampling: every must be >= 1, got %d", n)
		}
		opts.Every = n
	}

	return opts, nil
}
