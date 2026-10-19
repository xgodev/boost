// Package client serves POST /counters/:name, which increments a counter kept
// on a redis client of boost's factory. The health plugin of the client
// adds redis to GET /health.
//
// Run it next to a redis on 127.0.0.1:6379:
//
//	go run ./cmd/api
//	curl -X POST localhost:8080/counters/visits
//	curl localhost:8080/health
package client

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
