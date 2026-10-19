// Package health serves GET /ping and GET /health on a cassandra session of
// boost's factory built with the health plugin, which adds cassandra to the
// health checks.
//
// Run it next to a cassandra on 127.0.0.1:9042:
//
//	go run ./cmd/api
//	curl localhost:8080/health
package health

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
