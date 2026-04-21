package sampling

import (
	"math/rand"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

// Options controls how log entry sampling is performed.
type Options struct {
	// Rate is a value between 0.0 and 1.0 indicating the fraction of entries to keep.
	Rate float64
	// Every keeps every Nth entry (e.g. Every=10 keeps 1 in 10).
	Every int
	// Seed for reproducible sampling; 0 means use a random seed.
	Seed int64
}

// DefaultOptions returns sampling options that keep all entries.
func DefaultOptions() Options {
	return Options{
		Rate:  1.0,
		Every: 0,
		Seed:  0,
	}
}

// Run applies sampling to the provided entries and returns the sampled subset.
func Run(entries []filter.Entry, opts Options) []filter.Entry {
	if len(entries) == 0 {
		return entries
	}

	if opts.Every > 1 {
		return sampleEveryN(entries, opts.Every)
	}

	if opts.Rate >= 1.0 {
		return entries
	}

	if opts.Rate <= 0.0 {
		return []filter.Entry{}
	}

	seed := opts.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed)) //nolint:gosec

	out := make([]filter.Entry, 0, int(float64(len(entries))*opts.Rate)+1)
	for _, e := range entries {
		if rng.Float64() < opts.Rate {
			out = append(out, e)
		}
	}
	return out
}

func sampleEveryN(entries []filter.Entry, n int) []filter.Entry {
	out := make([]filter.Entry, 0, len(entries)/n+1)
	for i, e := range entries {
		if i%n == 0 {
			out = append(out, e)
		}
	}
	return out
}
