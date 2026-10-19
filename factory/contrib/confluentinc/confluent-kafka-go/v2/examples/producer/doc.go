// Package producer serves POST /messages, which produces the text of the body
// on a kafka topic through the producer of boost's factory
// (boost.factory.confluent).
//
// Run it next to a broker (boost.factory.confluent.brokers, default
// localhost:9092):
//
//	go run ./cmd/api
//	curl -d '{"text":"hi"}' -H 'Content-Type: application/json' localhost:8080/messages
package producer

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
