package annotation

import (
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Annotation is a single `// @Name(key=value, ...)` line found in a doc comment.
type Annotation struct {
	Name  string                 // Name of the annotation, without the "@"
	Value string                 // Text between the parentheses, as written in the source
	Map   map[string]interface{} // Parsed key=value pairs
}

// NewAnnotation builds an annotation from its name and the text between its
// parentheses. Values are split on top-level commas: commas inside double
// quotes or brackets belong to the value, and only the first "=" separates
// the key. Unquoted values become int, float64 or bool (true/false) when they
// parse as such; quoted values are always strings.
func NewAnnotation(name string, value string) Annotation {
	return Annotation{Name: name, Map: asMap(value), Value: value}
}

// RawValue returns the text between the parentheses, as written in the source.
func (m *Annotation) RawValue() string {
	return m.Value
}

// Decode copies the parsed values into a, matching keys with `attr` tags.
func (m *Annotation) Decode(a interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName: "attr",
		Result:  &a,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(m.Map)
}

func asMap(value string) map[string]interface{} {
	mp := make(map[string]interface{})
	items, _ := splitTopLevel(value)
	for _, item := range items {
		key, val, found := strings.Cut(item, "=")
		if !found {
			continue
		}
		mp[strings.ReplaceAll(strings.TrimSpace(key), " ", "")] = typedValue(strings.TrimSpace(val))
	}
	return mp
}

// splitTopLevel splits s on the commas that are outside double quotes and
// brackets. ok is false when quotes or brackets are unbalanced.
func splitTopLevel(s string) (items []string, ok bool) {
	var depth, start int
	var inQuote, escaped bool
	for i, r := range s {
		switch {
		case escaped:
			escaped = false
		case inQuote && r == '\\':
			escaped = true
		case r == '"':
			inQuote = !inQuote
		case inQuote:
		case r == '(' || r == '[' || r == '{':
			depth++
		case r == ')' || r == ']' || r == '}':
			depth--
			if depth < 0 {
				return nil, false
			}
		case r == ',' && depth == 0:
			items = append(items, s[start:i])
			start = i + 1
		}
	}
	if strings.TrimSpace(s) != "" {
		items = append(items, s[start:])
	}
	return items, !inQuote && depth == 0
}

func typedValue(value string) interface{} {
	if len(value) >= 2 && strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		if unquoted, err := strconv.Unquote(value); err == nil {
			return unquoted
		}
		return value[1 : len(value)-1]
	}
	return determineType(value)
}

func determineType(value string) interface{} {
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	if strings.ContainsAny(value, "0123456789") {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	if strings.EqualFold(value, "true") || strings.EqualFold(value, "false") {
		return strings.EqualFold(value, "true")
	}
	return value
}
