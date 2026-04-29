package pipeline

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/unique"
)

// applyUnique extracts unique values for the given fields from entries and
// replaces the entry slice with synthetic summary entries (one per field).
// If opts is nil the entries are returned unchanged.
func applyUnique(entries []filter.Entry, opts *unique.Options) ([]filter.Entry, error) {
	if opts == nil {
		return entries, nil
	}

	ue := unique.FromFilterEntries(entries)
	result, err := unique.Run(ue, *opts)
	if err != nil {
		return nil, fmt.Errorf("unique: %w", err)
	}

	// Produce one summary entry per field listing its unique values.
	fields := make([]string, 0, len(result))
	for f := range result {
		fields = append(fields, f)
	}
	sort.Strings(fields)

	out := make([]filter.Entry, 0, len(fields))
	for _, f := range fields {
		vals := result[f]
		out = append(out, filter.Entry{
			Fields: map[string]string{
				"field":  f,
				"values": strings.Join(vals, ", "),
				"count":  fmt.Sprintf("%d", len(vals)),
			},
		})
	}
	return out, nil
}
