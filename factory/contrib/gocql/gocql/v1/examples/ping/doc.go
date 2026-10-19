// Package ping serves GET /ping, which queries cassandra through a session of
// boost's factory (boost.factory.gocql).
//
// Run it next to a cassandra on 127.0.0.1:9042:
//
//	go run ./cmd/api
//	curl localhost:8080/ping
package ping

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
