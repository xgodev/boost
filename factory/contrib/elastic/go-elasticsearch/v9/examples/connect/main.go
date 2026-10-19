package main

import (
	"context"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v9"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	boost.Start()

	client, err := elasticsearch.NewClient(context.Background())
	if err != nil {
		log.Panic(err)
	}

	info, err := client.Info()
	if err != nil {
		log.Panic(err)
	}
	defer info.Body.Close()

	log.Infof("status: %v", info.StatusCode)
}
