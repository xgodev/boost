// Package config provides a wrapper around koanf.
package config

import (
	"time"

	"github.com/knadh/koanf"
)

// UnmarshalWithPath unmarshals a given key path into the given struct using the mapstructure lib.
func UnmarshalWithPath(path string, o interface{}) error {
	return instance.UnmarshalWithConf(path, &o, koanf.UnmarshalConf{Tag: "config"})
}

// Unmarshal unmarshals the given struct using the mapstructure lib. The whole map is unmarshalled.
func Unmarshal(o interface{}) error {
	return instance.UnmarshalWithConf("", &o, koanf.UnmarshalConf{Tag: "config"})
}

// Exists returns true if the given key path exists in the conf map.
func Exists(path string) bool {
	return instance.Exists(path)
}

// Int64 returns the int64 value of a given key path or 0 if the path
// does not exist or if the value is not a valid int64.
func Int64(path string) int64 {
	return instance.Int64(path)
}

// Int64s returns the []int64 slice value of a given key path or an
// empty []int64 slice if the path does not exist or if the value
// is not a valid int slice.
func Int64s(path string) []int64 {
	return instance.Int64s(path)
}

// Int64Map returns the map[string]int64 value of a given key path
// or an empty map[string]int64 if the path does not exist or if the
// value is not a valid int64 map.
func Int64Map(path string) map[string]int64 {
	return instance.Int64Map(path)
}

// Int returns the int value of a given key path or 0 if the path
// does not exist or if the value is not a valid int.
func Int(path string) int {
	return instance.Int(path)
}

// Ints returns the []int slice value of a given key path or an
// empty []int slice if the path does not exist or if the value
// is not a valid int slice.
func Ints(path string) []int {
	return instance.Ints(path)
}

// IntMap returns the map[string]int value of a given key path
// or an empty map[string]int if the path does not exist or if the
// value is not a valid int map.
func IntMap(path string) map[string]int {
	return instance.IntMap(path)
}

// Float64 returns the float64 value of a given key path or 0 if the path
// does not exist or if the value is not a valid float64.
func Float64(path string) float64 {
	return instance.Float64(path)
}

// Float64s returns the []float64 slice value of a given key path or an
// empty []float64 slice if the path does not exist or if the value
// is not a valid float64 slice.
func Float64s(path string) []float64 {
	return instance.Float64s(path)
}

// Float64Map returns the map[string]float64 value of a given key path
// or an empty map[string]float64 if the path does not exist or if the
// value is not a valid float64 map.
func Float64Map(path string) map[string]float64 {
	return instance.Float64Map(path)
}

// Duration returns the time.Duration value of a given key path assuming
// that the key contains a valid numeric value.
func Duration(path string) time.Duration {
	return instance.Duration(path)
}

// Time attempts to parse the value of a given key path and return time.Time
// representation. If the value is numeric, it is treated as a UNIX timestamp
// and if it's string, a parse is attempted with the given layout.
func Time(path, layout string) time.Time {
	return instance.Time(path, layout)
}

// String returns the string value of a given key path or "" if the path
// does not exist or if the value is not a valid string.
func String(path string) string {
	return instance.String(path)
}

// Strings returns the []string slice value of a given key path or an
// empty []string slice if the path does not exist or if the value
// is not a valid string slice.
func Strings(path string) []string {
	return instance.Strings(path)
}

// StringMap returns the map[string]string value of a given key path
// or an empty map[string]string if the path does not exist or if the
// value is not a valid string map.
func StringMap(path string) map[string]string {
	return instance.StringMap(path)
}

// Bytes returns the []byte value of a given key path or an empty
// []byte slice if the path does not exist or if the value is not a valid string.
func Bytes(path string) []byte {
	return instance.Bytes(path)
}

// Bool returns the bool value of a given key path or false if the path
// does not exist or if the value is not a valid bool representation.
// Accepted string representations of bool are the ones supported by strconv.ParseBool.
func Bool(path string) bool {
	return instance.Bool(path)
}

// Bools returns the []bool slice value of a given key path or an
// empty []bool slice if the path does not exist or if the value
// is not a valid bool slice.
func Bools(path string) []bool {
	return instance.Bools(path)
}

// BoolMap returns the map[string]bool value of a given key path
// or an empty map[string]bool if the path does not exist or if the
// value is not a valid bool map.
func BoolMap(path string) map[string]bool {
	return instance.BoolMap(path)
}

// All returns all configs
func All() map[string]interface{} {
	return instance.All()
}

// Get returns interface{} value
func Get(path string) interface{} {
	return instance.Get(path)
}
