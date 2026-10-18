// Package simple is an application whose doors come from builder
// annotations: app api serves REST routes on two echo servers (rest and
// admin) and a gRPC health service; app worker records the orders of
// order.created events.
//
// Regenerate gen/ with:
//
//	go generate ./builder/examples/simple
package simple

//go:generate go run github.com/xgodev/boost/builder/cmd/builder --boost.builder.output=gen
