// Package ping serves GET /ping, which pings mongo through the connection the
// generator builds on boost's factory (boost.factory.mongo).
//
// Run it next to a mongo on localhost:27017:
//
//	go run ./cmd/api
//	curl localhost:8080/ping
package ping

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
