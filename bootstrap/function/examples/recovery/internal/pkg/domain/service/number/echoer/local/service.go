// Package local echoes the number it gets.
package local

import (
	"context"

	"github.com/xgodev/boost/bootstrap/function/examples/recovery/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/bootstrap/function/examples/recovery/internal/pkg/domain/service/number/echoer"
)

type service struct{}

// New provides the echoer.
// @Provide(index=0)
func New() echoer.Service {
	return &service{}
}

func (s *service) Echo(_ context.Context, n *entity.Number) (*entity.Number, error) {
	return &entity.Number{Value: n.Value}, nil
}
