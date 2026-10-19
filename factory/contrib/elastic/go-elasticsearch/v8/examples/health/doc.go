// Package health serves GET /ping and GET /health on an elasticsearch client
// of boost's factory built with the health plugin, which adds elasticsearch
// to the health checks.
//
// Run it next to an elasticsearch on 127.0.0.1:9200:
//
//	go run ./cmd/api
//	curl localhost:8080/health
package health

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
