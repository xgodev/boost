// Package finder serves the StockFinder service of stock.v1.
package finder

import (
	"context"

	stockv1 "github.com/xgodev/boost/examples/stock/api/proto/stock/v1"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/finder"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
)

// Handler serves StockFinder.
type Handler struct {
	stockv1.UnimplementedStockFinderServer
	svc finder.Service
}

// New provides the handler, which the generated grpc door registers with
// stockv1.RegisterStockFinderServer.
// @Provide(index=0)
// @GRPC(register=stockv1.RegisterStockFinderServer)
func New(svc finder.Service) *Handler {
	return &Handler{svc: svc}
}

// Find returns the stock of the sku of the request.
func (h *Handler) Find(ctx context.Context, req *stockv1.FindRequest) (*stockv1.Stock, error) {
	s, err := h.svc.Find(ctx, req.GetSku())
	if err != nil {
		return nil, server.Error(err)
	}
	return &stockv1.Stock{Sku: s.SKU, Units: s.Units}, nil
}
