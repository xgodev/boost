package main

import (
	"context"
	"encoding/json"
	h "github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v7"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v7/plugins/local/extra/health"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
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
