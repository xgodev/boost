package main

import (
	"context"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	boost.Start()

	session, err := gocql.NewSession(context.Background())
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var version string
	if err := session.Query("SELECT release_version FROM system.local").Scan(&version); err != nil {
		panic(err)
	}

	log.Infof("cassandra %s", version)
}
