package koanf

import (
	"fmt"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/model"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"log"

	"github.com/gobeam/stringy"
	"github.com/knadh/koanf/providers/confmap"
	flag "github.com/spf13/pflag"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
)

var (
	instance *koanf.Koanf
	f        *flag.FlagSet
)

func init() {
	instance = koanf.New(".")
	flagLoad()
}

func flagLoad() {
	// Use the POSIX compliant pflag lib instead of Go's flag lib.
	f = flag.NewFlagSet("config", flag.ContinueOnError)

	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
}

// Load parsing and load flags, files and environments.
func Load(cfs []config.Config) {

	// Load flags
	parseFlags(cfs)

	m := make(map[string]interface{})

	for _, v := range cfs {

		switch v.Value.(type) {

		case map[string]string:
			toString, err := f.GetStringToString(v.Key)
			if err != nil {
				log.Println(err)
			}
			m[v.Key] = toString
		case map[string]int:
			toString, err := f.GetStringToInt(v.Key)
			if err != nil {
				log.Println(err)
			}
			m[v.Key] = toString
		case map[string]int64:
			toString, err := f.GetStringToInt64(v.Key)
			if err != nil {
				log.Println(err)
			}
			m[v.Key] = toString
		}

	}

	if err := instance.Load(confmap.Provider(m, "."), nil); err != nil {
		panic(err)
	}

	var files []string

	confEnv := os.Getenv(model.ConfEnvironment)
	if confEnv != "" {
		// Load the config files provided in the environment var.
		files = strings.Split(confEnv, ",")
	} else {
		// Load the config files provided in the commandline.
		files, _ = f.GetStringSlice(model.ConfArgument)
	}

	for _, c := range files {

		var parser koanf.Parser

		if filepath.Ext(c) == ".toml" {
			parser = toml.Parser()
		} else if filepath.Ext(c) == ".yaml" || filepath.Ext(c) == ".yml" {
			parser = yaml.Parser()
		} else if filepath.Ext(c) == ".json" {
			parser = json.Parser()
		} else {
			panic(fmt.Sprintf("error on check extension of file %s", c))
		}

		if err := instance.Load(file.Provider(c), parser); err != nil {
			panic(err)
		}
	}

	// Env vars
	err := instance.Load(env.Provider("", ".", func(s string) string {
		return parseEnv(s)
	}), nil)
	if err != nil {
		panic(err)
	}

	// Load flags
	flap := posflag.Provider(f, ".", instance)

	if err := instance.Load(flap, nil); err != nil {
		panic(err)
	}

}

func parseFlags(cfs []config.Config) {

	for _, v := range cfs {

		fl := f.Lookup(v.Key)
		if fl != nil {
			continue
		}

		switch t := v.Value.(type) {

		case string:
			f.String(v.Key, t, v.Description)
		case []string:
			f.StringSlice(v.Key, t, v.Description)
		case map[string]string:
			f.StringToString(v.Key, t, v.Description)
		case map[string]int:
			f.StringToInt(v.Key, t, v.Description)
		case map[string]int64:
			f.StringToInt64(v.Key, t, v.Description)
		case bool:
			f.Bool(v.Key, t, v.Description)
		case []bool:
			f.BoolSlice(v.Key, t, v.Description)
		case []int:
			f.IntSlice(v.Key, t, v.Description)
		case int:
			f.Int(v.Key, t, v.Description)
		case int64:
			f.Int64(v.Key, t, v.Description)
		case int32:
			f.Int32(v.Key, t, v.Description)
		case int16:
			f.Int16(v.Key, t, v.Description)
		case int8:
			f.Int8(v.Key, t, v.Description)
		case uint:
			f.Uint(v.Key, t, v.Description)
		case []uint:
			f.UintSlice(v.Key, t, v.Description)
		case uint64:
			f.Uint64(v.Key, t, v.Description)
		case uint32:
			f.Uint32(v.Key, t, v.Description)
		case uint16:
			f.Uint16(v.Key, t, v.Description)
		case uint8:
			f.Uint8(v.Key, t, v.Description)
		case time.Duration:
			f.Duration(v.Key, t, v.Description)
		case []time.Duration:
			f.DurationSlice(v.Key, t, v.Description)
		case []byte:
			f.BytesBase64(v.Key, t, v.Description)
		case float32:
			f.Float32(v.Key, t, v.Description)
		case float64:
			f.Float64(v.Key, t, v.Description)
		case net.IP:
			f.IP(v.Key, t, v.Description)
		case []net.IP:
			f.IPSlice(v.Key, t, v.Description)
		case net.IPMask:
			f.IPMask(v.Key, t, v.Description)
		}

	}

	flc := f.Lookup(model.ConfArgument)
	if flc == nil {
		// Path to one or more config files to load into koanf along with some config params.
		f.StringSlice(model.ConfArgument, nil, "path to one or more config files")
	}

	err := f.Parse(os.Args[0:])
	if err != nil {
		log.Println(err)
	}
}

func parseEnv(s string) string {

	strs := make([]string, 0)

	s = strings.ReplaceAll(s, "__", "-")

	for _, v := range strings.Split(s, "_") {

		var add string

		if strings.Contains(v, "-") {

			sgyl := stringy.New(strings.ToLower(v))
			sgylc := stringy.New(sgyl.CamelCase().Get())
			add = sgylc.LcFirst()

		} else {

			add = strings.ToLower(v)

		}

		strs = append(strs, add)

	}

	return strings.Join(strs, ".")
}
