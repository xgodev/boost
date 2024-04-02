package newrelic

import (
	"sync"

	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	newrelic "github.com/xgodev/boost/factory/contrib/newrelic/go-agent/v3"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for newrelic agent.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Invoke(
				newrelic.NewApplication,
			),
		)
	})

	return options
}
