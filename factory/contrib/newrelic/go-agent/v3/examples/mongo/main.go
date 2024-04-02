package main

import (
	"context"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	newrelic "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1/plugins/contrib/newrelic/go-agent/v3"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
)

func main() {

	config.Load()
	ilog.New()

	mongo.NewConn(context.Background(), newrelic.Register)
}
