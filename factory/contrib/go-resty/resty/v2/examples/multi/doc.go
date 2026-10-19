// Package multi serves GET /google and GET /bing, each on its own resty
// client of boost's factory (boost.factory.resty.clients.google and
// boost.factory.resty.clients.bing), both logging their requests.
//
// Run it with the client config:
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl localhost:8080/google
//	curl localhost:8080/bing
package multi

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
