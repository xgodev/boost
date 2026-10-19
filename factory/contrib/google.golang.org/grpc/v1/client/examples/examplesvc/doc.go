// Package examplesvc serves POST /messages, which sends the message to the
// Example service through a grpc client of boost's factory
// (boost.factory.grpc.client) and answers with the reply.
//
// Run it next to server/examples/examplesvcautotls, which configs/config.yaml
// reaches over TLS on localhost:9090:
//
//	go run ./cmd/api
//	curl -d '{"message":"hi"}' -H 'Content-Type: application/json' localhost:8080/messages
package examplesvc

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
