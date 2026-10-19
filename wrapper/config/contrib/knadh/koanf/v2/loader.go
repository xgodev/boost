package koanf

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gobeam/stringy"
	"github.com/knadh/koanf/providers/env/v2"
	k "github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/config/model"
)

var (
	instance *k.Koanf
	f        *flag.FlagSet
)

func init() {
	instance = k.New(".")
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

// Load registers cfs as flags, parses the command line and loads every config
// source into the package instance. It panics when a source fails to load.
func Load(cfs []config.Config) {
	parseFlags(cfs)

	if err := load(instance, f, cfs); err != nil {
		panic(err)
	}
}

// load merges into ko the map defaults of cfs, the config files, the
// environment and the flags of fs, each one overriding the previous.
func load(ko *k.Koanf, fs *flag.FlagSet, cfs []config.Config) error {
	if err := loadDefaults(ko, fs, cfs); err != nil {
		return err
	}

	if err := loadFiles(ko, confFiles(fs)); err != nil {
		return err
	}

	if err := loadEnv(ko); err != nil {
		return err
	}

	return loadFlags(ko, fs)
}

// loadDefaults merges the map defaults first, so a file or env var that sets
// one entry of the map keeps the others.
func loadDefaults(ko *k.Koanf, fs *flag.FlagSet, cfs []config.Config) error {
	for _, v := range cfs {

		var (
			value interface{}
			err   error
		)

		switch v.Value.(type) {
		case map[string]string:
			value, err = fs.GetStringToString(v.Key)
		case map[string]int:
			value, err = fs.GetStringToInt(v.Key)
		case map[string]int64:
			value, err = fs.GetStringToInt64(v.Key)
		default:
			continue
		}

		if err != nil {
			log.Println(err)
		}

		if err := ko.Set(v.Key, value); err != nil {
			return err
		}
	}

	return nil
}

// confFiles returns the config files listed in the CONF environment variable
// or, when it is empty, in the --conf flag.
func confFiles(fs *flag.FlagSet) []string {
	if confEnv := os.Getenv(model.ConfEnvironment); confEnv != "" {
		return strings.Split(confEnv, ",")
	}

	files, err := fs.GetStringSlice(model.ConfArgument)
	if err != nil {
		log.Println(err)
	}

	return files
}

func loadFiles(ko *k.Koanf, files []string) error {
	for _, c := range files {

		parser, err := parserFor(c)
		if err != nil {
			return err
		}

		if err := ko.Load(fileProvider(c), parser); err != nil {
			return err
		}
	}

	return nil
}

func loadEnv(ko *k.Koanf) error {
	return ko.Load(env.Provider(".", env.Opt{
		TransformFunc: func(key, value string) (string, any) {
			return parseEnv(key), value
		},
	}), nil)
}

func parseFlags(cfs []config.Config) {

	registerFlags(f, cfs)

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

// registerFlags declares one flag per config, typed by its default value.
// A key already declared keeps its first flag.
func registerFlags(fs *flag.FlagSet, cfs []config.Config) {

	for _, v := range cfs {

		fl := fs.Lookup(v.Key)
		if fl != nil {
			continue
		}

		switch t := v.Value.(type) {

		case string:
			fs.String(v.Key, t, v.Description)
		case []string:
			fs.StringSlice(v.Key, t, v.Description)
		case map[string]string:
			fs.StringToString(v.Key, t, v.Description)
		case map[string]int:
			fs.StringToInt(v.Key, t, v.Description)
		case map[string]int64:
			fs.StringToInt64(v.Key, t, v.Description)
		case bool:
			fs.Bool(v.Key, t, v.Description)
		case []bool:
			fs.BoolSlice(v.Key, t, v.Description)
		case []int:
			fs.IntSlice(v.Key, t, v.Description)
		case int:
			fs.Int(v.Key, t, v.Description)
		case int64:
			fs.Int64(v.Key, t, v.Description)
		case int32:
			fs.Int32(v.Key, t, v.Description)
		case int16:
			fs.Int16(v.Key, t, v.Description)
		case int8:
			fs.Int8(v.Key, t, v.Description)
		case uint:
			fs.Uint(v.Key, t, v.Description)
		case []uint:
			fs.UintSlice(v.Key, t, v.Description)
		case uint64:
			fs.Uint64(v.Key, t, v.Description)
		case uint32:
			fs.Uint32(v.Key, t, v.Description)
		case uint16:
			fs.Uint16(v.Key, t, v.Description)
		case uint8:
			fs.Uint8(v.Key, t, v.Description)
		case time.Duration:
			fs.Duration(v.Key, t, v.Description)
		case []time.Duration:
			fs.DurationSlice(v.Key, t, v.Description)
		case []byte:
			fs.BytesBase64(v.Key, t, v.Description)
		case float32:
			fs.Float32(v.Key, t, v.Description)
		case float64:
			fs.Float64(v.Key, t, v.Description)
		case net.IP:
			fs.IP(v.Key, t, v.Description)
		case []net.IP:
			fs.IPSlice(v.Key, t, v.Description)
		case net.IPMask:
			fs.IPMask(v.Key, t, v.Description)
		}

	}
}

// parseEnv turns an env var name into a config key: "_" splits the path and
// "__" joins camel case words, so K_CAMEL__CASE becomes k.camelCase.
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
