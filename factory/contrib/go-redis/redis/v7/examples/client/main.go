package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/factory/contrib/go-redis/redis/v7"

	"github.com/xgodev/boost/config"
	h "github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/factory/contrib/go-redis/redis/v7/plugins/local/extra/health"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	config.Load()

	ilog.New()

	var err error

	healthIntegrator := health.NewClientHealth()

	_, err = redis.NewClient(context.Background(), healthIntegrator.Register)
	if err != nil {
		log.Error(err)
	}

	all := h.CheckAll(context.Background())

	j, _ := json.Marshal(all)

	log.Info(string(j))

}
