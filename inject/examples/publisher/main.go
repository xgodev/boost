// Command publisher announces a created order through the publisher inject
// builds, on the driver boost.wrapper.publisher.driver.name chooses (noop
// by default).
//
// Regenerate the modules in gen/ with:
//
//	go generate ./inject/examples/publisher
package main

//go:generate go run github.com/xgodev/boost/inject/cmd/inject --boost.inject.output=gen

import (
	"context"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject/examples/publisher/domain/model/entity"
	"github.com/xgodev/boost/inject/examples/publisher/domain/publisher/ordercreated"
	"github.com/xgodev/boost/inject/examples/publisher/gen"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

func main() {
	boost.Start()

	var orders ordercreated.Publisher
	if err := fx.New(gen.Module(), fx.Populate(&orders), fx.NopLogger).Err(); err != nil {
		log.Fatalf("%v", err)
	}

	if err := orders.Publish(context.Background(), entity.Order{ID: "1", Item: "book"}); err != nil {
		log.Fatalf("%v", err)
	}
	log.Info("order 1 published")
}
