package pipeline

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/typecast"
)

// applyTypecast parses typecast rule expressions and applies them to entries.
// If exprs is empty, entries are returned unchanged.
func applyTypecast(entries []filter.Entry, exprs []string) ([]filter.Entry, error) {
	if len(exprs) == 0 {
		return entries, nil
	}
	opts, err := typecast.ParseConfig(exprs)
	if err != nil {
		return nil, fmt.Errorf("typecast: %w", err)
	}
	return typecast.RunFromConfig(entries, opts), nil
}
