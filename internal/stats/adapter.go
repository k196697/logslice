package stats

import (
	"github.com/user/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to stats.Entry slice
// so that stats can be computed without a direct dependency on filter internals.
func FromFilterEntries(entries []filter.Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		out = append(out, Entry{
			Timestamp: e.Timestamp,
			Fields:    fields,
		})
	}
	return out
}
