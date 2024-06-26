package koanf

import (
	"github.com/xgodev/boost/wrapper/config"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		key         string
		value       interface{}
		description string
		gotFunc     func(string) interface{}
		expected    interface{}
	}{
		{
			key:         "green",
			value:       map[string]string{"a": "A"},
			description: "test map string string",
			gotFunc: func(key string) interface{} {
				return config.StringMap(key)
			},
			expected: map[string]string{"a": "A"},
		},
		{
			key:         "blue",
			value:       "test",
			description: "test",
			gotFunc: func(key string) interface{} {
				return config.String(key)
			},
			expected: "test",
		},
		{
			key:         "red",
			value:       map[string]string{"h": "0.0.0.0"},
			description: "overriding default",
			gotFunc: func(key string) interface{} {
				return config.StringMap("red")
			},
			expected: map[string]string{"h": "127.0.0.14"},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			os.Args = []string{"--conf", "./testdata/config.yaml"}
			config.Set(New())
			Load([]config.Config{{Key: tt.key, Value: tt.value, Description: tt.description}})
			got := tt.gotFunc(tt.key)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v got %v", tt.expected, got)
			}
		})
	}
}
