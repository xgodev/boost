// Package greet serves the greet tool on an MCP server over streamable HTTP
// on :8090/mcp, with the log and OpenTelemetry plugins of the server.
//
//	go run ./cmd/api
//
// client/examples/greet calls it.
package greet

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
