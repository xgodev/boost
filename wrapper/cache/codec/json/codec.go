package json

import (
	"bytes"
	"encoding/json"
)

type Codec[T any] struct {
}

func (c *Codec[T]) Encode(data T) (b []byte, err error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	if err = enc.Encode(data); err != nil {
		return b, err
	}

	b = buf.Bytes()
	return b, err
}

func (c *Codec[T]) Decode(b []byte, data *T) (err error) {
	var buf bytes.Buffer
	buf.Write(b)
	dec := json.NewDecoder(&buf)
	if err = dec.Decode(data); err != nil {
		return err
	}
	return nil
}

func New[T any]() *Codec[T] {
	return &Codec[T]{}
}
