// Package multiecho serves one route on two echo servers, each on its own
// port; the second also answers whether both are up (the multi-server
// plugin).
//
// Run it with the port of the second server set:
//
//	CONF=configs/config.yaml go run ./cmd/api
package multiecho

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
