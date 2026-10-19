// Package tls serves GET / over TLS with a certificate generated at the
// start for localhost and answers with the protocol, host, remote address,
// method and path of the request.
//
// Run it with TLS on port 8443, then call it trusting the certificate:
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl -vk https://localhost:8443
package tls

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
