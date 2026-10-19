// Package api composes app api.
package api

import (
	"context"

	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/extra/middleware/examples/custom/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/extra/middleware/examples/custom/internal/pkg/middleware/counter"
)

// Counter is the custom middleware the greeter calls run through.
// @Provide(index=0)
func Counter() *counter.Counter[*entity.Greeting] {
	return counter.New[*entity.Greeting]()
}

// Wrapper is what the greeter calls run through: the counter.
// @Provide(index=0)
func Wrapper(ctx context.Context, c *counter.Counter[*entity.Greeting]) *middleware.AnyErrorWrapper[*entity.Greeting] {
	return middleware.NewAnyErrorWrapper[*entity.Greeting](ctx, "greeter", c)
}
