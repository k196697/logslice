package unique

import (
	"strings"

	"github.com/user/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to unique.Entry slice.
func FromFilterEntries(entries []filter.Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		ue := Entry{
			Fields: make(map[string]string, len(e.Fields)),
		}
		for k, v := range e.Fields {
			ue.Fields[k] = v
		}
		out = append(out, ue)
	}
	return out
}

// ParseConfig returns Options parsed from CLI flag values.
// fields is a comma-separated list of field names; limit is the cap (0 = none).
func ParseConfig(fields string, limit int) *Options {
	fields = strings.TrimSpace(fields)
	if fields == "" {
		return nil
	}
	parts := strings.Split(fields, ",")
	filtered := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	opts := DefaultOptions()
	opts.Fields = filtered
	if limit > 0 {
		opts.Limit = limit
	}
	return &opts
}
