package main

import (
	"context"
	"encoding/json"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/go-redis/redis.v7"
	"github.com/xgodev/boost/factory/go-redis/redis.v7/plugins/contrib/americanas-go/health.v1"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
	h "github.com/xgodev/boost/health"
	"github.com/xgodev/boost/log"
)

func main() {

	config.Load()

	ilog.New()

	var err error

	healthIntegrator := health.NewClusterHealth()

	_, err = redis.NewClusterClient(context.Background(), healthIntegrator.Register)
	if err != nil {
		log.Error(err)
	}

	all := h.CheckAll(context.Background())

	j, _ := json.Marshal(all)

	log.Info(string(j))

}
