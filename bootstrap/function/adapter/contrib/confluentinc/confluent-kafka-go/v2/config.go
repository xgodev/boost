package confluent

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

const (
	root         = adapter.Root + ".kafka_confluent"
	topics       = root + ".topics"
	timeOut      = root + ".timeOut"
	manualCommit = root + ".manualCommit"
)

func init() {
	config.Add(topics, []string{"changeme"}, "kafka topics")
	config.Add(timeOut, 2*time.Second, "kafka timeout")
	config.Add(manualCommit, true, "kafka manual commit")
}
