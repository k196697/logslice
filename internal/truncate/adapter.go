package truncate

import (
	"github.com/yourorg/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to truncate.Entry slice.
func FromFilterEntries(entries []filter.Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		te := Entry{Fields: make(map[string]string, len(e.Fields))}
		for k, v := range e.Fields {
			te.Fields[k] = v
		}
		out = append(out, te)
	}
	return out
}

// ToFilterEntries converts truncate.Entry slice back to filter.Entry slice,
// preserving the original timestamp from the source entries.
func ToFilterEntries(truncated []Entry, originals []filter.Entry) []filter.Entry {
	out := make([]filter.Entry, 0, len(truncated))
	for i, te := range truncated {
		fe := filter.Entry{
			Fields: make(map[string]string, len(te.Fields)),
		}
		if i < len(originals) {
			fe.Timestamp = originals[i].Timestamp
			fe.Raw = originals[i].Raw
		}
		for k, v := range te.Fields {
			fe.Fields[k] = v
		}
		out = append(out, fe)
	}
	return out
}

// ParseConfig builds Options from CLI-style parameters.
// fieldsExpr is a comma-separated list of field names (empty means all fields).
func ParseConfig(fieldsExpr string, maxLength int, suffix string) (*Options, error) {
	if maxLength <= 0 && fieldsExpr == "" && suffix == "" {
		return nil, nil
	}
	opts := DefaultOptions()
	if maxLength > 0 {
		opts.MaxLength = maxLength
	}
	if suffix != "" {
		opts.Suffix = suffix
	}
	opts.Fields = ParseFields(fieldsExpr)
	return &opts, nil
}

// RunFromConfig applies truncation using filter.Entry types.
func RunFromConfig(entries []filter.Entry, opts *Options) ([]filter.Entry, error) {
	if opts == nil {
		return entries, nil
	}
	truncated, err := Run(FromFilterEntries(entries), *opts)
	if err != nil {
		return nil, err
	}
	return ToFilterEntries(truncated, entries), nil
}
