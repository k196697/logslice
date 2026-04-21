package pipeline

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

const tsLayout = time.RFC3339

// applyFilter builds filter criteria from Options and runs them.
func applyFilter(opts Options, entries []filter.Entry) ([]filter.Entry, error) {
	var criteria []filter.Criterion

	if opts.From != "" {
		t, err := time.Parse(tsLayout, opts.From)
		if err != nil {
			return nil, fmt.Errorf("invalid --from value %q: %w", opts.From, err)
		}
		criteria = append(criteria, filter.After(t))
	}

	if opts.To != "" {
		t, err := time.Parse(tsLayout, opts.To)
		if err != nil {
			return nil, fmt.Errorf("invalid --to value %q: %w", opts.To, err)
		}
		criteria = append(criteria, filter.Before(t))
	}

	for k, v := range opts.Fields {
		criteria = append(criteria, filter.FieldEquals(k, v))
	}

	result := filter.Apply(entries, criteria...)

	if opts.TailN > 0 && opts.TailN < len(result) {
		result = result[len(result)-opts.TailN:]
	}

	return result, nil
}
