package highlight

import (
	"fmt"
	"strings"
)

// Color ANSI codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Rule defines a field-based highlight rule.
type Rule struct {
	Field string
	Value string
	Color string
}

// Highlighter applies color rules to log line output.
type Highlighter struct {
	rules   []Rule
	enabled bool
}

// New creates a Highlighter. If enabled is false, all methods return plain text.
func New(enabled bool, rules []Rule) *Highlighter {
	return &Highlighter{rules: rules, enabled: enabled}
}

// ColorizeField wraps a field value in the matching color, if any rule applies.
func (h *Highlighter) ColorizeField(field, value string) string {
	if !h.enabled {
		return value
	}
	for _, r := range h.rules {
		if r.Field == field && (r.Value == "*" || r.Value == value) {
			return fmt.Sprintf("%s%s%s", r.Color, value, Reset)
		}
	}
	return value
}

// ColorizeLine applies a color to an entire line if the line contains the given keyword.
func (h *Highlighter) ColorizeLine(line, keyword, color string) string {
	if !h.enabled || keyword == "" {
		return line
	}
	if strings.Contains(line, keyword) {
		return fmt.Sprintf("%s%s%s", color, line, Reset)
	}
	return line
}

// ParseColor maps a color name string to its ANSI escape code.
func ParseColor(name string) (string, bool) {
	switch strings.ToLower(name) {
	case "red":
		return Red, true
	case "green":
		return Green, true
	case "yellow":
		return Yellow, true
	case "blue":
		return Blue, true
	case "cyan":
		return Cyan, true
	case "bold":
		return Bold, true
	}
	return "", false
}
