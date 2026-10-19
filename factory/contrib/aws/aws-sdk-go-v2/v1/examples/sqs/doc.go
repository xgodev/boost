// Package sqs serves POST /messages, which sends the text of the body to a
// queue through the sqs client of boost's factory (boost.factory.aws).
//
// Run it with aws credentials, or next to a localstack set as the custom
// endpoint of SQS (boost.factory.aws.customEndpoint):
//
//	go run ./cmd/api
//	curl -d '{"text":"hi"}' -H 'Content-Type: application/json' localhost:8080/messages
package sqs

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
