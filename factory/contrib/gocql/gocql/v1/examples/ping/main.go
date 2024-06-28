package main

import (
	"context"
	"github.com/xgodev/boost/wrapper/config"

	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
)

func main() {

	config.Load()

	ilog.New()

	session, err := gocql.NewSession(context.Background())
	if err != nil {
		panic(err)
	}

	defer session.Close()

	err = session.Query("void").Exec()
	if err != nil {
		panic(err)
	}

}
