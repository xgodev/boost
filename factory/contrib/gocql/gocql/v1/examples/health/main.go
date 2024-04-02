package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1"

	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/extra/health"
	h "github.com/xgodev/boost/factory/contrib/gocql/gocql/v1/plugins/local/extra/health"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
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
