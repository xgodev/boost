// Package plugin serves notes kept in a boost cache manager over freecache,
// with boost's log plugin on the manager: every Get and Set of the cache logs
// at TRACE with the cache name and key. PUT /notes/:key saves the text of a
// note, GET /notes/:key answers it, 404 for a key it does not have.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl -X PUT -H 'Content-Type: application/json' -d '{"text":"milk"}' localhost:8080/notes/shopping
//	curl localhost:8080/notes/shopping
package plugin

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
