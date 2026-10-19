// Package connect serves GET /ping, which pings the elasticsearch client of
// boost's factory and answers with the status of the cluster.
//
// Run it next to an elasticsearch on 127.0.0.1:9200:
//
//	go run ./cmd/api
//	curl localhost:8080/ping
package connect

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
