// Package simple serves GET /health with boost's health registry: it adds a
// checker of its own, tcp, which dials the dependency at app.dependency.addr,
// and answers with the result of every checker, 503 when a required one
// fails.
//
//	APP_DEPENDENCY_ADDR=localhost:6379 go run ./cmd/api
//	curl -i localhost:8080/health
package simple

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
