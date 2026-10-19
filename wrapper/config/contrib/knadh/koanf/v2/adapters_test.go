package koanf

import (
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	k "github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/model"
)

// strict returns a koanf instance that fails a merge changing the type of
// key, so each load step's error path can be reached.
func strict(t *testing.T, key string) *k.Koanf {
	t.Helper()
	ko := k.NewWithConf(k.Conf{Delim: ".", StrictMerge: true})
	if err := ko.Set(key, map[string]interface{}{"a": "b"}); err != nil {
		t.Fatal(err)
	}
	return ko
}

func TestLoadReturnsEachSourceError(t *testing.T) {
	t.Setenv(model.ConfEnvironment, "")
	t.Setenv("ZZSTRICT", "scalar")

	tests := []struct {
		name string
		ko   func(t *testing.T) *k.Koanf
		fs   func() *flag.FlagSet
		cfs  []config.Config
	}{
		{
			name: "default map",
			ko: func(t *testing.T) *k.Koanf {
				ko := k.NewWithConf(k.Conf{Delim: ".", StrictMerge: true})
				if err := ko.Set("m", "scalar"); err != nil {
					t.Fatal(err)
				}
				return ko
			},
			fs: func() *flag.FlagSet {
				fs := flag.NewFlagSet("t", flag.ContinueOnError)
				fs.StringToString("m", map[string]string{"a": "A"}, "")
				fs.StringSlice(model.ConfArgument, nil, "")
				return fs
			},
			cfs: []config.Config{{Key: "m", Value: map[string]string{"a": "A"}}},
		},
		{
			name: "file",
			ko:   func(t *testing.T) *k.Koanf { return k.New(".") },
			fs: func() *flag.FlagSet {
				fs := flag.NewFlagSet("t", flag.ContinueOnError)
				fs.StringSlice(model.ConfArgument, []string{"./testdata/config.ini"}, "")
				return fs
			},
		},
		{
			name: "env",
			ko:   func(t *testing.T) *k.Koanf { return strict(t, "zzstrict") },
			fs: func() *flag.FlagSet {
				fs := flag.NewFlagSet("t", flag.ContinueOnError)
				fs.StringSlice(model.ConfArgument, nil, "")
				return fs
			},
		},
		{
			name: "flag",
			ko:   func(t *testing.T) *k.Koanf { return strict(t, "zzflag") },
			fs: func() *flag.FlagSet {
				fs := flag.NewFlagSet("t", flag.ContinueOnError)
				fs.StringSlice(model.ConfArgument, nil, "")
				fs.String("zzflag", "", "")
				if err := fs.Parse([]string{"--zzflag", "x"}); err != nil {
					panic(err)
				}
				return fs
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ko := tt.ko(t)
			if tt.name != "env" {
				// only the env case keeps the conflicting ZZSTRICT key
				t.Setenv("ZZSTRICT", "")
			}
			if err := load(ko, tt.fs(), tt.cfs); err == nil {
				t.Errorf("expected the %s step to fail", tt.name)
			}
		})
	}
}

func TestLoadFlagsKeepsEarlierSourcesForUnsetFlags(t *testing.T) {
	ko := k.New(".")
	if err := ko.Set("from.file", "file"); err != nil {
		t.Fatal(err)
	}
	fs := flag.NewFlagSet("t", flag.ContinueOnError)
	fs.String("from.file", "default", "")
	fs.String("only.default", "default", "")
	fs.String("from.cli", "default", "")
	if err := fs.Parse([]string{"--from.cli", "cli"}); err != nil {
		t.Fatal(err)
	}

	if err := loadFlags(ko, fs); err != nil {
		t.Fatal(err)
	}

	for key, want := range map[string]string{"from.file": "file", "only.default": "default", "from.cli": "cli"} {
		if got := ko.String(key); got != want {
			t.Errorf("%s: expected %q got %q", key, want, got)
		}
	}
}

// custom is a pflag.Value of a type flagValue does not know.
type custom string

func (c *custom) String() string     { return string(*c) }
func (c *custom) Set(s string) error { *c = custom(s); return nil }
func (c *custom) Type() string       { return "custom" }

func TestFlagValueTypes(t *testing.T) {
	fs := flag.NewFlagSet("t", flag.ContinueOnError)
	cfs := []config.Config{
		{Key: "string", Value: "s"},
		{Key: "stringSlice", Value: []string{"a"}},
		{Key: "stringToString", Value: map[string]string{"a": "A"}},
		{Key: "stringToInt", Value: map[string]int{"a": 1}},
		{Key: "stringToInt64", Value: map[string]int64{"a": 1}},
		{Key: "bool", Value: true},
		{Key: "boolSlice", Value: []bool{true}},
		{Key: "intSlice", Value: []int{1}},
		{Key: "int", Value: 1},
		{Key: "int64", Value: int64(1)},
		{Key: "int32", Value: int32(1)},
		{Key: "int16", Value: int16(1)},
		{Key: "int8", Value: int8(1)},
		{Key: "uint", Value: uint(1)},
		{Key: "uintSlice", Value: []uint{1}},
		{Key: "uint64", Value: uint64(1)},
		{Key: "uint32", Value: uint32(1)},
		{Key: "uint16", Value: uint16(1)},
		{Key: "uint8", Value: uint8(1)},
		{Key: "duration", Value: time.Second},
		{Key: "durationSlice", Value: []time.Duration{time.Second}},
		{Key: "bytesBase64", Value: []byte("b")},
		{Key: "float32", Value: float32(1.5)},
		{Key: "float64", Value: 1.5},
		{Key: "ip", Value: net.IPv4(1, 2, 3, 4)},
		{Key: "ipSlice", Value: []net.IP{net.IPv4(1, 2, 3, 4)}},
		{Key: "ipMask", Value: net.IPv4Mask(255, 255, 255, 0)},
	}
	registerFlags(fs, cfs)
	var c custom = "raw"
	fs.Var(&c, "custom", "")

	want := map[string]interface{}{
		"string":         "s",
		"stringSlice":    []string{"a"},
		"stringToString": map[string]string{"a": "A"},
		"stringToInt":    map[string]int{"a": 1},
		"stringToInt64":  map[string]int64{"a": 1},
		"bool":           true,
		"boolSlice":      []bool{true},
		"intSlice":       []int{1},
		"int":            int64(1),
		"int64":          int64(1),
		"int32":          int64(1),
		"int16":          int64(1),
		"int8":           int64(1),
		"uint":           uint64(1),
		"uintSlice":      []uint{1},
		"uint64":         uint64(1),
		"uint32":         uint64(1),
		"uint16":         uint64(1),
		"uint8":          uint64(1),
		"duration":       time.Second,
		"durationSlice":  []time.Duration{time.Second},
		"bytesBase64":    []byte("b"),
		"float32":        float32(1.5),
		"float64":        1.5,
		"ip":             net.IPv4(1, 2, 3, 4),
		"ipSlice":        []net.IP{net.IPv4(1, 2, 3, 4)},
		"ipMask":         net.IPv4Mask(255, 255, 255, 0),
		"custom":         "raw",
	}

	for key, w := range want {
		t.Run(key, func(t *testing.T) {
			got, err := flagValue(fs, fs.Lookup(key))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, w) {
				t.Errorf("got %#v want %#v", got, w)
			}
		})
	}
}

func TestFlagValueFailsOnAFlagOfAnotherType(t *testing.T) {
	fs := flag.NewFlagSet("t", flag.ContinueOnError)
	fs.String("x", "", "")
	other := flag.NewFlagSet("o", flag.ContinueOnError)
	other.Int("x", 0, "")

	if _, err := flagValue(fs, other.Lookup("x")); err == nil {
		t.Error("expected an error reading a string flag as int")
	}
}

func TestFileProviderDoesNotRead(t *testing.T) {
	if _, err := fileProvider("./testdata/config.yaml").Read(); err == nil {
		t.Error("expected Read to be unsupported")
	}
}

func TestParsers(t *testing.T) {
	for _, ext := range []string{".json", ".yaml", ".yml", ".toml"} {
		t.Run(ext, func(t *testing.T) {
			p, err := parserFor("config" + ext)
			if err != nil {
				t.Fatal(err)
			}
			b, err := p.Marshal(map[string]interface{}{"a": map[string]interface{}{"b": "c"}})
			if err != nil {
				t.Fatal(err)
			}
			got, err := p.Unmarshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if v := got["a"].(map[string]interface{})["b"]; v != "c" {
				t.Errorf("round trip lost a.b: %#v", got)
			}
			if _, err := p.Unmarshal([]byte("{[")); err == nil {
				t.Error("expected an error on invalid input")
			}
		})
	}

	if _, err := parserFor("config.ini"); err == nil || !strings.Contains(err.Error(), "config.ini") {
		t.Errorf("expected an error naming config.ini, got %v", err)
	}
}
