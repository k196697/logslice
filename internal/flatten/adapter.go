package flatten

import (
	"github.com/user/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to flatten.Entry slice.
func FromFilterEntries(entries []filter.Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		out = append(out, Entry{
			Timestamp: e.Timestamp.UnixNano(),
			Fields:    fields,
		})
	}
	return out
}

// ToFilterEntries converts flatten.Entry slice back to filter.Entry slice.
func ToFilterEntries(entries []Entry) []filter.Entry {
	out := make([]filter.Entry, 0, len(entries))
	for _, e := range entries {
		fe := filter.Entry{
			Fields: make(map[string]string, len(e.Fields)),
		}
		for k, v := range e.Fields {
			switch sv := v.(type) {
			case string:
				fe.Fields[k] = sv
			default:
				fe.Fields[k] = formatValue(v)
			}
		}
		out = append(out, fe)
	}
	return out
}

// ParseConfig returns nil when no flattening is requested.
func ParseConfig(separator string, maxDepth int) *Options {
	if separator == "" && maxDepth == 0 {
		return nil
	}
	opts := DefaultOptions()
	if separator != "" {
		opts.Separator = separator
	}
	opts.MaxDepth = maxDepth
	return &opts
}

// RunFromConfig runs flattening only when opts is non-nil.
func RunFromConfig(entries []filter.Entry, opts *Options) []filter.Entry {
	if opts == nil {
		return entries
	}
	flat := FromFilterEntries(entries)
	result := Run(flat, *opts)
	return ToFilterEntries(result)
}
