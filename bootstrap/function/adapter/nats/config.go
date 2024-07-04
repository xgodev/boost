package nats

import (
	"github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root     = function.Root + ".nats"
	subjects = root + ".subjects"
	queue    = root + ".queue"
)

func init() {
	config.Add(subjects, []string{"changeme"}, "nats listener subjects")
	config.Add(queue, "changeme", "nats listener queue")
}
