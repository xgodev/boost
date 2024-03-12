package main

import (
	"context"
	"encoding/json"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/gocql/gocql.v0"
	h "github.com/xgodev/boost/factory/gocql/gocql.v0/plugins/contrib/americanas-go/health.v1"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
	"github.com/xgodev/boost/health"
	"github.com/xgodev/boost/log"
)

func main() {

	config.Load()

	ilog.New()

	i := h.NewHealth()

	session, err := gocql.NewSession(context.Background(), i.Register)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	all := health.CheckAll(context.Background())

	j, _ := json.Marshal(all)

	log.Info(string(j))

}
