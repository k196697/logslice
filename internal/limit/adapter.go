package limit

import "github.com/user/logslice/internal/filter"

// FromFilterEntries converts filter.Entry slice to limit.Entry slice.
func FromFilterEntries(in []filter.Entry) []Entry {
	out := make([]Entry, len(in))
	for i, e := range in {
		fields := make(map[string]string, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		out[i] = Entry{
			Fields:        fields,
			TimestampNano: e.Timestamp.UnixNano(),
		}
	}
	return out
}

// ToFilterEntries converts limit.Entry slice back to filter.Entry slice.
func ToFilterEntries(in []Entry, ref []filter.Entry) []filter.Entry {
	// Build a nano->index map from the reference slice for timestamp restoration.
	refMap := make(map[int64]filter.Entry, len(ref))
	for _, e := range ref {
		refMap[e.Timestamp.UnixNano()] = e
	}

	out := make([]filter.Entry, 0, len(in))
	for _, e := range in {
		if fe, ok := refMap[e.TimestampNano]; ok {
			out = append(out, fe)
		}
	}
	return out
}

// ParseConfig returns non-nil Options only when maxEntries > 0.
func ParseConfig(maxEntries int, fromEnd bool) *Options {
	if maxEntries <= 0 {
		return nil
	}
	opts := DefaultOptions()
	opts.MaxEntries = maxEntries
	opts.FromEnd = fromEnd
	return &opts
}

// RunFromConfig applies the limit using filter.Entry types.
// If opts is nil, entries are returned unchanged.
func RunFromConfig(entries []filter.Entry, opts *Options) []filter.Entry {
	if opts == nil {
		return entries
	}
	limited := Run(FromFilterEntries(entries), *opts)
	return ToFilterEntries(limited, entries)
}
