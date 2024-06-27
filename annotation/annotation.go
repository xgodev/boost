package annotation

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"strings"
)

type Annotation struct {
	Name  string
	Value string
	Map   map[string]interface{}
}

func NewAnnotation(name string, value string) Annotation {
	mp := asMap(value)
	v := valueMap(mp)
	return Annotation{Name: name, Map: mp, Value: v}
}

func (m *Annotation) RawValue() string {
	return m.Value
}

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

func valueMap(mp map[string]interface{}) string {
	var entries []string
	for k, v := range mp {
		entries = append(entries, strings.Join([]string{k, fmt.Sprintf("%v", v)}, "="))
	}
	return strings.Join(entries, ",")
}

func asMap(value string) map[string]interface{} {
	entries := strings.Split(value, ",")
	mp := make(map[string]interface{})
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		e := strings.Split(entry, "=")
		if len(e) == 2 {
			mp[strings.ReplaceAll(strings.TrimSpace(e[0]), " ", "")] = determineType(strings.TrimSpace(e[1]))
		}
	}
	return mp
}

func determineType(value string) interface{} {
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}
	if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue
	}
	return value
}
