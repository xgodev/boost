package main

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v1"
	"github.com/xgodev/boost/wrapper/config"
)

func main() {

	config.Load()

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
