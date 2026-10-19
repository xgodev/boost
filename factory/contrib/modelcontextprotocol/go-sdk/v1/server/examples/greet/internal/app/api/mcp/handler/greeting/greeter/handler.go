// Package greeter serves the greet tool.
package greeter

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/examples/greet/internal/pkg/domain/service/greeting/greeter"
)

// Handler greets.
type Handler struct {
	svc greeter.Service
}

// New provides the handler to the generated mcp door.
// @Provide(index=0)
func New(svc greeter.Service) *Handler {
	return &Handler{svc: svc}
}

// Input is who to greet.
type Input struct {
	Name string `json:"name" jsonschema:"who to greet"`
}

// Output is the greeting.
type Output struct {
	Greeting string `json:"greeting"`
}

// Greet says hello.
// @MCP(name=greet)
func (h *Handler) Greet(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, Output, error) {
	greeting, err := h.svc.Greet(ctx, in.Name)
	if err != nil {
		return nil, Output{}, err
	}
	return nil, Output{Greeting: greeting.Message}, nil
}
