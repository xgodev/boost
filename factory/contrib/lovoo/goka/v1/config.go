package goka

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root    = "boost.factory.goka"
	brokers = ".brokers"
	topic   = ".topic"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+brokers, []string{"localhost:9092"}, "defines brokers addresses")
	config.Add(path+topic, "changeme", "defines topic name")
}
