package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

type DBConfig struct {
	Username string `config:"username"`
}

type MyConfig struct {
	Addr  string
	Debug string
	DB    DBConfig `config:"db"`
	Redis struct {
		Host string `config:"h"`
	} `config:"red"`
}

func (s *WrapperSuite) TestConfigMethods() {

	tt := []struct {
		name string
		init func()
		got  func() interface{}
		want interface{}
	}{
		{
			name: "Pflag config returns the value (string) from the key",
			init: func() {
				flagLoad()
				Add("key", "value", "test")
			},
			got:  func() interface{} { return instance.String("key") },
			want: "value",
		},
		{
			name: "Env config returns the value (string) from the key",
			init: func() {
				os.Setenv("K_ENV", "value")
			},
			got:  func() interface{} { return instance.String("k.env") },
			want: "value",
		},
		{
			name: "Env config returns the camel case value (string) from the key",
			init: func() {
				os.Setenv("K_CAMEL-CASE", "value")
			},
			got:  func() interface{} { return instance.String("k.camelCase") },
			want: "value",
		},
		{
			name: "Env config returns the camel case two value (string) from the key",
			init: func() {
				os.Setenv("K_CAMEL__CASE__Two", "value")
			},
			got:  func() interface{} { return instance.String("k.camelCaseTwo") },
			want: "value",
		},
		{
			name: "Conf config returns the value (bool) from the key",
			init: func() {
				flagLoad()
				os.Args = []string{"--conf", "./testdata/config.json", "--conf", "./testdata/config.yaml"}
			},
			got:  func() interface{} { return instance.Bool("debug") },
			want: true,
		},
		{
			name: "Conf config returns the value (string) from the key",
			init: func() {
				flagLoad()
				os.Args = []string{"--conf", "./testdata/config.json", "--conf", "./testdata/config.yaml"}
			},
			got:  func() interface{} { return instance.String("redis.host") },
			want: "127.0.0.13",
		},
		{
			name: "Unmarshal config doesn't returns error",
			init: func() {
				flagLoad()
				os.Args = []string{"--conf", "./testdata/config.json", "--conf", "./testdata/config.yaml"}
			},
			got: func() interface{} {
				c := MyConfig{}
				return Unmarshal(&c)
			},
			want: nil,
		},
		{
			name: "Unmarshal config returns the Addr field value",
			init: func() {
				flagLoad()
				os.Args = []string{"--conf", "./testdata/config.json", "--conf", "./testdata/config.yaml"}
			},
			got: func() interface{} {
				c := MyConfig{}
				Unmarshal(&c)
				return c.Addr
			},
			want: ":8083",
		},
		{
			name: "Unmarshal config returns the DB.Username field value",
			init: func() {
				flagLoad()
				os.Args = []string{"--conf", "./testdata/config.json", "--conf", "./testdata/config.yaml"}
			},
			got: func() interface{} {
				c := MyConfig{}
				Unmarshal(&c)
				return c.DB.Username
			},
			want: "foosss",
		},
		{
			name: "Unmarshal config returns the Redis.Host field value",
			init: func() {
				flagLoad()
				os.Args = []string{"--conf", "./testdata/config.json", "--conf", "./testdata/config.yaml"}
			},
			got: func() interface{} {
				c := MyConfig{}
				Unmarshal(&c)
				return c.Redis.Host
			},
			want: "127.0.0.14",
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			t.init()
			Load()
			got := t.got()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}
