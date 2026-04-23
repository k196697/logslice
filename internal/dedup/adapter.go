package dedup

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// ParseConfig parses CLI-level dedup flags into a dedup Options struct.
// fields is a comma-separated list of field names to use as the dedup key.
// consecutive, if true, only removes consecutive duplicates.
// Returns nil if dedup is not enabled (no fields specified and not global).
func ParseConfig(fields string, consecutive bool, global bool) (*Options, error) {
	if fields == "" && !global {
		return nil, nil
	}

	opts := DefaultOptions()
	opts.ConsecutiveOnly = consecutive

	if fields != "" {
		parts := strings.Split(fields, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
			if parts[i] == "" {
				return nil, fmt.Errorf("dedup: empty field name in list %q", fields)
			}
		}
		opts.Fields = parts
	}

	return &opts, nil
}

// RunFromConfig applies deduplication to a slice of filter.Entry values using
// the provided Options. Returns the original slice if opts is nil.
func RunFromConfig(entries []filter.Entry, opts *Options) []filter.Entry {
	if opts == nil {
		return entries
	}
	return Run(entries, *opts)
}
