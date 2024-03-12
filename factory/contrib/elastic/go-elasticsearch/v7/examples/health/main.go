package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v7/plugins/local/health"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v7"
	ilog "github.com/xgodev/boost/factory/local/log"
	h "github.com/xgodev/boost/health"
	"github.com/xgodev/boost/log"
)

func main() {

	config.Load()

	ilog.New()

	_, err := elasticsearch.NewClient(context.Background(), health.Register)
	if err != nil {
		log.Panic(err)
	}

	all := h.CheckAll(context.Background())

	j, _ := json.Marshal(all)

	log.Info(string(j))
}
