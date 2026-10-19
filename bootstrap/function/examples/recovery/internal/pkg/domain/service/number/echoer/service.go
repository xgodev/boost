// Package echoer echoes numbers.
package echoer

import (
	"context"

	"github.com/xgodev/boost/bootstrap/function/examples/recovery/internal/pkg/domain/model/entity"
)

// Service echoes numbers.
type Service interface {
	Echo(ctx context.Context, n *entity.Number) (*entity.Number, error)
}
