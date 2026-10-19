// Package api composes app api.
package api

import (
	"context"

	hystrixgo "github.com/afex/hystrix-go/hystrix"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/extra/middleware/examples/service/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/extra/middleware/plugins/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/log"
)

// Wrapper is what the greeter calls run through: the log middleware, then the
// hystrix circuit breaker greeter.
// @Provide(index=0)
func Wrapper(ctx context.Context) *middleware.AnyErrorWrapper[*entity.Greeting] {
	return middleware.NewAnyErrorWrapper[*entity.Greeting](ctx, "greeter",
		log.NewAnyErrorMiddleware[*entity.Greeting](ctx),
		hystrix.NewAnyErrorMiddlewareWithConfig[*entity.Greeting](ctx, "greeter", hystrixgo.CommandConfig{
			Timeout:               1000,
			MaxConcurrentRequests: 100,
			ErrorPercentThreshold: 50,
		}),
	)
}
