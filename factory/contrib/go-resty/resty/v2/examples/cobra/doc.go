// Package cobra was a command line with one subcommand per site; as a
// builder app each subcommand is a tool of an MCP server on stdio, which an
// MCP client starts as a command: the google tool calls Google on its resty
// client, and the bing tool calls Bing on another one that logs its
// requests.
//
// Run it on stdio with the client config:
//
//	CONF=configs/config.yaml go run ./cmd/api
package cobra

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
