package pubsub

import "github.com/xgodev/boost/config"

const (
	root     = "boost.factory.gocloud"
	resource = ".resource"
	tp       = ".type"
	region   = ".region"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+tp, "memory", "define queue type")
	config.Add(path+resource, "topicA", "define queue resource")
	config.Add(path+region, "", "define queue region")
}
