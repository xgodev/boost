// Package local greets through the middleware wrapper the app provides.
package local

import (
	"context"

	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/factory/local/extra/middleware/examples/simple/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/local/extra/middleware/examples/simple/internal/pkg/domain/service/greeting/greeter"
)

type service struct {
	wrapper *middleware.AnyErrorWrapper[*entity.Greeting]
}

// New provides the greeter on the wrapper its calls run through.
// @Provide(index=0)
func New(wrapper *middleware.AnyErrorWrapper[*entity.Greeting]) greeter.Service {
	return &service{wrapper: wrapper}
}

func (s *service) Greet(ctx context.Context, name string) (*entity.Greeting, error) {
	return s.wrapper.Exec(ctx, name,
		func(ctx context.Context) (*entity.Greeting, error) {
			return &entity.Greeting{Text: "hello, " + name}, nil
		},
		// the hystrix fallback, for a failed call or an open circuit
		func(_ context.Context, _ *entity.Greeting, err error) (*entity.Greeting, error) {
			return nil, err
		})
}
