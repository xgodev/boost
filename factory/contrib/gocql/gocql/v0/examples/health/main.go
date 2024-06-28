package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v0"
	"github.com/xgodev/boost/wrapper/config"

	"github.com/xgodev/boost/extra/health"
	h "github.com/xgodev/boost/factory/contrib/gocql/gocql/v0/plugins/local/extra/health"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
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
