// Command greet serves one MCP tool over streamable HTTP on :8090/mcp.
package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/plugins/contrib/go.opentelemetry.io/otel/v1"
	mcplog "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server/plugins/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

type greetIn struct {
	Name string `json:"name" jsonschema:"who to greet"`
}

type greetOut struct {
	Greeting string `json:"greeting"`
}

func greet(ctx context.Context, _ *mcp.CallToolRequest, in greetIn) (*mcp.CallToolResult, greetOut, error) {
	log.FromContext(ctx).Infof("greeting %s", in.Name)
	return nil, greetOut{Greeting: "hello " + in.Name}, nil
}

func main() {
	ctx := context.Background()
	boost.Start()

	srv, err := server.NewServer(ctx, mcplog.Register, otel.Register)
	if err != nil {
		log.Fatal(err)
	}
	mcp.AddTool(srv.Instance(), &mcp.Tool{Name: "greet", Description: "says hello"}, greet)

	srv.Serve(ctx)
}
