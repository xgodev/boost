package client

import (
	"sync"

	k "github.com/xgodev/boost/factory/contrib/k8s.io/client-go/v0"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for kubernetes clientset.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				k.NewClientset,
			),
		)
	})

	return options
}
