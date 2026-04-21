package highlight

import (
	"fmt"
	"strings"
)

// ParseRule parses a highlight rule expression of the form:
//   field=value:color   e.g. "level=error:red"
//   field=*:color       e.g. "host=*:cyan"
func ParseRule(expr string) (Rule, error) {
	// Split on last colon to get color
	colonIdx := strings.LastIndex(expr, ":")
	if colonIdx < 0 {
		return Rule{}, fmt.Errorf("highlight rule %q missing color suffix (expected field=value:color)", expr)
	}

	colorName := expr[colonIdx+1:]
	kv := expr[:colonIdx]

	eqIdx := strings.Index(kv, "=")
	if eqIdx < 0 {
		return Rule{}, fmt.Errorf("highlight rule %q missing '=' between field and value", expr)
	}

	field := kv[:eqIdx]
	value := kv[eqIdx+1:]

	if field == "" {
		return Rule{}, fmt.Errorf("highlight rule %q has empty field name", expr)
	}

	color, ok := ParseColor(colorName)
	if !ok {
		return Rule{}, fmt.Errorf("highlight rule %q has unknown color %q", expr, colorName)
	}

	return Rule{Field: field, Value: value, Color: color}, nil
}

// ParseRules parses multiple rule expressions and returns all valid rules.
// Errors are collected and returned together.
func ParseRules(exprs []string) ([]Rule, error) {
	var rules []Rule
	var errs []string
	for _, expr := range exprs {
		r, err := ParseRule(expr)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		rules = append(rules, r)
	}
	if len(errs) > 0 {
		return rules, fmt.Errorf("rule parse errors: %s", strings.Join(errs, "; "))
	}
	return rules, nil
}
