// Package stocktool serves the stock tools of app api over MCP.
package stocktool

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/builder/examples/simple/stock"
)

// Handler serves the stock tools.
type Handler struct {
	stock stock.Service
}

// New provides the handler; the generated mcp door gets it from here.
// @Provide(index=0)
func New(stock stock.Service) *Handler {
	return &Handler{stock: stock}
}

// Input is the item a stock tool asks about.
type Input struct {
	Item string `json:"item" jsonschema:"the item to check"`
}

// Output tells whether the item can be ordered.
type Output struct {
	Available bool `json:"available"`
}

// CheckStock tells whether an item can be ordered.
// @MCP
func (h *Handler) CheckStock(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, Output, error) {
	available, err := h.stock.Available(ctx, in.Item)
	if err != nil {
		return nil, Output{}, err
	}
	return nil, Output{Available: available}, nil
}
