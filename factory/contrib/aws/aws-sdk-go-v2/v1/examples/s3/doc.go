// Package s3 serves GET /files/:name, which tells when a file of a bucket was
// last modified, through the aws config of boost's factory (boost.factory.aws).
//
// Run it with aws credentials, or next to a localstack set as the custom
// endpoint of S3 (boost.factory.aws.customEndpoint):
//
//	go run ./cmd/api
//	curl localhost:8080/files/examplefile
package s3

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
