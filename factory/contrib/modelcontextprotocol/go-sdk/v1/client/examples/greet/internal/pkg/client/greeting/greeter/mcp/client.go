// Package mcp greets through the greet tool of an MCP server, on the MCP
// client of boost's factory (boost.factory.mcp.client).
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client/examples/greet/internal/pkg/domain/model/entity"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client/examples/greet/internal/pkg/domain/service/greeting/greeter"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client/plugins/local/wrapper/log"
	"github.com/xgodev/boost/model/errors"
)

type greeterClient struct {
	cli *client.Client
}

// New connects the MCP client with the log plugin and provides the greeter
// on it: the factory has no generated client, so the constructor calls it.
// @Provide(index=0)
func New(ctx context.Context) (greeter.Service, error) {
	cli, err := client.NewClient(ctx, log.Register)
	if err != nil {
		return nil, err
	}
	return &greeterClient{cli: cli}, nil
}

func (c *greeterClient) Greet(ctx context.Context, name string) (*entity.Greeting, error) {
	res, err := c.cli.Session().CallTool(ctx, &sdk.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": name},
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error calling the greet tool"))
	}
	if res.IsError {
		return nil, errors.Internalf("the greet tool failed: %v", res.Content)
	}
	out, _ := res.StructuredContent.(map[string]any)
	message, _ := out["greeting"].(string)
	return &entity.Greeting{Message: message}, nil
}
