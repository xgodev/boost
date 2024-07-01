package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	newrelic "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1/plugins/contrib/newrelic/go-agent/v3"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	boost.Start()

	ilog.New()

	conn, err := mongo.NewConn(context.Background(), newrelic.Register)
	if err != nil {
		log.Panic(err)
	}

	err = conn.Client.Ping(context.Background(), nil)
	if err != nil {
		log.Panic(err)
	}

}
