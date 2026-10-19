// Package freecache serves notes kept in a boost cache manager over an
// in-memory freecache, with the string codec: PUT /notes/:key saves the text
// of a note, GET /notes/:key answers it, 404 for a key it does not have.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl -X PUT -H 'Content-Type: application/json' -d '{"text":"milk"}' localhost:8080/notes/shopping
//	curl localhost:8080/notes/shopping
package freecache

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
