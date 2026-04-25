package redact

import (
	"github.com/user/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to redact.Entry slice.
func FromFilterEntries(in []filter.Entry) []Entry {
	out := make([]Entry, len(in))
	for i, e := range in {
		fields := make(map[string]string, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		out[i] = Entry{Fields: fields}
	}
	return out
}

// ToFilterEntries converts redact.Entry slice back to filter.Entry slice,
// preserving the original timestamps from the source slice.
func ToFilterEntries(redacted []Entry, original []filter.Entry) []filter.Entry {
	out := make([]filter.Entry, len(redacted))
	for i, e := range redacted {
		out[i] = filter.Entry{
			Timestamp: original[i].Timestamp,
			Fields:    e.Fields,
			Raw:       original[i].Raw,
		}
	}
	return out
}

// ParseConfig builds an Options from CLI flag values.
// exprs is the list of --redact flag values.
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

// RunFromConfig applies redaction to filter entries when opts is non-nil.
func RunFromConfig(entries []filter.Entry, opts *Options) []filter.Entry {
	if opts == nil {
		return entries
	}
	redacted := Run(FromFilterEntries(entries), *opts)
	return ToFilterEntries(redacted, entries)
}
