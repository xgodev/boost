// Package kinesis serves POST /messages, which puts the text of the body on a
// stream through the kinesis client of boost's factory (boost.factory.aws).
//
// Run it with aws credentials, or next to a localstack set as the custom
// endpoint of Kinesis (boost.factory.aws.customEndpoint):
//
//	go run ./cmd/api
//	curl -d '{"text":"hi"}' -H 'Content-Type: application/json' localhost:8080/messages
package kinesis

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
