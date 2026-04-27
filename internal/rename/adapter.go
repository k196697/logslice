package rename

import (
	"github.com/user/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to rename.Entry slice.
func FromFilterEntries(entries []filter.Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]string, len(e.Fields))
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

// ToFilterEntries converts rename.Entry slice back to filter.Entry slice.
func ToFilterEntries(entries []Entry, timeField string) []filter.Entry {
	if timeField == "" {
		timeField = "time"
	}
	out := make([]filter.Entry, 0, len(entries))
	for _, e := range entries {
		fe := filter.Entry{
			Fields: make(map[string]string, len(e.Fields)),
		}
		for k, v := range e.Fields {
			fe.Fields[k] = v
		}
		out = append(out, fe)
	}
	return out
}

// ParseConfig parses rename rule expressions and returns Options, or nil if exprs is empty.
func ParseConfig(exprs []string) (*Options, error) {
	if len(exprs) == 0 {
		return nil, nil
	}
	rules, err := ParseRules(exprs)
	if err != nil {
		return nil, err
	}
	opts := Options{Rules: rules}
	return &opts, nil
}

// RunFromConfig applies rename rules to filter entries if opts is non-nil.
func RunFromConfig(entries []filter.Entry, opts *Options) []filter.Entry {
	if opts == nil {
		return entries
	}
	renamed := Run(FromFilterEntries(entries), *opts)
	return ToFilterEntries(renamed, "")
}
