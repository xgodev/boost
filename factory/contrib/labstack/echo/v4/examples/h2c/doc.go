// Package h2c serves GET / over HTTP/2 without TLS (h2c) and answers with
// the protocol, host, remote address, method and path of the request.
//
// Run it with the protocol set, then call it with HTTP/2 prior knowledge:
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl -v --http2-prior-knowledge http://localhost:8080
package h2c

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
