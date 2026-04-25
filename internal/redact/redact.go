package redact

import (
	"regexp"
	"strings"
)

// Rule describes a single redaction rule.
type Rule struct {
	Field   string
	Pattern *regexp.Regexp // nil means redact entire field value
	Mask    string
}

// Options controls redaction behaviour.
type Options struct {
	Rules []Rule
}

// DefaultOptions returns Options with a sensible mask.
func DefaultOptions() Options {
	return Options{}
}

// Entry is the minimal interface redact operates on.
type Entry struct {
	Fields map[string]string
}

// Run applies all rules to every entry and returns the (possibly modified) slice.
func Run(entries []Entry, opts Options) []Entry {
	if len(opts.Rules) == 0 {
		return entries
	}
	out := make([]Entry, len(entries))
	for i, e := range entries {
		out[i] = applyRules(e, opts.Rules)
	}
	return out
}

func applyRules(e Entry, rules []Rule) Entry {
	copy := Entry{Fields: make(map[string]string, len(e.Fields))}
	for k, v := range e.Fields {
		copy.Fields[k] = v
	}
	for _, r := range rules {
		v, ok := copy.Fields[r.Field]
		if !ok {
			continue
		}
		mask := r.Mask
		if mask == "" {
			mask = "[REDACTED]"
		}
		if r.Pattern == nil {
			copy.Fields[r.Field] = mask
		} else {
			copy.Fields[r.Field] = r.Pattern.ReplaceAllString(v, mask)
		}
	}
	return copy
}

// ParseRule parses a rule expression of the form:
//   field                  — redact entire field
//   field:pattern          — replace regex matches with [REDACTED]
//   field:pattern=mask     — replace regex matches with custom mask
func ParseRule(expr string) (Rule, error) {
	if !strings.Contains(expr, ":") {
		return Rule{Field: expr}, nil
	}
	parts := strings.SplitN(expr, ":", 2)
	field := parts[0]
	rest := parts[1]
	mask := "[REDACTED]"
	if idx := strings.Index(rest, "="); idx != -1 {
		mask = rest[idx+1:]
		rest = rest[:idx]
	}
	re, err := regexp.Compile(rest)
	if err != nil {
		return Rule{}, err
	}
	return Rule{Field: field, Pattern: re, Mask: mask}, nil
}

// ParseRules parses multiple rule expressions.
func ParseRules(exprs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(exprs))
	for _, e := range exprs {
		r, err := ParseRule(e)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}
