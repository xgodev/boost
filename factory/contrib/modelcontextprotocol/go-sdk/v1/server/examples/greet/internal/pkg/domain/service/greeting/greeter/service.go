// Package greeter greets someone.
package greeter

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/examples/greet/internal/pkg/domain/model/entity"
)

// Service greets someone.
type Service interface {
	Greet(ctx context.Context, name string) (*entity.Greeting, error)
}
