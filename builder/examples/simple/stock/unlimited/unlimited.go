// Package unlimited is a stock without limits.
package unlimited

import (
	"context"

	"github.com/xgodev/boost/builder/examples/simple/stock"
)

type service struct{}

// New provides a stock service in which every item is available.
// @Provide(index=0)
func New() stock.Service {
	return &service{}
}

func (*service) Available(context.Context, string) (bool, error) {
	return true, nil
}
