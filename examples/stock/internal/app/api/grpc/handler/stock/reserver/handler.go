// Package reserver serves the StockReserver service of stock.v1.
package reserver

import (
	"context"

	stockv1 "github.com/xgodev/boost/examples/stock/api/proto/stock/v1"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/reserver"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
)

// Handler serves StockReserver.
type Handler struct {
	stockv1.UnimplementedStockReserverServer
	svc reserver.Service
}

// New provides the handler, which the generated grpc door registers with
// stockv1.RegisterStockReserverServer.
// @Provide(index=0)
// @GRPC(register=stockv1.RegisterStockReserverServer)
func New(svc reserver.Service) *Handler {
	return &Handler{svc: svc}
}

// Reserve takes one unit of the sku of the request.
func (h *Handler) Reserve(ctx context.Context, req *stockv1.ReserveRequest) (*stockv1.Stock, error) {
	left, err := h.svc.Reserve(ctx, req.GetSku())
	if err != nil {
		return nil, server.Error(err)
	}
	return &stockv1.Stock{Sku: left.SKU, Units: left.Units}, nil
}
