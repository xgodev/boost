package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	boost.Start()

	conn, err := mongo.NewConn(context.Background())
	if err != nil {
		log.Panic(err)
	}

	err = conn.Client.Ping(context.Background(), nil)
	if err != nil {
		log.Panic(err)
	}

}
