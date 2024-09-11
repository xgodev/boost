package msgpack

import (
	"github.com/vmihailenco/msgpack/v5"
)

type Codec[T any] struct {
}

// Encode serializa os dados como MsgPack (formato binário mais compacto que JSON)
func (c *Codec[T]) Encode(data T) (b []byte, err error) {
	b, err = msgpack.Marshal(data)
	return b, err
}

// Decode desserializa os dados de MsgPack para o tipo T
func (c *Codec[T]) Decode(b []byte, data *T) (err error) {
	err = msgpack.Unmarshal(b, data)
	return err
}

func New[T any]() *Codec[T] {
	return &Codec[T]{}
}
