package provider

import (
	"sync"

	"github.com/xgodev/boost/bootstrap/datastore"
	"go.uber.org/fx"
)

var once sync.Once

// Module returns fx module for initialization of event wrapper provider.
//
// The module is only loaded once.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			datastore.EventModule(),
			fx.Provide(
				NewEventWrapperProvider,
			),
		)
	})

	return options
}
