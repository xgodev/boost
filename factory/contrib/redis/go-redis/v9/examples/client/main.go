package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost"
	h "github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9/plugins/local/extra/health"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	boost.Start()

	var err error

	healthIntegrator, err := health.NewClientHealth()
	if err != nil {
		log.Fatalf(err.Error())
	}

	_, err = redis.NewClient(context.Background(), healthIntegrator.Register)
	if err != nil {
		log.Error(err)
	}

	all := h.CheckAll(context.Background())

	j, _ := json.Marshal(all)

	log.Info(string(j))

}
