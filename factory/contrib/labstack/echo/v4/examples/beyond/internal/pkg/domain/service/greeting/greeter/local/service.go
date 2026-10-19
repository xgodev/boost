// Package local greets in process.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/beyond/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/beyond/internal/pkg/domain/service/greeting/greeter"
)

type service struct{}

// New provides the greeter wrapped in the log middleware, which logs every
// call to Greet: what a tracing advice of beyond did before each handler.
// @Provide(index=0)
// @Middleware(log)
func New() greeter.Service {
	return &service{}
}

func (s *service) Greet(context.Context) (*entity.Greeting, error) {
	return &entity.Greeting{Message: "Hello World!!"}, nil
}
