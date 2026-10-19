// Package local greets in process, with the greeting of the repository.
package local

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/health/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/health/internal/pkg/domain/repository/greeting"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/health/internal/pkg/domain/service/greeting/greeter"
)

type service struct {
	repo greeting.Repository
}

// New provides the greeter.
// @Provide(index=0)
func New(repo greeting.Repository) greeter.Service {
	return &service{repo: repo}
}

func (s *service) Greet(ctx context.Context) (*entity.Greeting, error) {
	return s.repo.Get(ctx)
}
