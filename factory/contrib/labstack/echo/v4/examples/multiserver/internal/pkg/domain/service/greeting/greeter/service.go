// Package greeter greets whoever calls the app.
package greeter

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/examples/multiserver/internal/pkg/domain/model/entity"
)

// Service greets.
type Service interface {
	Greet(ctx context.Context) (*entity.Greeting, error)
}
