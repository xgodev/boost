package koanf

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

func TestProviderGetters(t *testing.T) {
	flagLoad()
	os.Args = []string{"--conf", "./testdata/getters.yaml"}
	p := New()
	config.Set(p)
	config.Load()

	var db DBConfig
	var whole struct {
		Getters struct {
			String string `config:"string"`
		} `config:"getters"`
	}

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"Exists", p.Exists("getters.int"), true},
		{"Exists missing", p.Exists("getters.missing"), false},
		{"Int64", p.Int64("getters.int"), int64(42)},
		{"Int64s", p.Int64s("getters.ints"), []int64{1, 2}},
		{"Int64Map", p.Int64Map("getters.int64map"), map[string]int64{"a": 2}},
		{"Int", p.Int("getters.int"), 42},
		{"Ints", p.Ints("getters.ints"), []int{1, 2}},
		{"IntMap", p.IntMap("getters.intmap"), map[string]int{"a": 1}},
		{"Float64", p.Float64("getters.float"), 1.5},
		{"Float64s", p.Float64s("getters.floats"), []float64{1.5, 2.5}},
		{"Float64Map", p.Float64Map("getters.floatmap"), map[string]float64{"a": 1.5}},
		{"Duration", p.Duration("getters.duration"), 3 * time.Second},
		{"Time", p.Time("getters.time", "2006-01-02"), time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)},
		{"String", p.String("getters.string"), "text"},
		{"Strings", p.Strings("getters.strings"), []string{"a", "b"}},
		{"StringMap", p.StringMap("getters.stringmap"), map[string]string{"a": "A"}},
		{"Bytes", p.Bytes("getters.string"), []byte("text")},
		{"Bool", p.Bool("getters.bool"), true},
		{"Bools", p.Bools("getters.bools"), []bool{true, false}},
		{"BoolMap", p.BoolMap("getters.boolmap"), map[string]bool{"a": true}},
		{"Get", p.Get("getters.string"), "text"},
		{"All", p.All()["getters.string"], "text"},
		{"UnmarshalWithPath error", p.UnmarshalWithPath("getters.db", &db), nil},
		{"UnmarshalWithPath value", db.Username, "getter"},
		{"Unmarshal error", p.Unmarshal(&whole), nil},
		{"Unmarshal value", whole.Getters.String, "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %#v want %#v", tt.got, tt.want)
			}
		})
	}
}
