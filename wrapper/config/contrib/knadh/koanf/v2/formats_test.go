package koanf

import (
	"os"
	"reflect"
	"testing"

	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/model"
)

func TestLoadTOMLFileFromArgument(t *testing.T) {
	flagLoad()
	os.Args = []string{"--conf", "./testdata/config.toml"}
	config.Set(New())
	config.Load()

	if got := instance.String("toml.host"); got != "127.0.0.20" {
		t.Errorf("expected 127.0.0.20 got %q", got)
	}
}

func TestLoadFilesFromConfEnvironment(t *testing.T) {
	t.Setenv(model.ConfEnvironment, "./testdata/config.toml,./testdata/config.json")
	flagLoad()
	os.Args = []string{}
	config.Set(New())
	config.Load()

	if got := instance.String("toml.host"); got != "127.0.0.20" {
		t.Errorf("expected toml.host 127.0.0.20 got %q", got)
	}
	if got := instance.String("db.username"); got != "foosss" {
		t.Errorf("expected db.username foosss got %q", got)
	}
}

func TestLoadPanicsOnUnknownFileExtension(t *testing.T) {
	flagLoad()
	os.Args = []string{"--conf", "./testdata/config.ini"}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic for a .ini file")
		}
	}()
	Load(nil)
}

func TestLoadPanicsOnMissingFile(t *testing.T) {
	flagLoad()
	os.Args = []string{"--conf", "./testdata/missing.yaml"}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic for a missing file")
		}
	}()
	Load(nil)
}

func TestLoadIntMapDefaults(t *testing.T) {
	flagLoad()
	os.Args = []string{}
	config.Set(New())
	Load([]config.Config{
		{Key: "maps.int", Value: map[string]int{"a": 1}, Description: "int map"},
		{Key: "maps.int64", Value: map[string]int64{"b": 2}, Description: "int64 map"},
	})

	if got, want := config.IntMap("maps.int"), map[string]int{"a": 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v got %v", want, got)
	}
	if got, want := config.Int64Map("maps.int64"), map[string]int64{"b": 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v got %v", want, got)
	}
}
