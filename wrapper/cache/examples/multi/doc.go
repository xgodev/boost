// Package multi serves notes kept in a boost cache manager over two
// freecache layers, a near one in front of a far one: a saved note goes to
// both, a read looks in the near layer first, and a note the far layer alone
// has is copied to the near one on the way back. PUT /notes/:key saves the
// text of a note, GET /notes/:key answers it, 404 for a key neither layer
// has.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl -X PUT -H 'Content-Type: application/json' -d '{"text":"milk"}' localhost:8080/notes/shopping
//	curl localhost:8080/notes/shopping
package multi

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
