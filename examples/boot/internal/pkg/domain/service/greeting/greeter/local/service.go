// Package local greets in process.
package local

import (
	"context"

	"github.com/xgodev/boost/examples/boot/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/examples/boot/internal/pkg/domain/service/greeting/greeter"
	"github.com/xgodev/boost/wrapper/log"
)

type service struct{}

// New provides the greeter.
// @Provide(index=0)
func New() greeter.Service {
	return &service{}
}

// Greet logs the greeting with the logger of ctx, the one boost.Start set.
func (s *service) Greet(ctx context.Context, name string) (*entity.Greeting, error) {
	text := "hello, " + name
	log.FromContext(ctx).Debugf("greeting %s", name)
	return &entity.Greeting{Text: text}, nil
}
