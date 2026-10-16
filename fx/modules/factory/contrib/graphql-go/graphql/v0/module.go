package graphql

import (
	"sync"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	g "github.com/xgodev/boost/factory/contrib/graphql-go/graphql/v0"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for the graphql http handler. The *graphql.Schema must be provided by the app.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				func(schema *graphql.Schema) *handler.Handler {
					return g.NewHandler(schema)
				},
			),
		)
	})

	return options
}
