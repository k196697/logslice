package pipeline

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

// parseLines dispatches to the correct parser based on format.
func parseLines(format, timeField string, lines []string) ([]filter.Entry, error) {
	reader := strings.NewReader(strings.Join(lines, "\n"))
	_ = reader // parsers accept []string directly

	switch strings.ToLower(format) {
	case "json":
		parsed, err := parser.ParseJSONLines(lines, timeField)
		if err != nil {
			return nil, err
		}
		return toFilterEntries(parsed), nil

	case "csv":
		parsed, err := parser.ParseCSVLines(lines, timeField)
		if err != nil {
			return nil, err
		}
		return toFilterEntries(parsed), nil

	case "logfmt":
		parsed, err := parser.ParseLogfmtLines(lines, timeField)
		if err != nil {
			return nil, err
		}
		return toFilterEntries(parsed), nil

	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}

// toFilterEntries converts parser.Entry slice to filter.Entry slice.
func toFilterEntries(parsed []parser.Entry) []filter.Entry {
	out := make([]filter.Entry, 0, len(parsed))
	for _, p := range parsed {
		out = append(out, filter.Entry{
			Timestamp: p.Timestamp,
			Fields:    p.Fields,
			Raw:       p.Raw,
		})
	}
	return out
}
