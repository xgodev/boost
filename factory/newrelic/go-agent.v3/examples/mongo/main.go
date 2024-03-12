package main

import (
	"context"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/go.mongodb.org/mongo-driver.v1"
	newrelic "github.com/xgodev/boost/factory/go.mongodb.org/mongo-driver.v1/plugins/contrib/newrelic/go-agent.v3"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
)

func main() {

	config.Load()
	ilog.New()

	mongo.NewConn(context.Background(), newrelic.Register)
}
