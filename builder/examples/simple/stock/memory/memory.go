// Package memory keeps the stock in memory. Its methods are plain: builder
// runs the middlewares @Middleware lists around each of them.
package memory

import (
	"context"

	"github.com/xgodev/boost/builder/examples/simple/stock"
)

type service struct {
	soldOut map[string]bool
}

// New provides the default stock service, out of the items Options lists;
// every consumer receives it wrapped in the log and prometheus middlewares.
// @Provide(index=0)
// @Default
// @Middleware(log, prometheus)
func New(o *Options) stock.Service {
	soldOut := map[string]bool{}
	for _, item := range o.SoldOut {
		soldOut[item] = true
	}
	return &service{soldOut: soldOut}
}

func (s *service) Available(_ context.Context, item string) (bool, error) {
	return !s.soldOut[item], nil
}
