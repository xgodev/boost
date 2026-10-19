// Command clients wires a NATS connection, a Pub/Sub client, gocql sessions
// of both drivers and a Bigtable client that inject builds from @Inject, and
// stops them with the app.
//
// Regenerate the modules in gen/ with:
//
//	go generate ./inject/examples/clients
package main

//go:generate go run github.com/xgodev/boost/inject/cmd/inject --boost.inject.output=gen

import (
	"context"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject/examples/clients/catalog"
	"github.com/xgodev/boost/inject/examples/clients/events"
	"github.com/xgodev/boost/inject/examples/clients/gen"
	"github.com/xgodev/boost/inject/examples/clients/ledger"
	"github.com/xgodev/boost/inject/examples/clients/orders"
	"github.com/xgodev/boost/inject/examples/clients/store"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

func main() {
	boost.Start()
	ctx := context.Background()

	var publisher *events.Publisher
	var topics *orders.Topics
	var cluster *store.Cluster
	var ledgerCluster *ledger.Cluster
	var prices *catalog.Prices
	app := fx.New(gen.Module(), fx.Populate(&publisher, &topics, &cluster, &ledgerCluster, &prices), fx.NopLogger)
	if err := app.Start(ctx); err != nil {
		log.Fatalf("%v", err)
	}

	if err := publisher.Publish("orders.created", []byte("42")); err != nil {
		log.Errorf("%v", err)
	}
	if err := topics.Create(ctx, "orders"); err != nil {
		log.Errorf("%v", err)
	}
	for _, c := range []interface{ Version() (string, error) }{cluster, ledgerCluster} {
		if version, err := c.Version(); err != nil {
			log.Errorf("%v", err)
		} else {
			log.Infof("cassandra %s", version)
		}
	}

	if err := prices.Set(ctx, "sku-1", "9.90"); err != nil {
		log.Errorf("%v", err)
	}

	// Stop closes the connection, the clients and the sessions.
	if err := app.Stop(ctx); err != nil {
		log.Fatalf("%v", err)
	}
}
