// Package api composes app api.
package api

import (
	"context"

	"github.com/xgodev/boost/extra/middleware"
	fmiddleware "github.com/xgodev/boost/factory/local/extra/middleware"
	"github.com/xgodev/boost/factory/local/extra/middleware/examples/simple/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/local/extra/middleware/plugins/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/factory/local/extra/middleware/plugins/local/wrapper/log"
)

// Wrapper is what the greeter calls run through, built by boost's middleware
// factory with the name in boost.factory.local.extra.middleware.name: the log
// middleware, then the hystrix circuit breaker of that name. Each plugin is
// turned off by its boost.factory.local.extra.middleware.plugins.<plugin>.<name>.enabled.
// @Provide(index=0)
func Wrapper(ctx context.Context) (*middleware.AnyErrorWrapper[*entity.Greeting], error) {
	return fmiddleware.NewAnyErrorWrapper[*entity.Greeting](ctx,
		log.NewAnyError[*entity.Greeting],
		hystrix.NewAnyError[*entity.Greeting],
	)
}
