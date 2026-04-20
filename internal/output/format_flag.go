package output

import (
	"fmt"
	"strings"
)

// FormatFlag implements flag.Value for the --format CLI flag.
type FormatFlag struct {
	Value Format
}

var validFormats = []Format{FormatJSON, FormatText, FormatCompact}

// String returns the current format value as a string.
func (f *FormatFlag) String() string {
	if f.Value == "" {
		return string(FormatJSON)
	}
	return string(f.Value)
}

// Set parses and validates the format flag value.
func (f *FormatFlag) Set(s string) error {
	norm := Format(strings.ToLower(strings.TrimSpace(s)))
	for _, v := range validFormats {
		if norm == v {
			f.Value = norm
			return nil
		}
	}
	return fmt.Errorf("invalid format %q: must be one of %s", s, formatList())
}

// Type returns the type name for help text.
func (f *FormatFlag) Type() string {
	return "format"
}

func formatList() string {
	strs := make([]string, len(validFormats))
	for i, v := range validFormats {
		strs[i] = string(v)
	}
	return strings.Join(strs, ", ")
}
