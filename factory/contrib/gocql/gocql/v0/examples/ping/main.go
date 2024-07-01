package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/gocql/gocql/v0"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
)

func main() {

	boost.Start()

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
