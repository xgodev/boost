package service

import (
	"github.com/xgodev/boost/bootstrap/function/adapter"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root   = adapter.Root + ".dapr.service.http"
	topics = root + ".topics"
)

func init() {
	config.Add(topics, []string{"changeme"}, "nats listener subjects")
}
