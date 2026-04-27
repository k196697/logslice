package mask

import (
	"strings"

	"github.com/user/logslice/internal/filter"
)

// FromFilterEntries converts filter.Entry slice to mask.Entry slice.
func FromFilterEntries(entries []filter.Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		me := Entry{Fields: make(map[string]string, len(e.Fields))}
		for k, v := range e.Fields {
			me.Fields[k] = fmt.Sprintf("%v", v)
		}
		out = append(out, me)
	}
	return out
}

// ToFilterEntries converts mask.Entry slice back to filter.Entry slice,
// preserving the original timestamps from the source slice.
func ToFilterEntries(masked []Entry, originals []filter.Entry) []filter.Entry {
	out := make([]filter.Entry, 0, len(masked))
	for i, me := range masked {
		fe := filter.Entry{
			Fields: make(map[string]interface{}, len(me.Fields)),
		}
		if i < len(originals) {
			fe.Timestamp = originals[i].Timestamp
		}
		for k, v := range me.Fields {
			fe.Fields[k] = v
		}
		out = append(out, fe)
	}
	return out
}

// ParseConfig returns an Options pointer if any --mask flags were provided,
// or nil when masking is disabled.
func ParseConfig(fields []string, keepPrefix, keepSuffix int, char string) *Options {
	if len(fields) == 0 {
		return nil
	}
	cleaned := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f != "" {
			cleaned = append(cleaned, f)
		}
	}
	if len(cleaned) == 0 {
		return nil
	}
	if char == "" {
		char = "*"
	}
	return &Options{
		Fields:     cleaned,
		Char:       char,
		KeepPrefix: keepPrefix,
		KeepSuffix: keepSuffix,
	}
}

// RunFromConfig applies masking using filter.Entry types if opts is non-nil.
func RunFromConfig(entries []filter.Entry, opts *Options) []filter.Entry {
	if opts == nil {
		return entries
	}
	me := FromFilterEntries(entries)
	masked, err := Run(me, *opts)
	if err != nil || len(masked) == 0 {
		return entries
	}
	return ToFilterEntries(masked, entries)
}
