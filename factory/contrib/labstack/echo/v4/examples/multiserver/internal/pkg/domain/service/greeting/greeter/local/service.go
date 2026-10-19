// Package local greets in process, with the configured message.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/multiserver/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/multiserver/internal/pkg/domain/service/greeting/greeter"
)

type service struct {
	message string
}

// New provides the greeter.
// @Provide(index=0)
func New(o *Options) greeter.Service {
	return &service{message: o.Message}
}

func (s *service) Greet(context.Context) (*entity.Greeting, error) {
	return &entity.Greeting{Message: s.message}, nil
}
