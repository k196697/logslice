package pipeline

import (
	"github.com/yourorg/logslice/internal/dedup"
	"github.com/yourorg/logslice/internal/filter"
)

// applyDedup removes duplicate entries from the pipeline based on the provided options.
// If opts is nil, deduplication is skipped and the original entries are returned.
func applyDedup(entries []filter.Entry, opts *dedup.Options) []filter.Entry {
	if opts == nil {
		return entries
	}

	result := dedup.Run(entries, *opts)
	return result
}
