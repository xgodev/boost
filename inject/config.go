package inject

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

// ConfigField is a config key of a @Config struct.
type ConfigField struct {
	Name    string // Go name of the field
	Segment string // key segment under the struct's root key
	Type    string // supported type of the field ("time.Duration" for a duration)
	Default string // Go literal of the default value, of the field type
	Desc    string // description of the key
	Hide    bool   // the value is a secret (config.WithHide)
}

// configTags are the struct tags that make a field a config key.
var configTags = []string{"default", "desc", "hide", "config"}

// envSegment is a key segment an environment variable can reach.
var envSegment = regexp.MustCompile(`^[a-z][a-zA-Z0-9]*$`)

// bitSizes are the bit sizes of the supported numeric types, 0 for int and
// uint (the platform size).
var bitSizes = map[string]int{
	"int": 0, "int8": 8, "int16": 16, "int32": 32, "int64": 64,
	"uint": 0, "uint8": 8, "uint16": 16, "uint32": 32, "uint64": 64,
	"float32": 32, "float64": 64,
}

// newConfig returns the component of a @Config struct: a provider of
// *Struct whose fields are config keys.
func newConfig(entry annotation.Entry, parsed []parsedAnnotation) (Component, bool, error) {
	if err := checkConfigDecl(entry, parsed); err != nil {
		return Component{}, false, err
	}

	comp := Component{
		Entry:    entry,
		Type:     AnnotationTypePROVIDE,
		Provides: map[int]Annotation{0: {}},
		Injects:  map[int]Annotation{},
	}
	for _, f := range entry.Fields {
		cf, ok, err := configField(entry, f)
		if err != nil {
			return Component{}, false, err
		}
		if ok {
			comp.Config = append(comp.Config, cf)
		}
	}
	return comp, true, nil
}

// checkConfigDecl fails for a @Config the generated code cannot use.
func checkConfigDecl(entry annotation.Entry, parsed []parsedAnnotation) error {
	where := gid(entry)
	switch {
	case !entry.IsStruct():
		return errors.NotValidf("@Config on %s: only a struct holds config keys", where)
	case !token.IsExported(entry.Struct):
		return errors.NotValidf("@Config on %s: the struct must be exported to be provided by the generated module", where)
	case entry.Package == "main":
		return errors.NotValidf("@Config on %s: package main cannot be imported by the generated module; move it to another package", where)
	}
	for _, p := range parsed {
		if p.kind != AnnotationTypeCONFIG {
			return errors.NotValidf("@Config on %s cannot be combined with @%s: the struct is provided by its generated module", where, upperFirst(strings.ToLower(p.kind.String())))
		}
	}
	return nil
}

// configField returns the config key of a struct field, or false when the
// field is not a key (config:"-", or unexported without config tags).
func configField(entry annotation.Entry, f annotation.EntryField) (ConfigField, bool, error) {
	where := gid(entry) + "." + f.Name
	tag := reflect.StructTag(f.Tag)

	switch {
	case f.Name == "":
		return ConfigField{}, false, errors.NotValidf("@Config on %s: embedded field %s is not supported; declare its fields", gid(entry), f.Type)
	case tag.Get("config") == "-":
		return ConfigField{}, false, nil
	case !token.IsExported(f.Name):
		if hasConfigTag(tag) {
			return ConfigField{}, false, errors.NotValidf("@Config field %s is unexported and cannot be filled from config; export it or remove its tags", where)
		}
		return ConfigField{}, false, nil
	}

	cf := ConfigField{Name: f.Name, Segment: fieldSegment(f.Name)}
	var err error
	if cf.Type, err = fieldType(entry, f, where); err != nil {
		return ConfigField{}, false, err
	}
	if err = applyTags(&cf, tag, where); err != nil {
		return ConfigField{}, false, err
	}
	return cf, true, nil
}

// applyTags reads desc, config, hide and default into cf.
func applyTags(cf *ConfigField, tag reflect.StructTag, where string) error {
	var ok bool
	if cf.Desc, ok = tag.Lookup("desc"); !ok || strings.TrimSpace(cf.Desc) == "" {
		return errors.NotValidf("@Config field %s needs a desc tag: every config key is documented", where)
	}
	if name, ok := tag.Lookup("config"); ok {
		if !envSegment.MatchString(name) {
			return errors.NotValidf("@Config field %s: config:%q must be lowerCamel (letters and digits) so an env var reaches it", where, name)
		}
		cf.Segment = name
	}
	if hide, ok := tag.Lookup("hide"); ok {
		h, err := strconv.ParseBool(hide)
		if err != nil {
			return errors.NotValidf("@Config field %s: hide:%q is not a bool", where, hide)
		}
		cf.Hide = h
	}

	value, set := tag.Lookup("default")
	literal, err := defaultLiteral(cf.Type, value, set)
	if err != nil {
		return errors.NotValidf("@Config field %s: default:%q is not a valid %s: %v", where, value, cf.Type, err)
	}
	cf.Default = literal
	return nil
}

func hasConfigTag(tag reflect.StructTag) bool {
	for _, name := range configTags {
		if _, ok := tag.Lookup(name); ok {
			return true
		}
	}
	return false
}

// fieldType returns the supported type of a field: a basic type, []string,
// or time.Duration written through the file's import of time.
func fieldType(entry annotation.Entry, f annotation.EntryField, where string) (string, error) {
	if _, ok := bitSizes[f.Type]; ok || f.Type == "string" || f.Type == "bool" || f.Type == "[]string" {
		return f.Type, nil
	}
	if expr, err := parser.ParseExpr(f.Type); err == nil {
		if sel, ok := expr.(*ast.SelectorExpr); ok && sel.Sel.Name == "Duration" {
			if pkg, ok := sel.X.(*ast.Ident); ok && entry.Imports[pkg.Name] == "time" {
				return "time.Duration", nil
			}
		}
	}
	return "", errors.NotValidf("@Config field %s has type %s; supported: string, bool, int and uint types, float32, float64, time.Duration and []string", where, f.Type)
}

// defaultLiteral parses the default of a field and returns it as a Go
// literal of the field type. Without a default it is the zero value.
func defaultLiteral(tp, value string, set bool) (string, error) {
	switch tp {
	case "string":
		return strconv.Quote(value), nil
	case "[]string":
		return stringsLiteral(value), nil
	case "time.Duration":
		return durationLiteral(value, set)
	}
	switch {
	case set:
	case tp == "bool":
		value = "false"
	default:
		value = "0"
	}
	return scalarLiteral(tp, value)
}

func scalarLiteral(tp, value string) (string, error) {
	var literal string
	var err error
	switch {
	case tp == "bool":
		var b bool
		b, err = strconv.ParseBool(value)
		return strconv.FormatBool(b), err
	case strings.HasPrefix(tp, "float"):
		var f float64
		f, err = strconv.ParseFloat(value, bitSizes[tp])
		literal = strconv.FormatFloat(f, 'g', -1, bitSizes[tp])
	case strings.HasPrefix(tp, "uint"):
		var u uint64
		u, err = strconv.ParseUint(value, 10, bitSizes[tp])
		literal = strconv.FormatUint(u, 10)
	default:
		var i int64
		i, err = strconv.ParseInt(value, 10, bitSizes[tp])
		literal = strconv.FormatInt(i, 10)
	}
	if err != nil || tp == "int" {
		return literal, err
	}
	return tp + "(" + literal + ")", nil
}

// stringsLiteral splits a comma list, trimming each item.
func stringsLiteral(value string) string {
	if strings.TrimSpace(value) == "" {
		return "[]string{}"
	}
	items := strings.Split(value, ",")
	for i, item := range items {
		items[i] = strconv.Quote(strings.TrimSpace(item))
	}
	return "[]string{" + strings.Join(items, ", ") + "}"
}

var durationUnits = []struct {
	unit time.Duration
	name string
}{
	{time.Hour, "time.Hour"},
	{time.Minute, "time.Minute"},
	{time.Second, "time.Second"},
	{time.Millisecond, "time.Millisecond"},
	{time.Microsecond, "time.Microsecond"},
	{time.Nanosecond, "time.Nanosecond"},
}

// durationLiteral writes a duration in its largest whole unit (90 *
// time.Second for 1m30s).
func durationLiteral(value string, set bool) (string, error) {
	if !set {
		return "time.Duration(0)", nil
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return "", err
	}
	if d == 0 {
		return "time.Duration(0)", nil
	}
	for _, u := range durationUnits {
		if d%u.unit == 0 {
			return fmt.Sprintf("%d * %s", d/u.unit, u.name), nil
		}
	}
	return "", nil // unreachable: every duration is a whole number of nanoseconds
}

// fieldWords splits a Go field name into words: "_" separates them, a
// capital starts one after a lowercase letter or a digit, and the last
// capital of an acronym starts the next word (HTTPTimeout -> HTTP, Timeout).
func fieldWords(name string) []string {
	runes := []rune(name)
	var words []string
	start := 0
	flush := func(end int) {
		if end > start {
			words = append(words, string(runes[start:end]))
		}
	}
	for i, r := range runes {
		switch {
		case r == '_':
			flush(i)
			start = i + 1
		case i > start && unicode.IsUpper(r) && startsWord(runes, i):
			flush(i)
			start = i
		}
	}
	flush(len(runes))
	return words
}

// startsWord reports whether the capital at i starts a word.
func startsWord(runes []rune, i int) bool {
	prev := runes[i-1]
	if unicode.IsLower(prev) || unicode.IsDigit(prev) {
		return true
	}
	return unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1])
}

// fieldSegment is the lowerCamel key segment of a field name (BaseURL ->
// baseUrl): what the env loader builds from BASE__URL.
func fieldSegment(name string) string {
	var b strings.Builder
	for i, word := range fieldWords(name) {
		word = strings.ToLower(word)
		if i > 0 {
			word = upperFirst(word)
		}
		b.WriteString(word)
	}
	return b.String()
}
