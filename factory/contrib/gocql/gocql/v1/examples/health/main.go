package main

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1"
	h "github.com/xgodev/boost/factory/contrib/gocql/gocql/v1/plugins/local/extra/health"
	"github.com/xgodev/boost/wrapper/config"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	config.Load()

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
