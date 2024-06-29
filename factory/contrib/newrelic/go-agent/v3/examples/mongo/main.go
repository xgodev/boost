package main

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	newrelic "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1/plugins/contrib/newrelic/go-agent/v3"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/config"
)

func main() {

	config.Load()
	ilog.New()

	mongo.NewConn(context.Background(), newrelic.Register)
}
