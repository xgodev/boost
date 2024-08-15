package goka

import (
	fxgoka "github.com/xgodev/boost/fx/modules/factory/contrib/lovoo/goka/v1"
	"github.com/xgodev/boost/wrapper/publisher/driver/contrib/lovoo/goka/v1"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fxgoka.Module(),
			fx.Provide(
				goka.New,
			),
		)
	})

	return options
}
