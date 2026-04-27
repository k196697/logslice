package typecast

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule defines a field-to-type casting instruction.
type Rule struct {
	Field string
	TargetType string // "int", "float", "bool", "string"
}

// Options holds configuration for the typecast run.
type Options struct {
	Rules []Rule
}

// DefaultOptions returns an Options with no rules.
func DefaultOptions() Options {
	return Options{}
}

// Entry is a log entry with raw string fields and a timestamp.
type Entry struct {
	Timestamp int64
	Fields    map[string]interface{}
}

// ParseRule parses a rule string of the form "field:type".
func ParseRule(expr string) (Rule, error) {
	parts := strings.SplitN(expr, ":", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("typecast: invalid rule %q: expected field:type", expr)
	}
	field := strings.TrimSpace(parts[0])
	target := strings.TrimSpace(parts[1])
	if field == "" {
		return Rule{}, fmt.Errorf("typecast: empty field name in rule %q", expr)
	}
	switch target {
	case "int", "float", "bool", "string":
	default:
		return Rule{}, fmt.Errorf("typecast: unsupported type %q in rule %q", target, expr)
	}
	return Rule{Field: field, TargetType: target}, nil
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

// Run applies typecast rules to each entry, converting field values in place.
func Run(entries []Entry, opts Options) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		for _, rule := range opts.Rules {
			val, ok := fields[rule.Field]
			if !ok {
				continue
			}
			str := fmt.Sprintf("%v", val)
			converted, err := castValue(str, rule.TargetType)
			if err == nil {
				fields[rule.Field] = converted
			}
		}
		result = append(result, Entry{Timestamp: e.Timestamp, Fields: fields})
	}
	return result
}

func castValue(s, targetType string) (interface{}, error) {
	switch targetType {
	case "int":
		return strconv.ParseInt(s, 10, 64)
	case "float":
		return strconv.ParseFloat(s, 64)
	case "bool":
		return strconv.ParseBool(s)
	case "string":
		return s, nil
	}
	return nil, fmt.Errorf("unsupported type: %s", targetType)
}
