package prometheus

import (
	"sync"

	p "github.com/xgodev/boost/factory/contrib/prometheus/client_golang/v1"
	"go.uber.org/fx"
)

var once sync.Once

// PusherModule fx module that starts the prometheus push gateway pusher
// (no-op unless boost.factory.prometheus.pushGateway.enabled).
func PusherModule() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fx.Invoke(
				p.StartPusher,
			),
		)
	})

	return options
}
