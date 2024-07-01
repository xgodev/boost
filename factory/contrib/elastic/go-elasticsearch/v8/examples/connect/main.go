package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	boost.Start()

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
