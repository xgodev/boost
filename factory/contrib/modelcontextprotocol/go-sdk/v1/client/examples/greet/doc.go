// Package greet serves GET /greetings/:name, which calls the greet tool of
// server/examples/greet on the MCP client of boost's factory
// (boost.factory.mcp.client, http://localhost:8090/mcp by default).
//
// Run it with the server up:
//
//	go run ./cmd/api
//	curl localhost:8080/greetings/boost
package greet

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
