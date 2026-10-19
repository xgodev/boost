// Package msgpack is a wrapper/cache codec on github.com/shamaton/msgpack/v3.
package msgpack

import (
	"github.com/shamaton/msgpack/v3"
)

// Codec encodes cache values as MessagePack with shamaton/msgpack v3.
type Codec[T any] struct {
}

// Encode serializes data as MessagePack.
func (c *Codec[T]) Encode(data T) (b []byte, err error) {
	b, err = msgpack.Marshal(data)
	return b, err
}

// Decode deserializes MessagePack bytes into data. A timestamp decodes into
// a time.Time in UTC (v2 decoded into Local); the instant is the same.
func (c *Codec[T]) Decode(b []byte, data *T) (err error) {
	err = msgpack.Unmarshal(b, data)
	return err
}

// New returns a shamaton/msgpack v3 codec for T.
func New[T any]() *Codec[T] {
	return &Codec[T]{}
}
