package nats

import (
	cloudevents "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root     = cloudevents.Root + ".nats"
	subjects = root + ".subjects"
	queue    = root + ".queue"
)

func init() {
	config.Add(subjects, []string{"changeme"}, "nats listener subjects")
	config.Add(queue, "changeme", "nats listener queue")
}
