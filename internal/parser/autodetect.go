package parser

import (
	"bufio"
	"io"
	"strings"
)

// DetectFormat attempts to detect the log format from the first non-empty line
// of the reader. It returns the format name ("json", "csv", "logfmt") or an
// empty string if detection fails.
func DetectFormat(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		return detectLine(line), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}

// DetectFormatFromString attempts to detect the log format from a single line string.
func DetectFormatFromString(line string) string {
	return detectLine(strings.TrimSpace(line))
}

func detectLine(line string) string {
	if line == "" {
		return ""
	}
	// JSON: starts with '{'
	if line[0] == '{' {
		return "json"
	}
	// CSV: contains commas and no '=' signs (naive heuristic)
	if strings.Contains(line, ",") && !strings.Contains(line, "=") {
		return "csv"
	}
	// logfmt: contains key=value pairs
	if strings.Contains(line, "=") {
		return "logfmt"
	}
	return ""
}
