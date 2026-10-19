// Package pubsub serves POST /messages, which sends the text of the body to a
// topic opened by boost's gocloud factory (boost.factory.gocloud).
//
// The factory opens an in-memory topic by default; type and resource point it
// at sqs, sns or google pubsub instead:
//
//	go run ./cmd/api
//	curl -d '{"text":"hi"}' -H 'Content-Type: application/json' localhost:8080/messages
//
//	BOOST_FACTORY_GOCLOUD_TYPE=sqs BOOST_FACTORY_GOCLOUD_REGION=us-east-1 \
//	BOOST_FACTORY_GOCLOUD_RESOURCE=https://sqs.us-east-1.amazonaws.com/000000000000/example \
//	go run ./cmd/api
package pubsub

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
