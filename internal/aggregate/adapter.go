package aggregate

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// AdapterEntry is a filter-compatible entry used for aggregate input/output.
type AdapterEntry = filter.Entry

// ParseMetric validates and normalises a metric name string.
func ParseMetric(metric string) (string, error) {
	m := strings.ToLower(strings.TrimSpace(metric))
	switch m {
	case "", "count":
		return "count", nil
	case "sum", "avg", "min", "max":
		return m, nil
	}
	return "", fmt.Errorf("unsupported metric %q: must be one of count, sum, avg, min, max", metric)
}

// FromFilterEntries converts filter.Entry slice to the aggregate Entry slice.
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

// ToFilterEntries converts aggregate Entry slice back to filter.Entry slice.
func ToFilterEntries(entries []Entry) []filter.Entry {
	out := make([]filter.Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		out = append(out, filter.Entry{
			Timestamp: e.Timestamp,
			Fields:    fields,
		})
	}
	return out
}
