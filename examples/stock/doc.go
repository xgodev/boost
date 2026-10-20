// Package stock is the complete example app of boost: a small stock with
// reservation, built from every boost resource it can use. It is a living
// example: a resource that lands in boost is added here in the same PR.
//
// App api serves the same three use cases (create, find, reserve) on REST
// (echo, :8080), gRPC (:9090, api/proto/stock/v1) and MCP (:8090/mcp, the
// find_stock tool). App worker reserves the item of each order.created
// CloudEvent it receives (HTTP on :8081 or NATS). A reservation is announced
// as a stock.reserved event on boost's publisher (noop by default). App
// stock-check is a job: it logs the units of the skus it is given and exits,
// 1 when a sku is not in stock.
//
// Defaults need no infra: the stock is in memory (seeded with tshirt=3 and
// mug=1) and the publisher is noop. Config chooses the rest:
//
//	INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=mongodb        stock in MongoDB (boost.factory.mongo.*)
//	INTERNAL_PKG_DOMAIN_REPOSITORY_STOCK_IMPLEMENTATION=oracle         stock in Oracle (boost.factory.godror.*, needs cgo)
//	INTERNAL_PKG_DOMAIN_SERVICE_STOCK_RESERVER_IMPLEMENTATION=resty    reserve by calling app api
//	BOOST_WRAPPER_PUBLISHER_DRIVER_NAME=nats                           announce on NATS (boost.factory.nats.*)
//	WORKER_FUNCTION_ADAPTER=nats                                       worker reads order.created from NATS
//
// Run, from this directory (in GoLand: a Go Build configuration on the
// package, working directory here, and the CONF variable):
//
//	CONF=configs/api.yaml go run ./cmd/api
//	CONF=configs/worker.yaml go run ./cmd/worker
//	INTERNAL_APP_STOCKCHECK_JOB_STOCK_CHECKER_SKUS=tshirt,mug go run ./cmd/stock-check
//
// and call it (README.md shows what each call answers):
//
//	curl localhost:8080/stocks/tshirt
//	curl -XPOST localhost:8080/stocks -d '{"sku":"cap","units":2}' -H 'Content-Type: application/json'
//	curl -XPOST localhost:8080/stocks/tshirt/reserve
//	grpcurl -plaintext -d '{"sku":"mug"}' localhost:9090 stock.v1.StockReserver/Reserve
//	curl localhost:8081 -H 'Ce-Id: 1' -H 'Ce-Source: shop' -H 'Ce-Type: order.created' -H 'Ce-Specversion: 1.0' -H 'Content-Type: application/json' -d '{"sku":"tshirt"}'
//
// Regenerate gen/ and the mains with go generate ./... (and
// api/proto/stock/v1 with the protoc line of api/proto/stock/v1/doc.go).
package stock

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
