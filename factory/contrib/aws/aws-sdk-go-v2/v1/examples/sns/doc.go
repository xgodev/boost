// Package sns serves POST /messages, which publishes the text of the body on
// a topic through the sns client of boost's factory (boost.factory.aws).
//
// Run it with aws credentials, or next to a localstack set as the custom
// endpoint of SNS (boost.factory.aws.customEndpoint):
//
//	go run ./cmd/api
//	curl -d '{"text":"hi"}' -H 'Content-Type: application/json' localhost:8080/messages
package sns

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
