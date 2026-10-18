// Command greet calls the greet tool of server/examples/greet.
package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	mcplog "github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client/plugins/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {
	ctx := context.Background()
	boost.Start()

	cli, err := client.NewClient(ctx, mcplog.Register)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	res, err := cli.Session().CallTool(ctx, &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "boost"},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%v", res.StructuredContent)
}
