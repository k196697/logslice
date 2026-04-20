package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

// ParseLogfmtLines reads logfmt-encoded lines from r and returns a slice of
// parsed log entries. Lines that cannot be parsed are silently skipped.
// timestampField specifies which key holds the timestamp value.
func ParseLogfmtLines(r io.Reader, timestampField string) ([]map[string]interface{}, error) {
	if timestampField == "" {
		timestampField = "time"
	}

	var entries []map[string]interface{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry, err := parseLogfmtLine(line)
		if err != nil {
			continue
		}

		if raw, ok := entry[timestampField]; ok {
			if s, ok := raw.(string); ok {
				if t, err := parseTimestampString(s); err == nil {
					entry[timestampField] = t.Format(time.RFC3339Nano)
				}
			}
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("logfmt scanner error: %w", err)
	}

	return entries, nil
}

// parseLogfmtLine parses a single logfmt line into a map.
// Format: key=value key2="value with spaces" bare_key
func parseLogfmtLine(line string) (map[string]interface{}, error) {
	entry := make(map[string]interface{})
	remaining := strings.TrimSpace(line)

	for remaining != "" {
		// Find key
		eqIdx := strings.IndexByte(remaining, '=')
		if eqIdx == -1 {
			// bare key with no value
			key := strings.TrimSpace(remaining)
			if key != "" {
				entry[key] = true
			}
			break
		}

		key := strings.TrimSpace(remaining[:eqIdx])
		remaining = remaining[eqIdx+1:]

		var value string
		if strings.HasPrefix(remaining, "\"") {
			// quoted value
			closing := strings.Index(remaining[1:], "\"")
			if closing == -1 {
				return nil, fmt.Errorf("unclosed quote in logfmt line")
			}
			value = remaining[1 : closing+1]
			remaining = strings.TrimSpace(remaining[closing+2:])
		} else {
			// unquoted value — ends at next space
			spaceIdx := strings.IndexByte(remaining, ' ')
			if spaceIdx == -1 {
				value = remaining
				remaining = ""
			} else {
				value = remaining[:spaceIdx]
				remaining = strings.TrimSpace(remaining[spaceIdx+1:])
			}
		}

		if key != "" {
			entry[key] = value
		}
	}

	if len(entry) == 0 {
		return nil, fmt.Errorf("empty entry")
	}
	return entry, nil
}
