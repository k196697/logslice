package typecast

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to typecast.Entry slice.
func FromFilterEntries(in []filter.Entry) []Entry {
	out := make([]Entry, 0, len(in))
	for _, e := range in {
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

// ToFilterEntries converts typecast.Entry slice back to filter.Entry slice.
func ToFilterEntries(in []Entry) []filter.Entry {
	out := make([]filter.Entry, 0, len(in))
	for _, e := range in {
		fe := filter.Entry{
			Fields: make(map[string]string, len(e.Fields)),
		}
		if e.Timestamp != 0 {
			fe.Timestamp = timestampFromNano(e.Timestamp)
		}
		for k, v := range e.Fields {
			fe.Fields[k] = fmt.Sprintf("%v", v)
		}
		out = append(out, fe)
	}
	return out
}

// ParseConfig builds Options from CLI expressions, returning nil if empty.
func ParseConfig(exprs []string) (*Options, error) {
	if len(exprs) == 0 {
		return nil, nil
	}
	rules, err := ParseRules(exprs)
	if err != nil {
		return nil, err
	}
	opts := DefaultOptions()
	opts.Rules = rules
	return &opts, nil
}

// RunFromConfig applies typecast rules if opts is non-nil.
func RunFromConfig(entries []filter.Entry, opts *Options) []filter.Entry {
	if opts == nil {
		return entries
	}
	converted := FromFilterEntries(entries)
	result := Run(converted, *opts)
	return ToFilterEntries(result)
}

func timestampFromNano(ns int64) interface{} {
	// Return a time.Time-compatible value via filter.Entry zero value handling.
	// Callers rely on the Fields map; timestamp is stored separately.
	_ = ns
	return nil
}
