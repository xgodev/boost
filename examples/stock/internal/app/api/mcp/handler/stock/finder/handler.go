// Package finder serves the find_stock tool of the mcp server of app api.
package finder

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/examples/stock/internal/pkg/domain/service/stock/finder"
)

// Handler serves the find_stock tool.
type Handler struct {
	svc finder.Service
}

// New provides the handler; the generated mcp door gets it from here.
// @Provide(index=0)
func New(svc finder.Service) *Handler {
	return &Handler{svc: svc}
}

// Input is the item the tool asks about.
type Input struct {
	SKU string `json:"sku" jsonschema:"the sku of the item"`
}

// Output is the stock of the item.
type Output struct {
	SKU   string `json:"sku"`
	Units int64  `json:"units"`
}

// FindStock tells how many units of an item are in stock.
// @MCP
func (h *Handler) FindStock(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, Output, error) {
	s, err := h.svc.Find(ctx, in.SKU)
	if err != nil {
		return nil, Output{}, err
	}
	return nil, Output{SKU: s.SKU, Units: s.Units}, nil
}
