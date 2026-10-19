package main

import (
	"context"
	"encoding/json"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/health"
	"github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2"
	h "github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2/plugins/local/extra/health"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	boost.Start()

	session, err := gocql.NewSession(context.Background(), h.Register)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	all := health.CheckAll(context.Background())

	j, err := json.Marshal(all)
	if err != nil {
		panic(err)
	}

	log.Info(string(j))
}
