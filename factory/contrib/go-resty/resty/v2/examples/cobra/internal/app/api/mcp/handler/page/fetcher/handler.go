// Package fetcher serves the google and bing tools.
package fetcher

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/examples/cobra/internal/pkg/domain/service/page/fetcher"
)

// Handler fetches a page and tells how it answered.
type Handler struct {
	google fetcher.Service
	bing   fetcher.Service
}

// New provides the handler to the generated mcp door with both fetchers.
// @Provide(index=0)
// @Inject(name=google, index=0)
// @Inject(name=bing, index=1)
func New(google, bing fetcher.Service) *Handler {
	return &Handler{google: google, bing: bing}
}

// Input of the tools: they take no argument.
type Input struct{}

// Output tells how the site answered.
type Output struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

// Google calls Google.
// @MCP
func (h *Handler) Google(ctx context.Context, _ *mcp.CallToolRequest, _ Input) (*mcp.CallToolResult, Output, error) {
	return fetch(ctx, h.google)
}

// Bing calls Bing.
// @MCP
func (h *Handler) Bing(ctx context.Context, _ *mcp.CallToolRequest, _ Input) (*mcp.CallToolResult, Output, error) {
	return fetch(ctx, h.bing)
}

func fetch(ctx context.Context, svc fetcher.Service) (*mcp.CallToolResult, Output, error) {
	page, err := svc.Fetch(ctx)
	if err != nil {
		return nil, Output{}, err
	}
	return nil, Output{URL: page.URL, Status: page.Status}, nil
}
