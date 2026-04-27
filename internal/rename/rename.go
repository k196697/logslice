package rename

import (
	"fmt"
	"strings"
)

// Entry represents a log entry with fields and an optional timestamp.
type Entry struct {
	Timestamp int64
	Fields    map[string]string
}

// Rule maps an old field name to a new field name.
type Rule struct {
	From string
	To   string
}

// Options configures the rename operation.
type Options struct {
	Rules []Rule
}

// DefaultOptions returns an Options with no rules.
func DefaultOptions() Options {
	return Options{}
}

// ParseRule parses a rename expression of the form "old:new".
func ParseRule(expr string) (Rule, error) {
	parts := strings.SplitN(expr, ":", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("rename: invalid rule %q: expected format old:new", expr)
	}
	from := strings.TrimSpace(parts[0])
	to := strings.TrimSpace(parts[1])
	if from == "" || to == "" {
		return Rule{}, fmt.Errorf("rename: invalid rule %q: field names must not be empty", expr)
	}
	return Rule{From: from, To: to}, nil
}

// ParseRules parses multiple rename expressions.
func ParseRules(exprs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(exprs))
	for _, expr := range exprs {
		r, err := ParseRule(expr)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

// Run applies rename rules to each entry, returning updated entries.
func Run(entries []Entry, opts Options) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]string, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		for _, rule := range opts.Rules {
			if val, ok := fields[rule.From]; ok {
				delete(fields, rule.From)
				fields[rule.To] = val
			}
		}
		result = append(result, Entry{Timestamp: e.Timestamp, Fields: fields})
	}
	return result
}
