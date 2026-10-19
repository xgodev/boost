// Package freecachebinary serves the last reading of each sensor, kept in a
// boost cache manager over freecache with the msgpack codec: the datastore
// stores a reading as a packed struct. PUT /readings/:sensor saves a
// reading, GET /readings/:sensor answers it, 404 for a sensor it does not
// have.
//
//	CONF=configs/config.yaml go run ./cmd/api
//	curl -X PUT -H 'Content-Type: application/json' -d '{"location":1233,"timestamp":123452123,"temperature":12}' localhost:8080/readings/1
//	curl localhost:8080/readings/1
package freecachebinary

//go:generate go run github.com/xgodev/boost/builder/cmd/builder
