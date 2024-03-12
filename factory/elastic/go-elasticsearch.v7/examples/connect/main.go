package main

import (
	"context"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/elastic/go-elasticsearch.v7"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
	"github.com/xgodev/boost/log"
)

func main() {

	config.Load()

	ilog.New()

	client, err := elasticsearch.NewClient(context.Background())
	if err != nil {
		log.Panic(err)
	}

	ping, err := client.Ping(client.Ping.WithPretty())
	if err != nil {
		log.Panic(err)
	}

	log.Infof("status: %v", ping.StatusCode)

}
