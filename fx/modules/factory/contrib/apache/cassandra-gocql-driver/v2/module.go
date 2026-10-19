package gocql

import (
	"context"
	"sync"

	"github.com/apache/cassandra-gocql-driver/v2"
	g "github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

type params struct {
	fx.In
	Plugins []g.Plugin `optional:"true"`
}

// Module provides the gocql v2 *gocql.Session, built from
// boost.factory.gocql.* with the []g.Plugin in the graph, and closes it on
// app stop. It yields its options on the first call only, so any number of
// dependents can include it.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			module(),
		)
	})

	return options
}

// module provides the session; it is the seam the tests call, since Module
// spends its once.
func module() fx.Option {
	return fx.Provide(
		func(ctx context.Context, lc fx.Lifecycle, p params) (*gocql.Session, error) {
			s, err := g.NewSession(ctx, p.Plugins...)
			if err != nil {
				return nil, err
			}
			lc.Append(fx.StopHook(s.Close))
			return s, nil
		},
	)
}
