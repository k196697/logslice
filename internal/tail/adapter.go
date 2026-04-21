package tail

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// ParsedEntry mirrors the parser entry type used across the project.
type ParsedEntry = parser.Entry

// ReadLastNEntries reads the last n raw lines from r, then parses them
// using the provided parse function. Lines that fail to parse are skipped.
func ReadLastNEntries(
	r io.Reader,
	n int,
	parseFn func(io.Reader, string) ([]ParsedEntry, error),
	timeField string,
) ([]ParsedEntry, error) {
	lines, err := ReadLastN(r, n)
	if err != nil {
		return nil, fmt.Errorf("tail adapter: %w", err)
	}
	if len(lines) == 0 {
		return []ParsedEntry{}, nil
	}

	joined := strings.Join(lines, "\n")
	entries, err := parseFn(strings.NewReader(joined), timeField)
	if err != nil {
		return nil, fmt.Errorf("tail adapter: parse: %w", err)
	}
	return entries, nil
}
