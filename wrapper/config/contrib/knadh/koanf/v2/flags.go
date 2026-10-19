package koanf

import (
	k "github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

// loadFlags merges the flags of fs into ko the way koanf's posflag provider
// does: a flag the command line did not set is merged only when no earlier
// source set its key, so its default never hides a file or env value.
func loadFlags(ko *k.Koanf, fs *flag.FlagSet) error {
	values := make(map[string]interface{})

	var err error
	fs.VisitAll(func(fl *flag.Flag) {
		if err != nil || (!fl.Changed && ko.Exists(fl.Name)) {
			return
		}
		values[fl.Name], err = flagValue(fs, fl)
	})
	if err != nil {
		return err
	}

	for key, value := range values {
		if err := ko.Set(key, value); err != nil {
			return err
		}
	}

	return nil
}

// flagGetters reads a flag by its pflag type name, returning the same Go
// types as koanf's posflag.FlagVal (signed ints as int64, unsigned as uint64).
var flagGetters = map[string]func(fs *flag.FlagSet, name string) (interface{}, error){
	"int":   func(fs *flag.FlagSet, n string) (interface{}, error) { v, err := fs.GetInt(n); return int64(v), err },
	"int8":  func(fs *flag.FlagSet, n string) (interface{}, error) { v, err := fs.GetInt8(n); return int64(v), err },
	"int16": func(fs *flag.FlagSet, n string) (interface{}, error) { v, err := fs.GetInt16(n); return int64(v), err },
	"int32": func(fs *flag.FlagSet, n string) (interface{}, error) { v, err := fs.GetInt32(n); return int64(v), err },
	"int64": func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetInt64(n) },
	"uint":  func(fs *flag.FlagSet, n string) (interface{}, error) { v, err := fs.GetUint(n); return uint64(v), err },
	"uint8": func(fs *flag.FlagSet, n string) (interface{}, error) { v, err := fs.GetUint8(n); return uint64(v), err },
	"uint16": func(fs *flag.FlagSet, n string) (interface{}, error) {
		v, err := fs.GetUint16(n)
		return uint64(v), err
	},
	"uint32": func(fs *flag.FlagSet, n string) (interface{}, error) {
		v, err := fs.GetUint32(n)
		return uint64(v), err
	},
	"uint64":         func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetUint64(n) },
	"float32":        func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetFloat32(n) },
	"float64":        func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetFloat64(n) },
	"bool":           func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetBool(n) },
	"duration":       func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetDuration(n) },
	"ip":             func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetIP(n) },
	"ipMask":         func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetIPv4Mask(n) },
	"bytesBase64":    func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetBytesBase64(n) },
	"string":         func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetString(n) },
	"stringSlice":    func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetStringSlice(n) },
	"intSlice":       func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetIntSlice(n) },
	"uintSlice":      func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetUintSlice(n) },
	"boolSlice":      func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetBoolSlice(n) },
	"durationSlice":  func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetDurationSlice(n) },
	"ipSlice":        func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetIPSlice(n) },
	"stringToString": func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetStringToString(n) },
	"stringToInt":    func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetStringToInt(n) },
	"stringToInt64":  func(fs *flag.FlagSet, n string) (interface{}, error) { return fs.GetStringToInt64(n) },
}

// flagValue returns the typed value of fl read from fs, or its string form
// for a type registerFlags never declares.
func flagValue(fs *flag.FlagSet, fl *flag.Flag) (interface{}, error) {
	get, ok := flagGetters[fl.Value.Type()]
	if !ok {
		return fl.Value.String(), nil
	}
	return get(fs, fl.Name)
}
