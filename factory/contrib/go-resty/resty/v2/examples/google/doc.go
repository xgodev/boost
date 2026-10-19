// Package google serves GET /google, which calls Google on the resty client
// of boost's factory. The health plugin of the client adds Google to
// GET /health.
//
// Run it with the client config:
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl localhost:8080/google
//	curl localhost:8080/health
package google

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
