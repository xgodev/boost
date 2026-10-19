// Package examplesvcautotls serves the Example service of
// server/examples/examplesvc over TLS, on a certificate the grpc server of
// boost's factory generates for localhost (configs/config.yaml).
//
//	go run ./cmd/api
package examplesvcautotls

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
