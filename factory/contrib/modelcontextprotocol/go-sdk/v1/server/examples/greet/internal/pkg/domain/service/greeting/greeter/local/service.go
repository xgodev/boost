// Package local greets with hello.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/examples/greet/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/examples/greet/internal/pkg/domain/service/greeting/greeter"
	"github.com/xgodev/boost/wrapper/log"
)

type service struct{}

// New provides the greeter.
// @Provide(index=0)
func New() greeter.Service {
	return &service{}
}

func (s *service) Greet(ctx context.Context, name string) (*entity.Greeting, error) {
	log.FromContext(ctx).Infof("greeting %s", name)
	return &entity.Greeting{Message: "hello " + name}, nil
}
