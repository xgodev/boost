package config

import (
	"time"
)

type Provider interface {
	Load([]Config)

	// UnmarshalWithPath unmarshals a given key path into the given struct using the mapstructure lib.
	UnmarshalWithPath(path string, o interface{}) error

	// Unmarshal unmarshals the given struct using the mapstructure lib. The whole map is unmarshalled.
	Unmarshal(o interface{}) error

	// Exists returns true if the given key path exists in the conf map.
	Exists(path string) bool

	// Int64 returns the int64 value of a given key path or 0 if the path
	// does not exist or if the value is not a valid int64.
	Int64(path string) int64

	// Int64s returns the []int64 slice value of a given key path or an
	// empty []int64 slice if the path does not exist or if the value
	// is not a valid int slice.
	Int64s(path string) []int64

	// Int64Map returns the map[string]int64 value of a given key path
	// or an empty map[string]int64 if the path does not exist or if the
	// value is not a valid int64 map.
	Int64Map(path string) map[string]int64

	// Int returns the int value of a given key path or 0 if the path
	// does not exist or if the value is not a valid int.
	Int(path string) int

	// Ints returns the []int slice value of a given key path or an
	// empty []int slice if the path does not exist or if the value
	// is not a valid int slice.
	Ints(path string) []int

	// IntMap returns the map[string]int value of a given key path
	// or an empty map[string]int if the path does not exist or if the
	// value is not a valid int map.
	IntMap(path string) map[string]int

	// Float64 returns the float64 value of a given key path or 0 if the path
	// does not exist or if the value is not a valid float64.
	Float64(path string) float64

	// Float64s returns the []float64 slice value of a given key path or an
	// empty []float64 slice if the path does not exist or if the value
	// is not a valid float64 slice.
	Float64s(path string) []float64

	// Float64Map returns the map[string]float64 value of a given key path
	// or an empty map[string]float64 if the path does not exist or if the
	// value is not a valid float64 map.
	Float64Map(path string) map[string]float64

	// Duration returns the time.Duration value of a given key path assuming
	// that the key contains a valid numeric value.
	Duration(path string) time.Duration

	// Time attempts to parse the value of a given key path and return time.Time
	// representation. If the value is numeric, it is treated as a UNIX timestamp
	// and if it's string, a parse is attempted with the given layout.
	Time(path, layout string) time.Time

	// String returns the string value of a given key path or "" if the path
	// does not exist or if the value is not a valid string.
	String(path string) string

	// Strings returns the []string slice value of a given key path or an
	// empty []string slice if the path does not exist or if the value
	// is not a valid string slice.
	Strings(path string) []string
	// StringMap returns the map[string]string value of a given key path
	// or an empty map[string]string if the path does not exist or if the
	// value is not a valid string map.
	StringMap(path string) map[string]string

	// Bytes returns the []byte value of a given key path or an empty
	// []byte slice if the path does not exist or if the value is not a valid string.
	Bytes(path string) []byte
	// Bool returns the bool value of a given key path or false if the path
	// does not exist or if the value is not a valid bool representation.
	// Accepted string representations of bool are the ones supported by strconv.ParseBool.
	Bool(path string) bool

	// Bools returns the []bool slice value of a given key path or an
	// empty []bool slice if the path does not exist or if the value
	// is not a valid bool slice.
	Bools(path string) []bool

	// BoolMap returns the map[string]bool value of a given key path
	// or an empty map[string]bool if the path does not exist or if the
	// value is not a valid bool map.
	BoolMap(path string) map[string]bool

	// All returns all configs
	All() map[string]interface{}

	// Get returns interface{} value
	Get(path string) interface{}
}
