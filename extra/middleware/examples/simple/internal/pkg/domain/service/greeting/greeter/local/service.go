// Package local greets through the log middleware and a hystrix circuit
// breaker.
package local

import (
	"context"

	hystrixgo "github.com/afex/hystrix-go/hystrix"
	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/extra/middleware/examples/simple/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/extra/middleware/examples/simple/internal/pkg/domain/service/greeting/greeter"
	"github.com/xgodev/boost/extra/middleware/plugins/contrib/afex/hystrix-go/v0"
	"github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/log"
)

type service struct {
	wrapper *middleware.AnyErrorWrapper[*entity.Greeting]
}

// New provides the greeter. Its calls run through the wrapper it builds: the
// log middleware, then the hystrix circuit breaker greeter.
// @Provide(index=0)
func New(ctx context.Context) greeter.Service {
	return &service{wrapper: middleware.NewAnyErrorWrapper[*entity.Greeting](ctx, "greeter",
		log.NewAnyErrorMiddleware[*entity.Greeting](ctx),
		hystrix.NewAnyErrorMiddlewareWithConfig[*entity.Greeting](ctx, "greeter", hystrixgo.CommandConfig{
			Timeout:               1000,
			MaxConcurrentRequests: 100,
			ErrorPercentThreshold: 50,
		}),
	)}
}

func (s *service) Greet(ctx context.Context, name string) (*entity.Greeting, error) {
	return s.wrapper.Exec(ctx, name,
		func(ctx context.Context) (*entity.Greeting, error) {
			return &entity.Greeting{Text: "hello, " + name}, nil
		},
		// the fallback of the circuit breaker, for a failed call or an open
		// circuit: an app could answer from a cache here
		func(_ context.Context, _ *entity.Greeting, err error) (*entity.Greeting, error) {
			return nil, err
		})
}
