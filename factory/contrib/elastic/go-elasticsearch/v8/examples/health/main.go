package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost"
	h "github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8/plugins/local/extra/health"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	boost.Start()

	_, err := elasticsearch.NewClient(context.Background(), health.Register)
	if err != nil {
		log.Panic(err)
	}

	all := h.CheckAll(context.Background())

	j, _ := json.Marshal(all)

	log.Info(string(j))
}
