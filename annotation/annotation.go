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
	var sc scanner
	start := 0
	for i, r := range s {
		sep, balanced := sc.step(r)
		if !balanced {
			return nil, false
		}
		if sep {
			items = append(items, s[start:i])
			start = i + 1
		}
	}
	if strings.TrimSpace(s) != "" {
		items = append(items, s[start:])
	}
	return items, !sc.inQuote && sc.depth == 0
}

// scanner tracks double quotes and bracket depth while reading a value.
type scanner struct {
	depth   int
	inQuote bool
	escaped bool
}

// step reads r and reports whether it is a top-level comma and whether no
// closing bracket has appeared without its opening one.
func (s *scanner) step(r rune) (sep bool, balanced bool) {
	switch {
	case s.escaped:
		s.escaped = false
	case s.inQuote:
		s.escaped = r == '\\'
		s.inQuote = r != '"'
	case r == '"':
		s.inQuote = true
	case strings.ContainsRune("([{", r):
		s.depth++
	case strings.ContainsRune(")]}", r):
		s.depth--
	case r == ',':
		return s.depth == 0, true
	}
	return false, s.depth >= 0
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
