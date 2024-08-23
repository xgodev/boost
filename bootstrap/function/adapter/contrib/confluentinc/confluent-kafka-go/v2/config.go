package confluent

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

const (
	root    = adapter.Root + ".kafka_confluent"
	topics  = root + ".topics"
	timeOut = root + ".timeOut"
)

func init() {
	config.Add(topics, []string{"changeme"}, "nats listener subjects")
	config.Add(timeOut, 2*time.Second, "nats listener subjects")
}
