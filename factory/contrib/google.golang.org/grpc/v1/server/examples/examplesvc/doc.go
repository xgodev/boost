// Package examplesvc serves the Example service of example.proto on the grpc
// server of boost's factory (boost.factory.grpc.server, port 9090).
//
//	go run ./cmd/api
//
// client/examples/examplesvc calls it; pb is generated as README.md says.
package examplesvc

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
