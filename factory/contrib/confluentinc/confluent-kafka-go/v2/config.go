package confluent

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root            = "boost.factory.confluent"
	brokers         = ".brokers"
	consumer        = ".consumer"
	topics          = consumer + ".topics"
	groupId         = consumer + ".groupId"
	autoOffsetReset = consumer + ".autoOffsetReset"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+brokers, "localhost:9092", "defines brokers addresses")
	config.Add(path+topics, []string{"changeme"}, "defines topics")
	config.Add(path+groupId, "changeme", "defines consumer groupid")
	config.Add(path+autoOffsetReset, "earliest", "defines consumer auto offset reset")
}
