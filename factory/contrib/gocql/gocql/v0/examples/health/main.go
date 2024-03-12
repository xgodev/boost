package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v0"

	"github.com/xgodev/boost/config"
	h "github.com/xgodev/boost/factory/contrib/gocql/gocql/v0/plugins/local/health"
	ilog "github.com/xgodev/boost/factory/local/log"
	"github.com/xgodev/boost/health"
	"github.com/xgodev/boost/log"
)

func main() {

	config.Load()

	ilog.New()

	session, err := gocql.NewSession(context.Background(), h.Register)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	all := health.CheckAll(context.Background())

	j, _ := json.Marshal(all)

	log.Info(string(j))

}
